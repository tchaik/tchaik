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

type Track struct {
	index.Track

	group index.Group
}

func (t *Track) getString(field string) string {
	// Always return the default value for these fields.
	if field == "ID" || field == "Name" {
		return t.Track.GetString(field)
	}

	if t.group.Field(field) != nil {
		return ""
	}
	return t.Track.GetString(field)
}

func (t *Track) getStrings(field string) []string {
	if t.group.Field(field) != nil {
		return nil
	}
	return t.Track.GetStrings(field)
}

func (t *Track) getInt(field string) int {
	if field == "TotalTime" {
		return t.Track.GetInt(field)
	}

	if t.group.Field(field) != nil {
		return 0
	}
	return t.Track.GetInt(field)
}

func (t *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(track{
		ID:          t.GetString("ID"),
		Name:        t.GetString("Name"),
		TotalTime:   t.GetInt("TotalTime"),
		Artist:      t.getStrings("Artist"),
		AlbumArtist: t.getStrings("AlbumArtist"),
		Composer:    t.getStrings("Composer"),
		Album:       t.getString("Album"),
		Kind:        t.getString("Kind"),
		Year:        t.getInt("Year"),
		DiscNumber:  t.getInt("DiscNumber"),
		BitRate:     t.getInt("BitRate"),
	})
}

func (g *Group) build() group {
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
		return buildCollection(h, c)
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, &Track{
			Track: t,
			group: g,
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

func (g *Group) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.build())
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
}
