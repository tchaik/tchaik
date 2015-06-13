// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tchaik/tchaik/index"
)

type LibraryAPI struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string][]index.FilterItem
	recent      []index.Path
	searcher    index.Searcher
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
	getField := func(f string, g index.Group, c index.Collection) interface{} {
		if g.Field(f) == c.Field(f) {
			return nil
		}
		return g.Field(f)
	}

	for _, k := range c.Keys() {
		g := c.Get(k)
		g = index.FirstTrackAttr(index.StringAttr("AlbumArtist"), g)
		g = index.CommonGroupAttr([]index.Attr{index.StringAttr("Artist")}, g)
		h.Groups = append(h.Groups, group{
			Name:        g.Name(),
			Key:         k,
			AlbumArtist: getField("AlbumArtist", g, c),
			Artist:      getField("Artist", g, c),
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
		if g.Field(field) != nil {
			return ""
		}
		return t.GetString(field)
	}

	getInt := func(t index.Track, field string) int {
		if g.Field(field) != nil {
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

func (l *LibraryAPI) ExpandPaths(paths []index.Path) group {
	return build(index.NewPathsCollection(l.collections["Root"], paths), index.Key("Root"))
}
