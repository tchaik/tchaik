// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"tchaik.com/index"
	"tchaik.com/index/attr"
	"tchaik.com/index/favourite"
	"tchaik.com/store"
)

// Library is a type which encompases the components which form a full library.
type Library struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string][]index.FilterItem
	recent      []index.Path
	searcher    index.Searcher
	favourites  favourite.Store
}

type libraryFileSystem struct {
	store.FileSystem
	index.Library
}

// Open implements store.FileSystem and rewrites ID values to their corresponding Location
// values using the index.Library.
func (l *libraryFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	t, ok := l.Library.Track(strings.Trim(path, "/")) // IDs arrive with leading slash
	if !ok {
		return nil, fmt.Errorf("could not find track: %v", path)
	}

	loc := t.GetString("Location")
	if loc == "" {
		return nil, fmt.Errorf("invalid (empty) location for track: %v", path)
	}
	return l.FileSystem.Open(ctx, loc)
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
	ID          interface{} `json:",omitempty"`
	Year        interface{} `json:",omitempty"`
	Groups      []group     `json:",omitempty"`
	Tracks      []track     `json:",omitempty"`
	Favourite   bool        `json:",omitempty"`
}

type track struct {
	ID          string   `json:",omitempty"`
	Name        string   `json:",omitempty"`
	Album       string   `json:",omitempty"`
	Artist      []string `json:",omitempty"`
	AlbumArtist []string `json:",omitempty"`
	Composer    []string `json:",omitempty"`
	Year        int      `json:",omitempty"`
	DiscNumber  int      `json:",omitempty"`
	TotalTime   int      `json:",omitempty"`
	BitRate     int      `json:",omitempty"`
	Favourite   bool     `json:",omitempty"`
}

// StringSliceEqual is a function used to compare two interface{} types which are assumed
// to be of type []string (or interface{}(nil)).
func StringSliceEqual(x, y interface{}) bool {
	// Annoyingly we have to cater for zero values from map[string]interface{}
	// which don't have the correct type wrapping the nil.
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	xs := x.([]string) // NB: panics here are acceptable: should not be called on a non-'Strings' field.
	ys := y.([]string)
	if len(xs) != len(ys) {
		return false
	}
	for i, xss := range xs {
		if ys[i] != xss {
			return false
		}
	}
	return true
}

func buildCollection(h group, c index.Collection) group {
	getField := func(f string, g index.Group, c index.Collection) interface{} {
		if StringSliceEqual(g.Field(f), c.Field(f)) {
			return nil
		}
		return g.Field(f)
	}

	for _, k := range c.Keys() {
		g := c.Get(k)
		g = index.FirstTrackAttr(attr.Strings("AlbumArtist"), g)
		g = index.CommonGroupAttr([]attr.Interface{attr.Strings("Artist")}, g)
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
		ID:          g.Field("ID"),
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

	getStrings := func(t index.Track, field string) []string {
		if g.Field(field) != nil {
			return nil
		}
		return t.GetStrings(field)
	}

	getInt := func(t index.Track, field string) int {
		if g.Field(field) != nil {
			return 0
		}
		return t.GetInt(field)
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, track{
			ID:        t.GetString("ID"),
			Name:      t.GetString("Name"),
			TotalTime: t.GetInt("TotalTime"),
			// Potentially common fields (don't want to re-transmit everything)
			Artist:      getStrings(t, "Artist"),
			AlbumArtist: getStrings(t, "AlbumArtist"),
			Composer:    getStrings(t, "Composer"),
			Album:       getString(t, "Album"),
			Year:        getInt(t, "Year"),
			DiscNumber:  getInt(t, "DiscNumber"),
			BitRate:     getInt(t, "BitRate"),
		})
	}
	return h
}

// Fetch returns a group from the collection with the given path.
func (l *Library) Fetch(c index.Collection, path []string) (group, error) {
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
	g = index.Transform(g, index.SplitList("Artist", "AlbumArtist", "Composer"))
	c = index.Collect(g, index.ByPrefix("Name"))
	g = index.SubTransform(c, index.TrimEnumPrefix)
	g = index.SumGroupIntAttr("TotalTime", g)
	commonFields := []attr.Interface{
		attr.String("Album"),
		attr.Strings("Artist"),
		attr.Strings("AlbumArtist"),
		attr.Strings("Composer"),
		attr.Int("Year"),
		attr.Int("BitRate"),
		attr.Int("DiscNumber"),
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
	g = index.FirstTrackAttr(attr.String("ID"), g)

	return l.annotateFavourites(path, build(g, k)), nil
}

func (l *Library) annotateFavourites(path []string, g group) group {
	keyPath := make(index.Path, len(path)+1)
	keyPath[0] = "Root"
	for i, p := range path {
		keyPath[i+1] = index.Key(p)
	}

	if len(g.Tracks) > 0 || len(g.Groups) > 0 {
		g.Favourite = l.favourites.Get(keyPath)
	}
	return g
}

// FileSystem wraps the http.FileSystem in a library lookup which will translate /ID
// requests into their corresponding track paths.
func (l *Library) FileSystem(fs store.FileSystem) store.FileSystem {
	return store.Trace(&libraryFileSystem{fs, l.Library}, "libraryFileSystem")
}

// ExpandPaths constructs a collection (group) whose sub-groups are taken from the "Root"
// collection.
func (l *Library) ExpandPaths(paths []index.Path) group {
	return build(index.NewPathsCollection(l.collections["Root"], paths), index.Key("Root"))
}
