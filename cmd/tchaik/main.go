// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/websocket"

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
			Name:      "Media",
			FileSystem: mediaFileSystem,
		}
		artworkFileSystem = store.LogFileSystem{
			Name:      "Artwork",
			FileSystem: artworkFileSystem,
		}
	}

	mediaFileSystem = &libraryFileSystem{mediaFileSystem, l}
	artworkFileSystem = &libraryFileSystem{artworkFileSystem, l}

	libAPI := LibraryAPI{
		Library:  l,
		root:     root,
		searcher: searcher,
	}

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
	w.Handle("/socket", websocket.Handler(socketHandler(l)))
	return w
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, "ui/tchaik.html")
}


// Websocket handling
type socket struct {
	io.ReadWriter
	done chan struct{}
}

func (s *socket) Close() {
	select {
	case <-s.done:
		return
	default:
	}
	close(s.done)
}

type Command struct {
	Action string
	Input  string
	Path   []string
}

const (
	FetchAction  string = "FETCH"
	SearchAction string = "SEARCH"
)

func socketHandler(l LibraryAPI) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		s := socket{ws, make(chan struct{})}
		out, in := make(chan interface{}), make(chan *Command)
		errCh := make(chan error)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		// Encode messages from process and encode to the client
		enc := json.NewEncoder(s)
		go func() {
			defer wg.Done()
			defer s.Close()

			for {
				select {
				case x, ok := <-out:
					if !ok {
						return
					}

					if debug {
						b, err := json.MarshalIndent(x, "", "  ")
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(string(b))
					}

					if err := enc.Encode(x); err != nil {
						errCh <- fmt.Errorf("encode: %v", err)
						return
					}
				case <-s.done:
					return
				}
			}
		}()

		// Decode messages from the client and send them on the in channel
		go func() {
			defer wg.Done()
			defer s.Close()

			dec := json.NewDecoder(s)
			for {
				c := &Command{}
				if err := dec.Decode(c); err != nil {
					if err == io.EOF && debug {
						fmt.Println("websocket closed")
						return
					}
					errCh <- err
					return
				}
				in <- c
			}
		}()

		go func() {
			for x := range in {
				if debug {
					fmt.Printf("command received: %#v\n", x)
				}
				switch x.Action {
				case FetchAction:
					handleCollectionList(l, x, out)
				case SearchAction:
					handleSearch(l, x, out)
				default:
					fmt.Printf("unknown command: %v", x.Action)
				}
			}
		}()

		go func() {
			wg.Wait()

			close(in)
			close(out)
			close(errCh)
		}()

		go func() {
			for err := range errCh {
				if err != nil {
					fmt.Printf("websocket handler: %v\n", err)
				}
			}
		}()

		select {}
	}
}

func handleCollectionList(l LibraryAPI, x *Command, out chan<- interface{}) {
	if len(x.Path) < 1 {
		fmt.Printf("invalid path: %v\n", x.Path)
		return
	}

	g, err := l.Fetch(l.root, x.Path[1:])
	if err != nil {
		fmt.Printf("error in Fetch: %v (path: %#v)", err, x.Path[1:])
		return
	}

	o := struct {
		Action string
		Data   interface{}
	}{
		x.Action,
		struct {
			Path []string
			Item group
		}{
			x.Path,
			g,
		},
	}
	out <- o
}

func handleSearch(l LibraryAPI, x *Command, out chan<- interface{}) {
	paths := l.searcher.Search(x.Input)
	o := struct {
		Action string
		Data   interface{}
	}{
		Action: x.Action,
		Data:   paths,
	}
	out <- o
}
