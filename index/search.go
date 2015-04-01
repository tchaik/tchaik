// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Path is type which represents a position in the index heirarchy.  Each level has a key, and so the path
// is a slice of strings where each element is the key of some index element (group or track).
type Path []string

// Encode returns a string representation of the Path.
func (p Path) Encode() string {
	return strings.Join([]string(p), ">>")
}

// PathSlice is a wrapper type implementing sort.Interface (and index.Swapper)
type PathSlice []Path

// Swap implements sort.Interface (and index.Swapper).
func (p PathSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less implements sort.Interface.
func (p PathSlice) Less(i, j int) bool { return p[i].Encode() < p[j].Encode() }

// Len implements sort.Interface.
func (p PathSlice) Len() int { return len(p) }

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

var transformer = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

func removeNonAlphaNumeric(s string) string {
	in := []rune(s)
	res := make([]rune, len(in))
	i := 0
	for _, x := range s {
		if x == '-' {
			res[i] = ' '
			i++
			continue
		}
		if unicode.IsLetter(x) || unicode.IsDigit(x) || unicode.IsSpace(x) {
			res[i] = unicode.ToLower(x)
			i++
		}
	}
	result, _, _ := transform.Bytes(transformer, []byte(string(res[:i])))
	return string(result)
}

// Searcher is an interface which defines the Search method.
type Searcher interface {
	// Search uses the given string to filter a list of paths.
	Search(string) []Path
}

// WordIndex is an interface which defines the Words method.
type WordIndex interface {
	// Words returns all the words in the index.
	Words() []string
}

type wordIndex struct {
	fields []string
	words  map[string][]Path // mapping from word -> paths
}

func (s *wordIndex) Words() []string {
	words := make([]string, 0, len(s.words))
	for k := range s.words {
		words = append(words, k)
	}
	return words
}

// Search returns a list of paths for the given word.
func (s *wordIndex) Search(w string) []Path {
	return s.words[w]
}

// AddWord adds the given string (assumed to be a single word) to the index.
func (s *wordIndex) AddWord(w string, p Path) {
	s.words[w] = append(s.words[w], p)
}

// AddString adds the string to the index, removing non-alphanumeric characters,
// normalising modified characters, and splitting into words.
func (s *wordIndex) AddString(x string, p Path) {
	x = removeNonAlphaNumeric(x)
	w := strings.Fields(x)
	for _, x := range w {
		s.AddWord(x, p)
	}
}

func (s *wordIndex) AddCollection(c Collection, p Path) {
	for _, k := range c.Keys() {
		np := make(Path, len(p), len(p)+1)
		copy(p, np)
		np = append(np, string(k))
		s.AddGroup(c.Get(k), np)
	}
}

// AddGroup adds the given group to the word index, using the Path as root
func (s *wordIndex) AddGroup(g Group, p Path) {
	if c, ok := g.(Collection); ok {
		s.AddCollection(c, p)
		return
	}
	// for i, t := range g.Tracks() {
	for _, t := range g.Tracks() {
		// np := make(Path, len(p), len(p)+1)
		// copy(p, np)
		// np = append(np, strconv.Itoa(i))
		for _, f := range s.fields {
			// s.AddString(t.GetString(f), np)
			s.AddString(t.GetString(f), p)
		}
	}
}

// Expander is an interface which implements the Expand method.
type Expander interface {
	// Expand the given string, returning the result and true if succesful, or false
	// otherwise.
	Expand(string) ([]string, bool)
}

// MinPrefix is the minimum number of characters that can be used in a prefix.
const MinPrefix = 3

// PrefixMultiExpand is a type which implements Expander
type PrefixMultiExpand struct {
	words map[string][]string
	size  int
}

// BuildPrefixMultiExpander builds an Expander with given length n.  All words mapped to by prefixes will
// have length greater than or equal to MinPrefix.
func BuildPrefixMultiExpander(words []string, n int) PrefixMultiExpand {
	if n < MinPrefix {
		panic(fmt.Sprintf("Size must be greater than MinPrefix (%d)", MinPrefix))
	}

	m := make(map[string][]string)
	for _, w := range words {
		if len(w) >= MinPrefix {
			last := n
			if len(w) <= n {
				last = len(w)
			}
			for i := MinPrefix; i <= last; i++ {
				m[w[:i]] = append(m[w[:i]], w)
			}
		}
	}

	return PrefixMultiExpand{
		words: m,
		size:  n,
	}
}

// Expand uses the prefix mapping to return a list of words which can be expanded from s.
func (p PrefixMultiExpand) Expand(s string) ([]string, bool) {
	if len(s) < MinPrefix {
		return nil, false
	}
	if len(s) > p.size {
		// TODO: filter this with edit distance
		return p.words[s[:p.size]], true
	}
	return p.words[s], true
}

// expandSearcher is an implementation of Searcher which applies the Expander to search
// input and then performs a search on each of the expanded outputs, union the results
// and returns as the Search result.
type expandSearcher struct {
	Expander
	Searcher
}

// Search implements Searcher, and uses the internal Expander to expand all words
// in the search expression and unions the results (deduping). Returns nil when
// the search term isn't above the MinPrefix.
func (es *expandSearcher) Search(s string) []Path {
	e, ok := es.Expander.Expand(s)
	if !ok {
		return nil
	}
	ps := make([][]Path, len(e))
	for i, w := range e {
		ps[i] = es.Searcher.Search(w)
	}
	return Union(ps...)
}

// BuildPrefixExpandSearcher
func BuildPrefixExpandSearcher(s Searcher, w WordIndex, n int) Searcher {
	return &expandSearcher{BuildPrefixMultiExpander(w.Words(), n), s}
}

func BuildWordIndex(c Collection, fields []string) *wordIndex {
	wi := &wordIndex{
		fields: fields,
		words:  make(map[string][]Path),
	}
	wi.AddGroup(c, Path([]string{"Root"}))
	return wi
}

type wordSearchIntersect struct {
	Searcher
	min int // the minimum input string length before search returns something non-trivial.
}

func (s *wordSearchIntersect) Search(x string) []Path {
	if len(x) < s.min {
		return make([]Path, 0)
	}
	words := strings.Fields(strings.ToLower(x))
	paths := make([][]Path, 0, len(words))
	for _, w := range words {
		if len(w) >= s.min {
			paths = append(paths, s.Searcher.Search(w))
		}
	}
	return OrderedIntersection(paths...)
}

type FlatSearcher struct {
	Searcher
}

func (f FlatSearcher) Search(s string) []Path {
	return f.Searcher.Search(removeNonAlphaNumeric(s))
}

// WordsSearchIntersect calls Search on the Searcher for each word in the input string
// and then returns the ordered intersection (Paths are ordered by the number of times
// they appear.
func WordsIntersectSearcher(s Searcher) Searcher {
	return &wordSearchIntersect{
		Searcher: s,
		min:      3,
	}
}

type countKeyPath struct {
	n int
	k string
}

type countKeyPathSlice []countKeyPath

func (c countKeyPathSlice) Len() int      { return len(c) }
func (c countKeyPathSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c countKeyPathSlice) Less(i, j int) bool {
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
	count := make([]countKeyPath, 0, len(cnt))
	for k, v := range cnt {
		result = append(result, enc[k])
		count = append(count, countKeyPath{v, k})
	}

	sort.Sort(&ParallelSort{sort.Reverse(countKeyPathSlice(count)), PathSlice(result)})
	return result
}

// Union returns a []Path which is the union (deduped) of the given slices of []Path.
func Union(l ...[]Path) []Path {
	if len(l) == 0 {
		return []Path{}
	}
	if len(l) == 1 {
		return l[0]
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
