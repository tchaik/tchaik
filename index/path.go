// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PathSeparator is a string used to separate path components.
const PathSeparator string = ":"

// Key represents a unique value used to represent a group within a collection.
type Key string

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}

// Path is type which represents a position in the index heirarchy.  Each level has a key, and so the path
// is a slice of strings where each element is the key of some index element (group or track).
type Path []Key

// String implements Stringer.
func (p Path) String() string {
	return p.Encode()
}

// Encode returns a string representation of the Path, that is a PathSeparator'ed string where each
// component is a Key from the Path.
func (p Path) Encode() string {
	l := make([]string, len(p))
	for i, k := range p {
		l[i] = string(k)
	}
	return strings.Join(l, PathSeparator)
}

// Equal returns true iff q is Equal to p.
func (p Path) Equal(q Path) bool {
	if len(p) != len(q) {
		return false
	}
	for i, x := range p {
		if x != q[i] {
			return false
		}
	}
	return true
}

// Contains returns true iff q is contained within p.
func (p Path) Contains(q Path) bool {
	if len(p) == 0 || len(p) > len(q) {
		return false
	}
	for i, x := range p {
		if q[i] != x {
			return false
		}
	}
	return true
}

// NewPath creates a Path from the string representation.
func NewPath(x string) Path {
	split := strings.Split(x, PathSeparator)
	p := make([]Key, len(split))
	for i, s := range split {
		p[i] = Key(s)
	}
	return Path(p)
}

// PathFromStringSlice creates a path from the given []string.
func PathFromStringSlice(s []string) Path {
	p := make([]Key, len(s))
	for i, x := range s {
		p[i] = Key(x)
	}
	return p
}

// PathFromInterface reconstructs a Path from the given JSON-parsed interface{}
func PathFromJSONInterface(raw interface{}) (Path, error) {
	rawSlice, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected a slice of interface{}, got '%T'", raw)
	}

	path := make([]Key, len(rawSlice))
	for i, x := range rawSlice {
		s, ok := x.(string)
		if !ok {
			fl, ok := x.(float64)
			if !ok {
				return nil, fmt.Errorf("expected elements of type 'string' or 'float64', got '%T'", x)
			}
			s = strconv.Itoa(int(fl))
		}
		path[i] = Key(s)
	}
	return path, nil
}

// PathSlice is a wrapper type implementing sort.Interface (and index.Swapper).
type PathSlice []Path

// Swap implements sort.Interface (and index.Swapper).
func (p PathSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less implements sort.Interface.
func (p PathSlice) Less(i, j int) bool { return p[i].Encode() < p[j].Encode() }

// Len implements sort.Interface.
func (p PathSlice) Len() int { return len(p) }

// stringFreq is a helper type to count the number of occurances of a string.
type stringFreq struct {
	n int
	k string
}

// stringFreqSlice is a convenience type for sorting a slice of stringFreqs by the frequency,
// and then the alphabetical order of the strings.
type stringFreqSlice []stringFreq

func (c stringFreqSlice) Len() int      { return len(c) }
func (c stringFreqSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c stringFreqSlice) Less(i, j int) bool {
	if c[i].n < c[j].n {
		return true
	}
	if c[i].n > c[j].n {
		return false
	}
	return c[i].k < c[j].k
}

// Compute the intersection of the given lists of paths.
func OrderedIntersection(paths ...[]Path) []Path {
	if len(paths) == 0 {
		return []Path{}
	}

	enc := make(map[string]Path) // Encoding -> Path
	set := make(map[string]bool) // Set of Encoding
	cnt := make(map[string]int)

	for _, v := range paths[0] {
		e := v.Encode()
		set[e] = true
		enc[e] = v
		cnt[e]++
	}

	if len(paths) > 1 {
		for _, list := range paths[1:] {
			miss := make(map[string]bool)
			for k := range set {
				miss[k] = true
			}

			for _, v := range list {
				e := v.Encode()
				if set[e] {
					delete(miss, e)
					cnt[e]++
				}
			}

			for k := range miss {
				delete(set, k)
				delete(cnt, k)
			}
		}
	}

	result := make([]Path, 0, len(cnt))
	count := make([]stringFreq, 0, len(cnt))
	for k, v := range cnt {
		result = append(result, enc[k])
		count = append(count, stringFreq{v, k})
	}

	sort.Sort(ParallelSort(sort.Reverse(stringFreqSlice(count)), PathSlice(result)))
	return result
}

// Union returns a []Path which is the union (deduped) of the given slices of []Path.
func Union(l ...[]Path) []Path {
	if len(l) == 0 {
		return []Path{}
	}

	done := make(map[string]bool)
	res := make([]Path, 0)
	for _, paths := range l {
		for _, path := range paths {
			e := path.Encode()
			if !done[e] {
				done[e] = true
				res = append(res, path)
			}
		}
	}
	return res
}
