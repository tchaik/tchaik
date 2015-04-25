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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhowden/httpauth"
	"github.com/dhowden/itl"

	"github.com/dhowden/tchaik/index"
	"github.com/dhowden/tchaik/store"
	"github.com/dhowden/tchaik/store/cmdflag"
)

var debug bool
var itlXML, tchLib string

var listenAddr string
var certFile, keyFile string

var auth bool

func init() {
	flag.BoolVar(&debug, "debug", false, "print debugging information")

	flag.StringVar(&listenAddr, "listen", "localhost:8080", "bind address to http listen")
	flag.StringVar(&certFile, "tls-cert", "", "path to a certificate file, must also specify -tls-key")
	flag.StringVar(&keyFile, "tls-key", "", "path to a certificate key file, must also specify -tls-cert")

	flag.StringVar(&itlXML, "itlXML", "", "path to iTunes Library XML file")
	flag.StringVar(&tchLib, "lib", "", "path to Tchaik library file")

	flag.BoolVar(&auth, "auth", false, "use basic HTTP authentication")
}

var creds = httpauth.Creds(map[string]string{
	"user": "password",
})

func readLibrary() (index.Library, error) {
	if itlXML == "" && tchLib == "" {
		return nil, fmt.Errorf("must specify one library file (-itlXML or -lib)")
	}

	if itlXML != "" && tchLib != "" {
		return nil, fmt.Errorf("must only specify one library file (-itlXML or -lib)")
	}

	var l index.Library
	if itlXML != "" {
		f, err := os.Open(itlXML)
		if err != nil {
			return nil, fmt.Errorf("could not open iTunes library file: %v", err)
		}

		fmt.Printf("Parsing %v...", itlXML)
		it, err := itl.ReadFromXML(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing iTunes library file: %v", err)
		}
		f.Close()
		fmt.Println("done.")

		fmt.Printf("Building Tchaik Library...")
		l = index.Convert(index.NewITunesLibrary(&it), "TrackID")
		fmt.Println("done.")
		return l, nil
	}

	f, err := os.Open(tchLib)
	if err != nil {
		return nil, fmt.Errorf("could not open Tchaik library file: %v", err)
	}

	fmt.Printf("Parsing %v...", tchLib)
	l, err = index.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing Tchaik library file: %v\n", err)
	}
	fmt.Println("done.")
	return l, nil
}

func buildRootCollection(l index.Library) index.Collection {
	root := index.Collect(l, index.ByAttr(index.StringAttr("Album")))
	index.SortKeysByGroupName(root)
	return root
}

func buildSearchIndex(c index.Collection) index.Searcher {
	wi := index.BuildWordIndex(c, []string{"Composer", "Artist", "Album", "Name"})
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

	fmt.Printf("Building search index...")
	searcher := buildSearchIndex(root)
	fmt.Println("done.")

	mediaFileSystem, artworkFileSystem, err := cmdflag.Stores()
	if err != nil {
		fmt.Println("error setting up stores:", err)
		os.Exit(1)
	}

	if debug {
		mediaFileSystem = store.LogFileSystem{
			Name:       "Media",
			FileSystem: mediaFileSystem,
		}
		artworkFileSystem = store.LogFileSystem{
			Name:       "Artwork",
			FileSystem: artworkFileSystem,
		}
	}

	libAPI := LibraryAPI{
		Library: l,
		collections: map[string]index.Collection{
			"Root": root,
		},
		filters: map[string][]index.FilterItem{
			"Artist": artists,
		},
		searcher: searcher,
		sessions: newSessions(),
	}

	mediaFileSystem = libAPI.FileSystem(mediaFileSystem)
	artworkFileSystem = libAPI.FileSystem(artworkFileSystem)

	m := buildMainHandler(libAPI, mediaFileSystem, artworkFileSystem)

	if certFile != "" && keyFile != "" {
		fmt.Printf("Web server is running on https://%v\n", listenAddr)
		fmt.Println("Quit the server with CTRL-C.")

		log.Fatal(http.ListenAndServeTLS(listenAddr, certFile, keyFile, m))
	}

	fmt.Printf("Web server is running on http://%v\n", listenAddr)
	fmt.Println("Quit the server with CTRL-C.")

	log.Fatal(http.ListenAndServe(listenAddr, m))
}

func buildMainHandler(l LibraryAPI, mediaFileSystem, artworkFileSystem http.FileSystem) http.Handler {
	var c httpauth.Checker = httpauth.None{}
	if auth {
		c = creds
	}

	w := httpauth.NewServeMux(c, http.NewServeMux())
	w.HandleFunc("/", rootHandler)
	w.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	w.Handle("/track/", http.StripPrefix("/track/", http.FileServer(mediaFileSystem)))
	w.Handle("/artwork/", http.StripPrefix("/artwork/", http.FileServer(artworkFileSystem)))
	w.Handle("/icon/", http.StripPrefix("/icon/", http.FileServer(store.FaviconFileSystem{artworkFileSystem})))
	w.Handle("/socket", l.WebsocketHandler())
	w.Handle("/api/ctrl/", http.StripPrefix("/api/ctrl/", ctrlHandler(l)))
	return w
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, "ui/tchaik.html")
}
