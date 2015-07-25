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

	"github.com/tchaik/tchaik/index"
	"github.com/tchaik/tchaik/index/attr"
	"github.com/tchaik/tchaik/index/history"
	"github.com/tchaik/tchaik/index/itl"
	"github.com/tchaik/tchaik/index/walk"
	"github.com/tchaik/tchaik/store"
	"github.com/tchaik/tchaik/store/cmdflag"
)

var debug bool
var itlXML, tchLib, walkPath string

var playHistoryPath string

var listenAddr string
var staticDir string
var certFile, keyFile string

var authUser, authPassword string

var traceListenAddr string

func init() {
	flag.BoolVar(&debug, "debug", false, "print debugging information")

	flag.StringVar(&listenAddr, "listen", "localhost:8080", "bind address to http listen")
	flag.StringVar(&certFile, "tls-cert", "", "path to a certificate file, must also specify -tls-key")
	flag.StringVar(&keyFile, "tls-key", "", "path to a certificate key file, must also specify -tls-cert")

	flag.StringVar(&itlXML, "itlXML", "", "path to iTunes Library XML file")
	flag.StringVar(&tchLib, "lib", "", "path to Tchaik library file")
	flag.StringVar(&walkPath, "path", "", "path to directory containing music files (to build index from)")

	flag.StringVar(&playHistoryPath, "play-history", "history.json", "path to play history file")

	flag.StringVar(&staticDir, "static-dir", "ui/static", "Path to the static asset directory")

	flag.StringVar(&authUser, "auth-user", "", "username to use for HTTP authentication (set to enable)")
	flag.StringVar(&authPassword, "auth-password", "", "password to use for HTTP authentication")

	flag.StringVar(&traceListenAddr, "trace-listen", "", "bind address for trace HTTP server")
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
			"Artist": artists,
		},
		recent:   recent,
		searcher: searcher,
	}

	h := NewHandler(lib, hs, mediaFileSystem, artworkFileSystem)

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
