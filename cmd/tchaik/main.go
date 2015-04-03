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
	"net/http/httputil"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/dhowden/httpauth"
	"github.com/dhowden/itl"

	"github.com/dhowden/tchaik/index"
	"github.com/dhowden/tchaik/store"
	"github.com/dhowden/tchaik/store/cmdflag"
)

var debug bool
var xml, tchJSON string

var listenAddr string
var certFile, keyFile string

var auth bool

var trimLocationPrefix, addLocationPrefix string

func init() {
	flag.BoolVar(&debug, "debug", false, "print debugging information")

	flag.StringVar(&listenAddr, "http", "localhost:8080", "bind address to http listen")
	flag.StringVar(&certFile, "cert", "", "path to an SSL certificate file.  Must also specify -key")
	flag.StringVar(&keyFile, "key", "", "path to an SSL certificate key file.  Must also specify -cert")

	flag.StringVar(&xml, "xml", "", "path to iTunes Library XML file")
	flag.StringVar(&tchJSON, "tchJSON", "", "path to Tchaik Library file")

	flag.BoolVar(&auth, "auth", false, "use basic HTTP authentication")

	flag.StringVar(&trimLocationPrefix, "trim", "", "trim the location prefix by the given string")
	flag.StringVar(&addLocationPrefix, "prefix", "", "add the given prefix to location")
}

var creds = httpauth.Creds(map[string]string{
	"user": "password",
})

func readLibrary(xml, tchJSON string) (index.Library, error) {
	if xml == "" && tchJSON == "" {
		return nil, fmt.Errorf("must specify at least one library file (xml or tchJSON)")
	}

	if xml != "" && tchJSON != "" {
		return nil, fmt.Errorf("must only specify one library file")
	}

	var l index.Library
	if xml != "" {
		f, err := os.Open(xml)
		if err != nil {
			return nil, fmt.Errorf("could not open iTunes library file: %v", err)
		}

		fmt.Printf("Parsing %v...\n", xml)
		it, err := itl.ReadFromXML(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing iTunes library file: %v", err)
		}
		f.Close()

		fmt.Println("Building Tchaik Library...")
		l = index.Convert(index.NewITunesLibrary(&it), "TrackID")
		return l, nil
	}

	f, err := os.Open(tchJSON)
	if err != nil {
		return nil, fmt.Errorf("could not open Tchaik library file: %v", err)
	}

	fmt.Printf("Parsing %v...", tchJSON)
	l, err = index.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing Tchaik library file: %v\n", err)
	}
	return l, nil
}

func main() {
	flag.Parse()
	l, err := readLibrary(xml, tchJSON)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Grouping into albums...")
	albums := index.Collect(l, index.ByAttr(index.StringAttr("Album")))
	fmt.Println("done.")
	fmt.Printf(" - sorting albums...")
	index.SortKeysByGroupName(albums)
	fmt.Println("done.")

	fmt.Printf("Building search index...")
	wi := index.BuildWordIndex(albums, []string{"Composer", "Artist", "Album", "Name"})
	s := index.FlatSearcher{index.WordsIntersectSearcher(index.BuildPrefixExpandSearcher(wi, wi, 10))}
	fmt.Println("done.")

	fileSystem, artworkFileSystem, err := cmdflag.Stores()
	if err != nil {
		fmt.Println("error setting up stores:", err)
		os.Exit(1)
	}

	libAPI := LibraryAPI{
		Library:        l,
		trackHandler:   http.FileServer(store.LogFileSystem{"fileserver", fileSystem}),
		artworkHandler: http.FileServer(artworkFileSystem),
	}

	var c httpauth.Checker = httpauth.None{}
	if auth {
		c = creds
	}

	httpauth.HandleFunc(c, "/", rootHandler)
	httpauth.Handle(c, "/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	httpauth.HandleFunc(c, "/track/", libAPI.TrackHandler)
	httpauth.HandleFunc(c, "/artwork/", libAPI.ArtworkHandler)
	httpauth.Handle(c, "/socket", websocket.Handler(socketHandler(libAPI, albums, s)))

	if certFile != "" && keyFile != "" {
		fmt.Printf("Web server is running on https://%v\n", listenAddr)
		fmt.Println("Quit the server with CTRL-C.")

		log.Fatal(http.ListenAndServeTLS(listenAddr, certFile, keyFile, nil))
	}

	fmt.Printf("Web server is running on http://%v\n", listenAddr)
	fmt.Println("Quit the server with CTRL-C.")

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, "ui/tchaik.html")
}

func rewriteLocation(l string) string {
	l = strings.TrimPrefix(l, trimLocationPrefix)
	l = addLocationPrefix + l
	return l
}

func debugDumpRequest(r *http.Request) {
	if debug {
		rb, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println("could not dump request:", err)
		}
		fmt.Println(string(rb))
	}
}

type LibraryAPI struct {
	index.Library

	trackHandler   http.Handler
	artworkHandler http.Handler
}

func (l *LibraryAPI) locationForRequest(r *http.Request, prefix string) (string, error) {
	id := strings.TrimPrefix(r.URL.Path, prefix)
	t, ok := l.Track(id)
	if !ok {
		return "", fmt.Errorf("could not find track: %v\n", id)
	}

	loc := t.GetString("Location")
	if loc == "" {
		return "", fmt.Errorf("invalid (empty) location for track: %v", id)
	}
	return rewriteLocation(loc), nil
}

func (l *LibraryAPI) TrackHandler(w http.ResponseWriter, r *http.Request) {
	debugDumpRequest(r)

	loc, err := l.locationForRequest(r, "/track/")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	r.URL.Path = loc
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	l.trackHandler.ServeHTTP(w, r)
}

func (l *LibraryAPI) ArtworkHandler(w http.ResponseWriter, r *http.Request) {
	debugDumpRequest(r)

	loc, err := l.locationForRequest(r, "/artwork/")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	r.URL.Path = loc
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	l.artworkHandler.ServeHTTP(w, r)
}

type group struct {
	Name        string
	Key         index.Key
	TotalTime   interface{} `json:",omitempty"`
	Artist      interface{} `json:",omitempty"`
	AlbumArtist interface{} `json:",omitempty"`
	Composer    interface{} `json:",omitempty"`
	ListStyle   interface{} `json:",omitempty"`
	TrackID     interface{} `json:",omitempty"`
	Year        interface{} `json:",omitempty"`
	Groups      []group     `json:",omitempty"`
	Tracks      []track     `json:",omitempty"`
}

type track struct {
	TrackID     string `json:",omitempty"`
	Name        string `json:",omitempty"`
	Album       string `json:",omitempty"`
	Artist      string `json:",omitempty"`
	AlbumArtist string `json:",omitempty"`
	Composer    string `json:",omitempty"`
	Year        int    `json:",omitempty"`
	DiscNumber  int    `json:",omitempty"`
	TotalTime   int    `json:",omitempty"`
}

func (l *LibraryAPI) build(g index.Group, key index.Key) group {
	h := group{
		Name:        g.Name(),
		Key:         key,
		TotalTime:   g.Field("TotalTime"),
		Artist:      g.Field("Artist"),
		AlbumArtist: g.Field("AlbumArtist"),
		Composer:    g.Field("Composer"),
		Year:        g.Field("Year"),
		ListStyle:   g.Field("ListStyle"),
		TrackID:     g.Field("TrackID"),
	}

	if c, ok := g.(index.Collection); ok {
		for _, k := range c.Keys() {
			sg := c.Get(k)
			sg = index.FirstTrackAttr(index.StringAttr("TrackID"), sg)
			h.Groups = append(h.Groups, group{
				Name:    sg.Name(),
				Key:     k,
				TrackID: sg.Field("TrackID"),
			})
		}
		return h
	}

	getString := func(t index.Track, field string) string {
		if g.Field(field) != "" {
			return ""
		}
		return t.GetString(field)
	}

	getInt := func(t index.Track, field string) int {
		if g.Field(field) != 0 {
			return 0
		}
		return t.GetInt(field)
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, track{
			TrackID:    t.GetString("TrackID"),
			Name:       t.GetString("Name"),
			TotalTime:  t.GetInt("TotalTime"),
			DiscNumber: t.GetInt("DiscNumber"),
			// Potentially common fields (don't want to re-transmit everything)
			Album:       getString(t, "Album"),
			Artist:      getString(t, "Artist"),
			AlbumArtist: getString(t, "AlbumArtist"),
			Composer:    getString(t, "Composer"),
			Year:        getInt(t, "Year"),
		})
	}
	return h
}

func (l *LibraryAPI) Fetch(c index.Collection, path []string) group {
	var k index.Key
	var g index.Group
	g = c

	if debug {
		fmt.Printf("buildFromPath: %#v\n", path)
	}

	if len(path) > 0 {
		k = index.Key(path[0])
		g = c.Get(k)

		if g == nil {
			fmt.Printf("error: invalid path: '%#v' near '%v'\n", path, path[0])
			return group{}
		}

		index.Sort(g.Tracks(), index.MultiSort(index.SortByInt("DiscNumber"), index.SortByInt("TrackNumber")))
		c = index.ByPrefix("Name").Collect(g)
		c = index.SubTransform(c, index.TrimEnumPrefix)
		g = c
		g = index.SumGroupIntAttr("TotalTime", g)
		commonFields := []index.Attr{
			index.StringAttr("Album"),
			index.StringAttr("Artist"),
			index.StringAttr("AlbumArtist"),
			index.StringAttr("Composer"),
			index.IntAttr("Year"),
		}
		g = index.CommonGroupAttr(commonFields, g)
		g = index.RemoveEmptyCollections(g)

		for i, p := range path[1:] {
			var ok bool
			c, ok = g.(index.Collection)
			if !ok {
				fmt.Printf("error: group for path '%#v' is not a Collection\n", path)
				return group{}
			}
			k = index.Key(p)
			g = c.Get(k)

			if g == nil {
				fmt.Printf("error: invalid path: '%#v' near '%v'\n", path, path[1:][i])
				return group{}
			}

			c, ok = g.(index.Collection)
			if !ok {
				if i == len(path[1:])-1 {
					break
				}
				fmt.Printf("error: retrieved group isn't a collection: %v (path component: %v)\n", path, p)
				return group{}
			}
		}
		if g == nil {
			fmt.Println("error: could not find group")
			return group{}
		}
		g = index.FirstTrackAttr(index.StringAttr("TrackID"), g)
	} else {
		k = index.Key("Root")
	}
	return l.build(g, k)
}

// Websocket handling
type socket struct {
	io.ReadWriter
	done <-chan struct{}
}

type Command struct {
	Action string
	Input  string
	Path   []string
}

type Response struct {
	Tag  string
	Data interface{}
}

const (
	FetchAction  string = "FETCH"
	SearchAction string = "SEARCH"
)

func socketHandler(l LibraryAPI, collection index.Collection, searcher index.Searcher) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		s := socket{ws, make(chan struct{})}
		out, in := make(chan interface{}), make(chan *Command)
		errc := make(chan error, 1)

		wg := &sync.WaitGroup{}
		wg.Add(3)

		// Encode messages from process and encode to the client
		enc := json.NewEncoder(s)
		go func() {
			defer wg.Done()
			for x := range out {
				if debug {
					b, err := json.MarshalIndent(x, "", "  ")
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(string(b))
				}

				if err := enc.Encode(x); err != nil {
					fmt.Printf("error sending %v: %v\n", x, err)
					errc <- err
					return
				}
			}
		}()

		// Decode messages from the client and send them on the in channel
		go func() {
			defer wg.Done()
			dec := json.NewDecoder(s)
			for {
				c := &Command{}
				if err := dec.Decode(c); err != nil {
					fmt.Println("decode:", err)
					errc <- err
					return
				}
				in <- c
			}
		}()

		go func() {
			defer wg.Done()
			for x := range in {
				if debug {
					fmt.Printf("Command Received: %#v\n", x)
				}
				switch x.Action {
				case FetchAction:
					handleCollectionList(l, collection, x, out)
				case SearchAction:
					handleSearch(searcher, x, out)
				default:
					fmt.Printf("unknown command: %v", x.Action)
				}
			}
		}()

		select {}
	}
}

func handleCollectionList(l LibraryAPI, c index.Collection, x *Command, out chan<- interface{}) {
	if len(x.Path) < 1 {
		fmt.Printf("invalid path: %v\n", x.Path)
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
			l.Fetch(c, x.Path[1:]),
		},
	}
	out <- o
}

func handleSearch(s index.Searcher, x *Command, out chan<- interface{}) {
	paths := s.Search(x.Input)
	o := struct {
		Action string
		Data   interface{}
	}{
		Action: x.Action,
		Data:   paths,
	}
	out <- o
}
