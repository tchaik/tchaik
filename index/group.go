// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package group defines the concept of a group which is an ordered list of
// playable items
package index

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

// Key represents a unique value used to represent a group within a collection.
type Key string

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}

// Tracker is an interface which defines the Tracks method which returns a list
// of Tracks.
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

// Collection is an interface which represents an ordered series of Groups.
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

// NewCollection creates a new collection from a source collection `c` which will have the groups
// represented by the given list of paths.  All the paths are assumed to be unique, an of at least
// length 2.
func NewPathsCollection(src Collection, paths []Path) Collection {
	keys := make([]Key, len(paths))
	for i, path := range paths {
		keys[i] = path[1]
	}

	return pathsCollection{
		Collection: src,
		name:       "paths collection",
		keys:       keys,
	}
}

type pathsCollection struct {
	Collection
	name string
	keys []Key
}

func (c pathsCollection) Keys() []Key  { return c.keys }
func (c pathsCollection) Name() string { return c.name }

func (c pathsCollection) Field(string) interface{} { return nil }

func (c pathsCollection) Tracks() []Track {
	// TODO: Do something better here - this method shouldn't really get called
	return nil
}

// col is a basic implementation of Collection. It assumes that all Groups have unique names
// and so uses Group names for the keys.
type col struct {
	keys []Key
	name string
	grps map[Key]group
	flds map[string]interface{}
}

func newCol(name string) col {
	return col{
		name: name,
		grps: make(map[Key]group),
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
	if g, ok := c.grps[k]; ok {
		g.tracks = append(g.tracks, t)
		c.grps[k] = g
		return
	}
	g := group{
		name:   n,
		tracks: make([]Track, 1),
	}
	g.tracks[0] = t
	c.grps[k] = g
	c.keys = append(c.keys, k)
}

// add adds the track t to the collection, using the name n as the key.
func (c *col) add(n string, t Track) {
	k := Key(fmt.Sprintf("%x", sha1.Sum([]byte(n)))[:6])
	c.addTrack(n, k, t)
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
type WalkFn func(Track, Path)

func walkCollection(c Collection, p Path, f WalkFn) {
	for _, k := range c.Keys() {
		g := c.Get(k)
		np := make(Path, len(p)+1)
		copy(np, p)
		np[len(p)] = k
		Walk(g, np, f)
	}
}

// Walk transverses the Group g and calls the WalkFn f on each Track.
func Walk(g Group, p Path, f WalkFn) {
	if gc, ok := g.(Collection); ok {
		walkCollection(gc, p, f)
		return
	}
	for _, t := range g.Tracks() {
		f(t, p)
	}
}

type FilterItem interface {
	Name() string
	Fields() map[string]interface{}
	Paths() []Path
}

type FilterItemSlice []FilterItem

func (f FilterItemSlice) Len() int           { return len(f) }
func (f FilterItemSlice) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FilterItemSlice) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

type filterItem struct {
	name   string
	fields map[string]interface{}
	paths  []Path
}

func (f *filterItem) Name() string                   { return f.name }
func (f *filterItem) Fields() map[string]interface{} { return f.fields }
func (f *filterItem) Paths() []Path                  { return f.paths }

func Filter(c Collection, field string) []FilterItem {
	m := make(map[string][]Path)
	walkfn := func(t Track, p Path) {
		f := t.GetString(field)
		m[f] = append(m[f], p)
	}
	Walk(c, Path([]Key{"Root"}), walkfn)

	result := make([]FilterItem, 0, len(m))
	for k, v := range m {
		result = append(result, &filterItem{
			name:   k,
			fields: make(map[string]interface{}),
			paths:  Union(v),
		})
	}
	sort.Sort(FilterItemSlice(result))
	return result
}

type trackPath struct {
	t Track
	p Path
}

type trackPathSorter struct {
	tp []trackPath
	fn LessFn
}

func (tps trackPathSorter) Len() int           { return len(tps.tp) }
func (tps trackPathSorter) Swap(i, j int)      { tps.tp[i], tps.tp[j] = tps.tp[j], tps.tp[i] }
func (tps trackPathSorter) Less(i, j int) bool { return tps.fn(tps.tp[i].t, tps.tp[j].t) }

func Recent(c Collection, n int) []Path {
	var trackPaths []trackPath
	walkfn := func(t Track, p Path) {
		trackPaths = append(trackPaths, trackPath{t, p})
	}
	Walk(c, Path([]Key{"Root"}), walkfn)

	sort.Sort(sort.Reverse(trackPathSorter{trackPaths, SortByTime("DateAdded")}))

	dedup := make(map[string]bool)
	result := make([]Path, 0, n)
	for _, tp := range trackPaths {
		e := tp.p.Encode()
		if !dedup[e] {
			dedup[e] = true
			result = append(result, tp.p)
		}
		if len(result) == n {
			break
		}
	}
	return result
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
		gg.add(fmt.Sprintf("%v", ga.Value(t)), t)
	}
	return gg
}
