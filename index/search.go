// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

var transformer = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

func removeNonAlphaNumeric(s string) string {
	in := []rune(s)
	res := make([]rune, len(in))
	i := 0
	for _, x := range s {
		switch {
		case x == '-':
			res[i] = ' '
			i++

		case unicode.IsLetter(x), unicode.IsDigit(x), unicode.IsSpace(x):
			res[i] = unicode.ToLower(x)
			i++
		}
	}
	result, _, _ := transform.String(transformer, string(res[:i]))
	return result
}

// Searcher is an interface which defines the Search method.
type Searcher interface {
	// Search uses the given string to filter a list of paths.
	Search(string) []Path
}

// WordIndex is an interface which defines the Words method.
type WordIndex interface {
	Searcher

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

// Expander is an interface which implements the Expand method.
type Expander interface {
	// Expand the given string into a list of alternatives.
	Expand(string) []string
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
func (p PrefixMultiExpand) Expand(s string) []string {
	if len(s) < MinPrefix {
		return []string{s}
	}
	if len(s) > p.size {
		// TODO: filter this with edit distance
		return p.words[s[:p.size]]
	}
	return p.words[s]
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
	e := es.Expander.Expand(s)
	ps := make([][]Path, len(e))
	for i, w := range e {
		ps[i] = es.Searcher.Search(w)
	}
	return Union(ps...)
}

// BuildPrefixExpandSearcher constructs a prefix expander which wraps the given Searcher
// by expanding each word in the search input using the WordIndex.
func BuildPrefixExpandSearcher(s Searcher, w WordIndex, n int) Searcher {
	return &expandSearcher{BuildPrefixMultiExpander(w.Words(), n), s}
}

type trackWordIndex struct {
	*wordIndex

	fields []string
}

// AddCollection recursively adds all Groups within the collection to the word index.
func (w *trackWordIndex) AddCollection(c Collection, p Path) {
	for _, k := range c.Keys() {
		np := make(Path, len(p), len(p)+1)
		copy(np, p)
		np = append(np, k)
		w.AddGroup(c.Get(k), np)
	}
}

// AddGroup adds the Group tracks to the word index, using the Path as root.
func (w *trackWordIndex) AddGroup(g Group, p Path) {
	if c, ok := g.(Collection); ok {
		w.AddCollection(c, p)
		return
	}
	// for i, t := range g.Tracks() {
	for _, t := range g.Tracks() {
		// np := make(Path, len(p), len(p)+1)
		// copy(p, np)
		// np = append(np, strconv.Itoa(i))
		for _, f := range w.fields {
			// s.AddString(t.GetString(f), np)
			w.AddString(t.GetString(f), p)
		}
	}
}

// BuildCollectionWordIndex creates a WordIndex using the given Collection, taking data from
// the given fields.
func BuildCollectionWordIndex(c Collection, fields []string) WordIndex {
	w := &trackWordIndex{
		wordIndex: &wordIndex{
			words: make(map[string][]Path),
		},
		fields: fields,
	}
	w.AddGroup(c, Path([]Key{"Root"}))
	return w
}

type wordSearchIntersect struct {
	Searcher
	min int // the minimum input string length before search returns something non-trivial.
}

func (s *wordSearchIntersect) Search(x string) []Path {
	if len(x) < s.min {
		return make([]Path, 0)
	}
	words := strings.Fields(x)
	paths := make([][]Path, 0, len(words))
	for _, w := range words {
		if len(w) >= s.min {
			paths = append(paths, s.Searcher.Search(w))
		}
	}
	return OrderedIntersection(paths...)
}

// WordsSearchIntersect calls Search on the Searcher for each word in the input string
// and then returns the ordered intersection (Paths are ordered by the number of times
// they appear).
func WordsIntersectSearcher(s Searcher) Searcher {
	return &wordSearchIntersect{
		Searcher: s,
		min:      3,
	}
}

// FlatSearcher is a Searcher wrapper which flattens input strings (replaces any accented
// characters with their un-accented equivalents).
type FlatSearcher struct {
	Searcher
}

// Search implements Searcher.
func (f FlatSearcher) Search(s string) []Path {
	return f.Searcher.Search(removeNonAlphaNumeric(s))
}
