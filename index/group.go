// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package group defines the concept of a group which is an ordered list of
// playable items
package index

import "fmt"

// Key represents a unique value used to represent a group within a collection.
type Key string

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}

// Tracker is an interface which defines the Tracks method which returns a list
// ot Tracks.
type Tracker interface {
	Tracks() []Track
}

// Group is an interface which represents a named group of Tracks.
type Group interface {
	Tracker

	// Name returns the name of the group.
	Name() string

	// Field returns the value of a field.
	Field(string) interface{}
}

// Collection represents a reorganisation of Tracks into a series of Groups.
type Collection interface {
	Group

	// Keys returns a slice of Keys which give an ordering for Groups in the Collection.
	Keys() []Key

	// Get returns the Group corresponding to the given Key.
	Get(Key) Group
}

// Collect applies the Collector to the Tracker and returns the resulting Collection.
func Collect(t Tracker, c Collector) Collection {
	return c.Collect(t)
}

// col is a basic implementation of Collection. It assumes that all Groups have unique names
// and so uses Group names for the keys.
type col struct {
	keys []Key
	name string
	grps map[Key]Group
	flds map[string]interface{}
}

func newCol(name string) col {
	return col{
		name: name,
		grps: make(map[Key]Group),
		flds: make(map[string]interface{}),
	}
}

func (c col) Keys() []Key     { return c.keys }
func (c col) Name() string    { return c.name }
func (c col) Get(k Key) Group { return c.grps[k] }

func (c col) Field(field string) interface{} { return c.flds[field] }

func (c col) Tracks() []Track {
	return collectionTracks(c)
}

// addTrack adds the track t to the collection by adding it to a group with key k.  If no such
// group exists in the collection, then a new group is created with name n.
func (c *col) addTrack(n string, k Key, t Track) {
	if _, ok := c.grps[k]; !ok {
		ng := group{name: n, tracks: make([]Track, 1)}
		ng.tracks[0] = t
		c.grps[k] = ng
		c.keys = append(c.keys, k)
		return
	}
	ng := c.grps[k]
	ngg := ng.(group)
	ngg.tracks = append(ngg.tracks, t)
	c.grps[k] = ngg
}

// add adds the track t to the collection, using the name n as the key.
func (c *col) add(n string, t Track) {
	c.addTrack(n, Key(n), t)
}

// collectionTracks iterates over all the tracks in all the groups of the collection to construct a
// slice of Tracks.
func collectionTracks(c Collection) []Track {
	keys := c.Keys()
	var tracks []Track
	for _, k := range keys {
		tracks = append(tracks, c.Get(k).Tracks()...)
	}
	return tracks
}

// Collector is an interface which defines the Collect method.
type Collector interface {
	Collect(Tracker) Collection
}

// group is a basic implementation of Group
type group struct {
	name   string
	tracks []Track
	fields map[string]interface{}
}

// Name implements Group.
func (g group) Name() string { return g.name }

// Tracks implements Tracker.
func (g group) Tracks() []Track { return g.tracks }

// Field implements Group.
func (g group) Field(field string) interface{} { return g.fields[field] }

// subCol is a wrapper around an existing collection which overrides the Get
// method.
type subCol struct {
	Collection
	grps map[Key]Group
	flds map[string]interface{}
}

// Get implements Collection.
func (sg subCol) Get(k Key) Group {
	return sg.grps[k]
}

// Field implements Group.
func (sg subCol) Field(f string) interface{} {
	if x, ok := sg.flds[f]; ok {
		return x
	}
	return sg.Collection.Field(f)
}

// SubCollect applies the given Collector to each of the "leaf-Groups"
// in the Collection.
func SubCollect(c Collection, r Collector) Collection {
	keys := c.Keys()
	nc := subCol{
		Collection: c,
		grps:       make(map[Key]Group, len(keys)),
		flds:       make(map[string]interface{}),
	}
	for _, k := range keys {
		g := c.Get(k)
		if gc, ok := g.(Collection); ok {
			nc.grps[k] = SubCollect(gc, r)
			continue
		}
		nc.grps[k] = r.Collect(g)
	}
	return nc
}

// WalkFn is the type of the function called for each Track visited by Walk.
type WalkFn func(Track, []string)

func walkCollection(c Collection, p []string, f WalkFn) {
	for _, k := range c.Keys() {
		g := c.Get(k)
		np := make([]string, len(p)+1)
		copy(np, p)
		np[len(p)] = string(k)
		Walk(g, np, f)
	}
}

// Walk transverses the Group g and calls the WalkFn f on each Track.
func Walk(g Group, path []string, f WalkFn) {
	if gc, ok := g.(Collection); ok {
		walkCollection(gc, path, f)
		return
	}
	for _, t := range g.Tracks() {
		f(t, path)
	}
}
// ByAttr is a type which implements Collector, and groups elements by the value of
// the attribute given by the underlying Attr instance.
type ByAttr Attr

func (a ByAttr) Collect(tracker Tracker) Collection {
	ga := Attr(a)
	name := "by " + ga.Field()
	if tg, ok := tracker.(Group); ok {
		name = tg.Name()
	}
	gg := newCol(name)
	for _, t := range tracker.Tracks() {
		gg.add(fmt.Sprintf("%v", ga.fn(t)), t)
	}
	return gg
}
