// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"tchaik.com/index"
	"tchaik.com/index/attr"

	"tchaik.com/store"
)

// Library is a type which encompases the components which form a full library.
type Library struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string][]index.FilterItem
	recent      []index.Path
	searcher    index.Searcher
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
	loc = filepath.ToSlash(loc)
	return l.FileSystem.Open(ctx, loc)
}

type group struct {
	Name        string      `json:"name"`
	Key         index.Key   `json:"key"`
	TotalTime   interface{} `json:"totalTime,omitempty"`
	Artist      interface{} `json:"artist,omitempty"`
	AlbumArtist interface{} `json:"albumArtist,omitempty"`
	Composer    interface{} `json:"composer,omitempty"`
	BitRate     interface{} `json:"bitRate,omitempty"`
	DiscNumber  interface{} `json:"discNumber,omitempty"`
	ListStyle   interface{} `json:"listStyle,omitempty"`
	ID          interface{} `json:"id,omitempty"`
	Year        interface{} `json:"year,omitempty"`
	Kind        interface{} `json:"kind,omitempty"`
	Groups      []group     `json:"groups,omitempty"`
	Tracks      []track     `json:"tracks,omitempty"`
	Favourite   bool        `json:"favourite,omitempty"`
	Checklist   bool        `json:"checklist,omitempty"`
}

type track struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Album       string   `json:"album,omitempty"`
	Artist      []string `json:"artist,omitempty"`
	AlbumArtist []string `json:"albumArtist,omitempty"`
	Composer    []string `json:"composer,omitempty"`
	Kind        string   `json:"kind,omitempty"`
	Year        int      `json:"year,omitempty"`
	DiscNumber  int      `json:"discNumber,omitempty"`
	TotalTime   int      `json:"totalTime,omitempty"`
	BitRate     int      `json:"bitRate,omitempty"`
	Favourite   bool     `json:"favourite,omitempty"`
	Checklist   bool     `json:"checklist,omitempty"`
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
		Kind:        g.Field("Kind"),
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
			Kind:        getString(t, "Kind"),
			Year:        getInt(t, "Year"),
			DiscNumber:  getInt(t, "DiscNumber"),
			BitRate:     getInt(t, "BitRate"),
		})
	}
	return h
}

type rootCollection struct {
	index.Collection
}

func (r *rootCollection) Get(k index.Key) index.Group {
	g := r.Collection.Get(k)
	if g == nil {
		return g
	}

	index.Sort(g.Tracks(), index.MultiSort(index.SortByString("Kind"), index.SortByInt("DiscNumber"), index.SortByInt("TrackNumber")))
	g = index.Transform(g, index.SplitList("Artist", "AlbumArtist", "Composer"))
	g = index.Transform(g, index.TrimTrackNumPrefix)
	c := index.Collect(g, index.ByPrefix("Name"))
	g = index.SubTransform(c, index.TrimEnumPrefix)
	g = index.SumGroupIntAttr("TotalTime", g)
	commonFields := []attr.Interface{
		attr.String("Album"),
		attr.Strings("Artist"),
		attr.Strings("AlbumArtist"),
		attr.Strings("Composer"),
		attr.String("Kind"),
		attr.Int("Year"),
		attr.Int("BitRate"),
		attr.Int("DiscNumber"),
	}
	g = index.CommonGroupAttr(commonFields, g)
	g = index.RemoveEmptyCollections(g)
	return g
}

func (l *Library) Build(c index.Collection, p index.Path) (index.Group, error) {
	if len(p) == 0 {
		return c, nil
	}

	g, err := index.GroupFromPath(&rootCollection{c}, p)
	if err != nil {
		return nil, err
	}
	g = index.FirstTrackAttr(attr.String("ID"), g)
	return g, nil
}

// Fetch returns a group from the collection with the given path.
func (l *Library) Fetch(c index.Collection, p index.Path) (group, error) {
	if len(p) == 0 {
		return build(c, index.Key("Root")), nil
	}

	k := index.Key(p[0])
	g, err := l.Build(c, p)
	if err != nil {
		return group{}, err
	}
	return build(g, k), nil
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
