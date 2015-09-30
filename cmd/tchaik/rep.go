// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"

	"tchaik.com/index"
	"tchaik.com/index/attr"
)

// Group is a wrapper type for an index.Group which implements MarshalJSON
// for transmitting groups.
type Group struct {
	index.Group
	Key index.Key
}

// MarshalJSON implements json.Marshaler.
func (g *Group) MarshalJSON() ([]byte, error) {
	h := group{
		Name:        g.Name(),
		Key:         g.Key,
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
		Favourite:   g.Field("Favourite"),
		Checklist:   g.Field("Checklist"),
	}

	if c, ok := g.Group.(index.Collection); ok {
		return json.Marshal(buildCollection(h, c))
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, &Track{
			Track: t,
			group: g,
		})
	}
	return json.Marshal(h)
}

// stringSliceEqual is a function used to compare two interface{} types which are assumed
// to be of type []string (or interface{}(nil)).
func stringSliceEqual(x, y interface{}) bool {
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

type groupField struct {
	index.Group

	collection index.Collection
}

func (g groupField) Field(field string) interface{} {
	if field == "AlbumArtist" || field == "Artist" {
		if stringSliceEqual(g.Group.Field(field), g.collection.Field(field)) {
			return nil
		}
	}
	return g.Group.Field(field)
}

func buildCollection(h group, c index.Collection) group {
	for _, k := range c.Keys() {
		g := c.Get(k)
		g = index.FirstTrackAttr(attr.Strings("AlbumArtist"), g)
		g = index.CommonGroupAttr([]attr.Interface{attr.Strings("Artist")}, g)

		g = groupField{
			Group:      g,
			collection: c,
		}

		h.Groups = append(h.Groups, group{
			Name:        g.Name(),
			Key:         k,
			AlbumArtist: g.Field("AlbumArtist"),
			Artist:      g.Field("Artist"),
		})
	}
	return h
}

// Track is a simple wrapper for an index.Track.
type Track struct {
	index.Track

	group index.Group
}

// GetString implements index.Track.
func (t *Track) GetString(field string) string {
	// Always return the default value for these fields.
	if field == "ID" || field == "Name" {
		return t.Track.GetString(field)
	}

	if t.group.Field(field) != nil {
		return ""
	}
	return t.Track.GetString(field)
}

// GetStrings implements index.Track.
func (t *Track) GetStrings(field string) []string {
	if t.group.Field(field) != nil {
		return nil
	}
	return t.Track.GetStrings(field)
}

// GetInt implements index.Track.
func (t *Track) GetInt(field string) int {
	if field == "TotalTime" {
		return t.Track.GetInt(field)
	}

	if t.group.Field(field) != nil {
		return 0
	}
	return t.Track.GetInt(field)
}

func (t *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
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
	}{
		ID:          t.GetString("ID"),
		Name:        t.GetString("Name"),
		TotalTime:   t.GetInt("TotalTime"),
		Artist:      t.GetStrings("Artist"),
		AlbumArtist: t.GetStrings("AlbumArtist"),
		Composer:    t.GetStrings("Composer"),
		Album:       t.GetString("Album"),
		Kind:        t.GetString("Kind"),
		Year:        t.GetInt("Year"),
		DiscNumber:  t.GetInt("DiscNumber"),
		BitRate:     t.GetInt("BitRate"),
	})
}

type group struct {
	Name        string        `json:"name"`
	Key         index.Key     `json:"key"`
	TotalTime   interface{}   `json:"totalTime,omitempty"`
	Artist      interface{}   `json:"artist,omitempty"`
	AlbumArtist interface{}   `json:"albumArtist,omitempty"`
	Composer    interface{}   `json:"composer,omitempty"`
	BitRate     interface{}   `json:"bitRate,omitempty"`
	DiscNumber  interface{}   `json:"discNumber,omitempty"`
	ListStyle   interface{}   `json:"listStyle,omitempty"`
	ID          interface{}   `json:"id,omitempty"`
	Year        interface{}   `json:"year,omitempty"`
	Kind        interface{}   `json:"kind,omitempty"`
	Favourite   interface{}   `json:"favourite,omitempty"`
	Checklist   interface{}   `json:"checklist,omitempty"`
	Groups      []group       `json:"groups,omitempty"`
	Tracks      []index.Track `json:"tracks,omitempty"`
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
