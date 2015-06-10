// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/tchaik/tchaik/index"
)

type LibraryAPI struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string][]index.FilterItem
	recent      []index.Path
	searcher    index.Searcher
	players     *players
}

type libraryFileSystem struct {
	http.FileSystem
	index.Library
}

// Open implements http.FileSystem and rewrites TrackID values to their corresponding Location
// values using the index.Library
func (l *libraryFileSystem) Open(path string) (http.File, error) {
	t, ok := l.Library.Track(strings.Trim(path, "/")) // TrackIDs arrive with leading slash
	if !ok {
		return nil, fmt.Errorf("could not find track: %v", path)
	}

	loc := t.GetString("Location")
	if loc == "" {
		return nil, fmt.Errorf("invalid (empty) location for track: %v", path)
	}
	return l.FileSystem.Open(loc)
}

type group struct {
	Name        string
	Key         index.Key
	TotalTime   interface{} `json:",omitempty"`
	Artist      interface{} `json:",omitempty"`
	AlbumArtist interface{} `json:",omitempty"`
	Composer    interface{} `json:",omitempty"`
	BitRate     interface{} `json:",omitempty"`
	DiscNumber  interface{} `json:",omitempty"`
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
	BitRate     int    `json:",omitempty"`
}

func buildCollection(h group, c index.Collection) group {
	for _, k := range c.Keys() {
		g := c.Get(k)
		g = index.FirstTrackAttr(index.StringAttr("AlbumArtist"), g)
		g = index.CommonGroupAttr([]index.Attr{index.StringAttr("Artist")}, g)
		h.Groups = append(h.Groups, group{
			Name:        g.Name(),
			Key:         k,
			AlbumArtist: g.Field("AlbumArtist"),
			Artist:      g.Field("Artist"),
		})
	}
	return h
}

func build(g index.Group, key index.Key) group {
	h := group{
		Name:        g.Name(),
		Key:         key,
		TotalTime:   g.Field("TotalTime"),
		Artist:      g.Field("Artist"),
		AlbumArtist: g.Field("AlbumArtist"),
		Composer:    g.Field("Composer"),
		Year:        g.Field("Year"),
		BitRate:     g.Field("BitRate"),
		DiscNumber:  g.Field("DiscNumber"),
		ListStyle:   g.Field("ListStyle"),
		TrackID:     g.Field("TrackID"),
	}

	if c, ok := g.(index.Collection); ok {
		return buildCollection(h, c)
	}

	getString := func(t index.Track, field string) string {
		if g.Field(field) != "" && g.Field(field) == t.GetString(field) {
			return ""
		}
		return t.GetString(field)
	}

	getInt := func(t index.Track, field string) int {
		if g.Field(field) != 0 && g.Field(field) == t.GetInt(field) {
			return 0
		}
		return t.GetInt(field)
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, track{
			TrackID:   t.GetString("TrackID"),
			Name:      t.GetString("Name"),
			TotalTime: t.GetInt("TotalTime"),
			// Potentially common fields (don't want to re-transmit everything)
			Album:       getString(t, "Album"),
			Artist:      getString(t, "Artist"),
			AlbumArtist: getString(t, "AlbumArtist"),
			Composer:    getString(t, "Composer"),
			Year:        getInt(t, "Year"),
			DiscNumber:  getInt(t, "DiscNumber"),
			BitRate:     getInt(t, "BitRate"),
		})
	}
	return h
}

func (l *LibraryAPI) Fetch(c index.Collection, path []string) (group, error) {
	if len(path) == 0 {
		return build(c, index.Key("Root")), nil
	}

	var g index.Group = c
	k := index.Key(path[0])
	g = c.Get(k)

	if g == nil {
		return group{}, fmt.Errorf("invalid path: near '%v'", path[0])
	}

	index.Sort(g.Tracks(), index.MultiSort(index.SortByInt("DiscNumber"), index.SortByInt("TrackNumber")))
	c = index.Collect(g, index.ByPrefix("Name"))
	g = index.SubTransform(c, index.TrimEnumPrefix)
	g = index.SumGroupIntAttr("TotalTime", g)
	commonFields := []index.Attr{
		index.StringAttr("Album"),
		index.StringAttr("Artist"),
		index.StringAttr("AlbumArtist"),
		index.StringAttr("Composer"),
		index.IntAttr("Year"),
		index.IntAttr("BitRate"),
		index.IntAttr("DiscNumber"),
	}
	g = index.CommonGroupAttr(commonFields, g)
	g = index.RemoveEmptyCollections(g)

	for i, p := range path[1:] {
		var ok bool
		c, ok = g.(index.Collection)
		if !ok {
			return group{}, fmt.Errorf("retrieved Group is not a Collection")
		}
		k = index.Key(p)
		g = c.Get(k)

		if g == nil {
			return group{}, fmt.Errorf("invalid path near '%v'", path[1:][i])
		}

		if _, ok = g.(index.Collection); !ok {
			if i == len(path[1:])-1 {
				break
			}
			return group{}, fmt.Errorf("retrieved Group isn't a Collection: %v", p)
		}
	}
	if g == nil {
		return group{}, fmt.Errorf("could not find group")
	}
	g = index.FirstTrackAttr(index.StringAttr("TrackID"), g)

	return build(g, k), nil
}

func (l *LibraryAPI) FileSystem(fs http.FileSystem) http.FileSystem {
	return &libraryFileSystem{fs, l.Library}
}

type players struct {
	sync.RWMutex
	m map[string]Player
}

func newPlayers() *players {
	return &players{m: make(map[string]Player)}
}

func (s *players) add(p Player) {
	s.Lock()
	defer s.Unlock()

	s.m[p.Key()] = p
}

func (s *players) remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)
}

func (s *players) get(key string) Player {
	s.RLock()
	defer s.RUnlock()

	return s.m[key]
}

func (s *players) list() []string {
	s.RLock()
	defer s.RUnlock()

	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *players) MarshalJSON() ([]byte, error) {
	keys := s.list()
	return json.Marshal(struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	})
}

func playersGet(l LibraryAPI, w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(l.players)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func createPlayer(l LibraryAPI, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	postData := struct {
		Key        string
		PlayerKeys []string
	}{}
	err := dec.Decode(&postData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	if p := l.players.get(postData.Key); p != nil {
		http.Error(w, "player key already exists", http.StatusBadRequest)
		return
	}

	if postData.PlayerKeys == nil || len(postData.PlayerKeys) == 0 {
		http.Error(w, "no player keys specified", http.StatusBadRequest)
		return
	}

	var players []Player
	for _, pk := range postData.PlayerKeys {
		p := l.players.get(pk)
		if p == nil {
			http.Error(w, fmt.Sprintf("invalid player key: %v", pk), http.StatusBadRequest)
			return
		}
		players = append(players, p)
	}
	l.players.add(MultiPlayer(postData.Key, players...))
	w.WriteHeader(http.StatusCreated)
}

func playerAction(p Player, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	putData := struct {
		Action string
		Value  interface{}
	}{}
	err := dec.Decode(&putData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	switch putData.Action {
	case "play":
		err = p.Play()

	case "pause":
		err = p.Pause()

	case "next":
		err = p.NextTrack()

	case "prev":
		err = p.PreviousTrack()

	case "togglePlayPause":
		err = p.TogglePlayPause()

	case "toggleMute":
		err = p.ToggleMute()

	case "setVolume":
		f, ok := putData.Value.(float64)
		if !ok {
			err = InvalidValueError("invalid volume value: expected float")
			break
		}
		err = p.SetVolume(f)

	case "setMute":
		b, ok := putData.Value.(bool)
		if !ok {
			err = InvalidValueError("invalid mute value: expected boolean")
			break
		}
		err = p.SetMute(b)

	case "setTime":
		f, ok := putData.Value.(float64)
		if !ok {
			err = InvalidValueError("invalid time value: expected float")
			break
		}
		err = p.SetTime(f)

	default:
		err = InvalidValueError("invalid action")
		return
	}

	if err != nil {
		if err, ok := err.(InvalidValueError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("error sending player command: %v", err), http.StatusInternalServerError)
	}
}

func playerView(p Player, w http.ResponseWriter, t *http.Request) {
	enc := json.NewEncoder(w)
	err := enc.Encode(p)
	if err != nil {
		log.Printf("error encoding player data: %v", err)
		return
	}
}

func playersHandler(l LibraryAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			switch r.Method {
			case "GET":
				playersGet(l, w, r)
			case "POST":
				createPlayer(l, w, r)
			}
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		if len(paths) != 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		p := l.players.get(paths[0])
		if p == nil {
			http.Error(w, "invalid player key", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "DELETE":
			l.players.remove(paths[0])
			w.WriteHeader(http.StatusNoContent)
			return

		case "PUT":
			playerAction(p, w, r)

		case "GET":
			playerView(p, w, r)
		}
	})
}
