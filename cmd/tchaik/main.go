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

	"tchaik.com/index/itl"
	"tchaik.com/index/walk"
	"tchaik.com/store"
	"tchaik.com/store/cmdflag"
)

var debug bool
var itlXML, tchLib, walkPath string

var playHistoryPath, favouritesPath, checklistPath, playlistPath, cursorPath string

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
	flag.StringVar(&playlistPath, "playlists", "playlists.json", "playlists `file`")
	flag.StringVar(&cursorPath, "cursors", "cursors.json", "cursors `file`")

	flag.StringVar(&uiDir, "ui-dir", "ui", "UI asset `directory`")

	flag.StringVar(&authUser, "auth-user", "", "`user` to use for HTTP authentication (set to enable)")
	flag.StringVar(&authPassword, "auth-password", "", "`password` to use for HTTP authentication")

	flag.StringVar(&traceListenAddr, "trace-listen", "", "bind `address` for trace HTTP server")
}

type assignedCount int

func (e *assignedCount) check(list ...string) {
	for _, x := range list {
		if x != "" {
			*e++
		}
	}
}

func readLibrary() (index.Library, error) {
	e := assignedCount(0)
	e.check(itlXML, tchLib, walkPath)

	switch {
	case e == 0:
		return nil, fmt.Errorf("must specify one library file or a path to build one from (-itlXML, -lib or -path)")
	case e > 1:
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

	fmt.Printf("Processing artist names and composers...")
	rootSplit := index.SubTransform(root, index.SplitList("Artist", "Composer"))
	fmt.Println("done.")

	fmt.Printf("Building artists filter...")
	artists := index.Filter(rootSplit, attr.Strings("Artist"))
	fmt.Println("done.")

	fmt.Printf("Building composers filter...")
	composers := index.Filter(rootSplit, attr.Strings("Composer"))
	fmt.Println("done.")

	fmt.Printf("Building recent index...")
	recent := index.Recent(root, 150)
	fmt.Println("done.")

	meta, err := loadLocalMeta()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
		searcher: newBootstrapSearcher(root),
	}

	h := NewHandler(lib, meta, mediaFileSystem, artworkFileSystem)

	if certFile != "" && keyFile != "" {
		fmt.Printf("Web server is running on https://%v\n", listenAddr)
		fmt.Println("Quit the server with CTRL-C.")

		server := &http.Server{
			Addr:    listenAddr,
			Handler: h,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS10,
			},
		}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	}

	fmt.Printf("Web server is running on http://%v\n", listenAddr)
	fmt.Println("Quit the server with CTRL-C.")

	log.Fatal(http.ListenAndServe(listenAddr, h))
}
