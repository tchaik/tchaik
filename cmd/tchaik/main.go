// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchaik creates a webserver which serves the web UI.

It is assumed that tchaik is run relatively local to the user (i.e. serving pages to the local machine, or a local
network).

All configuration is done through command line parameters.

A common use case is to begin by use using an existing iTunes Library file:

  tchaik -itlXML /path/to/iTunesMusicLibrary.xml

*/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"tchaik.com/index"
	"tchaik.com/index/attr"
	"tchaik.com/index/checklist"
	"tchaik.com/index/favourite"
	"tchaik.com/index/history"
	"tchaik.com/index/itl"
	"tchaik.com/index/walk"
	"tchaik.com/store"
	"tchaik.com/store/cmdflag"
)

var debug bool
var itlXML, tchLib, walkPath string

var playHistoryPath, favouritesPath, checklistPath string

var listenAddr string
var uiDir string
var certFile, keyFile string

var authUser, authPassword string

var traceListenAddr string

func init() {
	flag.BoolVar(&debug, "debug", false, "print debugging information")

	flag.StringVar(&listenAddr, "listen", "localhost:8080", "bind `address` for main HTTP server")
	flag.StringVar(&certFile, "tls-cert", "", "certificate `file`, must also specify -tls-key")
	flag.StringVar(&keyFile, "tls-key", "", "certificate key `file`, must also specify -tls-cert")

	flag.StringVar(&itlXML, "itlXML", "", "iTunes Library XML `file`")
	flag.StringVar(&tchLib, "lib", "", "Tchaik library `file`")
	flag.StringVar(&walkPath, "path", "", "`directory` containing music files")

	flag.StringVar(&playHistoryPath, "play-history", "history.json", "play history `file`")
	flag.StringVar(&favouritesPath, "favourites", "favourites.json", "favourites `file`")
	flag.StringVar(&checklistPath, "checklist", "checklist.json", "checklist `file`")

	flag.StringVar(&uiDir, "ui-dir", "ui", "UI asset `directory`")

	flag.StringVar(&authUser, "auth-user", "", "`user` to use for HTTP authentication (set to enable)")
	flag.StringVar(&authPassword, "auth-password", "", "`password` to use for HTTP authentication")

	flag.StringVar(&traceListenAddr, "trace-listen", "", "bind `address` for trace HTTP server")
}

func readLibrary() (index.Library, error) {
	var count int
	check := func(x string) {
		if x != "" {
			count++
		}
	}
	check(itlXML)
	check(tchLib)
	check(walkPath)

	switch {
	case count == 0:
		return nil, fmt.Errorf("must specify one library file or a path to build one from (-itlXML, -lib or -path)")
	case count > 1:
		return nil, fmt.Errorf("must only specify one library file or a path to build one from (-itlXML, -lib or -path)")
	}

	var lib index.Library
	switch {
	case tchLib != "":
		f, err := os.Open(tchLib)
		if err != nil {
			return nil, fmt.Errorf("could not open Tchaik library file: %v", err)
		}
		defer f.Close()

		fmt.Printf("Parsing %v...", tchLib)
		lib, err = index.ReadFrom(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing Tchaik library file: %v\n", err)
		}
		fmt.Println("done.")
		return lib, nil

	case itlXML != "":
		f, err := os.Open(itlXML)
		if err != nil {
			return nil, fmt.Errorf("could open iTunes library file: %v", err)
		}
		defer f.Close()

		lib, err = itl.ReadFrom(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing iTunes library file: %v", err)
		}

	case walkPath != "":
		fmt.Printf("Walking %v...\n", walkPath)
		lib = walk.NewLibrary(walkPath)
		fmt.Println("Finished walking.")
	}

	fmt.Printf("Building Tchaik Library...")
	lib = index.Convert(lib, "ID")
	fmt.Println("done.")
	return lib, nil
}

func buildRootCollection(l index.Library) index.Collection {
	root := index.Collect(l, index.By(attr.String("Album")))
	index.SortKeysByGroupName(root)
	return root
}

func buildSearchIndex(c index.Collection) index.Searcher {
	wi := index.BuildCollectionWordIndex(c, []string{"Composer", "Artist", "Album", "Name"})
	return index.FlatSearcher{
		Searcher: index.WordsIntersectSearcher(index.BuildPrefixExpandSearcher(wi, wi, 10)),
	}
}

func main() {
	flag.Parse()
	l, err := readLibrary()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Building root collection...")
	root := buildRootCollection(l)
	fmt.Println("done.")

	fmt.Printf("Building artists filter...")
	artists := index.Filter(root, "Artist")
	fmt.Println("done.")

	fmt.Printf("Building composers filter...")
	composers := index.Filter(root, "Composer")
	fmt.Println("done.")

	fmt.Printf("Building recent index...")
	recent := index.Recent(root, 150)
	fmt.Println("done.")

	fmt.Printf("Building search index...")
	searcher := buildSearchIndex(root)
	fmt.Println("done.")

	fmt.Printf("Loading play history...")
	hs, err := history.NewStore(playHistoryPath)
	if err != nil {
		fmt.Printf("\nerror loading play history: %v", err)
		os.Exit(1)
	}
	fmt.Println("done.")

	fmt.Printf("Loading favourites...")
	favourites, err := favourite.NewStore(favouritesPath)
	if err != nil {
		fmt.Printf("\nerror loading favourites: %v", err)
		os.Exit(1)
	}
	fmt.Println("done.")

	fmt.Printf("Loading checklist...")
	checklist, err := checklist.NewStore(checklistPath)
	if err != nil {
		fmt.Printf("\nerror loading checklist: %v", err)
		os.Exit(1)
	}
	fmt.Println("done.")

	// fmt.Printf("Loading playlists...")
	// ps, err := playlist.NewStore(playListPath)
	// if err != nil {
	// 	fmt.Printf("\nerror loading playlists: %v", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("done")

	mediaFileSystem, artworkFileSystem, err := cmdflag.Stores()
	if err != nil {
		fmt.Println("error setting up stores:", err)
		os.Exit(1)
	}

	if debug {
		mediaFileSystem = store.LogFileSystem("Media", mediaFileSystem)
		artworkFileSystem = store.LogFileSystem("Artwork", artworkFileSystem)
	}

	if traceListenAddr != "" {
		fmt.Printf("Starting trace server on http://%v\n", traceListenAddr)
		go func() {
			log.Fatal(http.ListenAndServe(traceListenAddr, nil))
		}()
	}

	lib := Library{
		Library: l,
		collections: map[string]index.Collection{
			"Root": root,
		},
		filters: map[string][]index.FilterItem{
			"Artist":   artists,
			"Composer": composers,
		},
		recent:   recent,
		searcher: searcher,
	}

	meta := &Meta{
		history:    hs,
		favourites: favourites,
		checklist:  checklist,
	}

	h := NewHandler(lib, meta, mediaFileSystem, artworkFileSystem)

	if certFile != "" && keyFile != "" {
		fmt.Printf("Web server is running on https://%v\n", listenAddr)
		fmt.Println("Quit the server with CTRL-C.")

		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS10}
		server := &http.Server{Addr: listenAddr, Handler: h, TLSConfig: tlsConfig}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	}

	fmt.Printf("Web server is running on http://%v\n", listenAddr)
	fmt.Println("Quit the server with CTRL-C.")

	log.Fatal(http.ListenAndServe(listenAddr, h))
}
