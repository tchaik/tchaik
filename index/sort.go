// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import "sort"

// LessFn is a function type used for evaluating
type LessFn func(s, t Track) bool

type trackSorter struct {
	fn     LessFn
	tracks []Track
}

func (o *trackSorter) Len() int           { return len(o.tracks) }
func (o *trackSorter) Swap(i, j int)      { o.tracks[i], o.tracks[j] = o.tracks[j], o.tracks[i] }
func (o *trackSorter) Less(i, j int) bool { return o.fn(o.tracks[i], o.tracks[j]) }

// Sort sorts the slice of tracks using the given LessFn.
func Sort(tracks []Track, f LessFn) {
	o := &trackSorter{f, tracks}
	sort.Sort(o)
}

// SortByString returns a LessFn which orders Tracks using the GetString Attr on the given
// field.
func SortByString(field string) LessFn {
	return func(s, t Track) bool {
		return s.GetString(field) < t.GetString(field)
	}
}

// SortByInt returns a LessFn which orders Tracks using the GetInt Attr on the given field.
func SortByInt(field string) LessFn {
	return func(s, t Track) bool {
		return s.GetInt(field) < t.GetInt(field)
	}
}

// SortByTime returns a LessFn which orders Tracks using the GetTime Attr on the given field.
func SortByTime(field string) LessFn {
	return func(s, t Track) bool {
		return s.GetTime(field).Before(t.GetTime(field))
	}
}

// MultiSort creates a LessFn for tracks using the given LessFns.
func MultiSort(fns ...LessFn) LessFn {
	return func(s, t Track) bool {
		for _, fn := range fns[:len(fns)-1] {
			switch {
			case fn(s, t):
				return true
			case fn(t, s):
				return false
			}
		}
		return fns[len(fns)-1](s, t)
	}
}

// Swapper is an interface which defines the Swap method.
type Swaper interface {
	// Swap the items at indices i and j.
	Swap(i, j int)
}

// ParallelSort combines a sort.Interface implementation with a Swaper, and performs the same
// swap operations to w as they are applied to s.
func ParallelSort(s sort.Interface, w Swaper) sort.Interface {
	return &parallelSort{
		Interface: s,
		sw:        w,
	}
}

// parallelSort is a type allows for a Swapper implementation to be reordered in parallel
// to an implementation of sort.Interface.
type parallelSort struct {
	sort.Interface
	sw Swaper
}

// Swap implements Swapper (and sort.Interface) so that swaps are done on both the
// sort.Interface and Swapper.
func (p *parallelSort) Swap(i, j int) {
	p.Interface.Swap(i, j)
	p.sw.Swap(i, j)
}

// KeySlice attaches the methods of Swaper to []Key
type keySlice []Key

// Implements Swapper.
func (k keySlice) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func nameKeyMap(c Collection) map[string]Key {
	m := make(map[string]Key)
	for _, k := range c.Keys() {
		m[c.Get(k).Name()] = k
	}
	return m
}

func names(c Collection) []string {
	keys := c.Keys()
	n := make([]string, 0, len(keys))
	for _, k := range keys {
		n = append(n, c.Get(k).Name())
	}
	return n
}

// SorkKeysByGroupName sorts the names of the given collection (in place).  In particular, this
// assumes that g.Names() returns the actual internal representation of the listing.
func SortKeysByGroupName(c Collection) {
	sort.Sort(ParallelSort(sort.StringSlice(names(c)), keySlice(c.Keys())))
}
