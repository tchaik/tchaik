// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"errors"
	"strconv"
	"strings"

	"github.com/dhowden/numerus"
)

// parseFn is a type which represents a parse function, given a string it will return
// a uint, or non-nil error if the string could not be parsed.
type parseFn func(string) (uint, error)

// parseUInt attempts to parse the given string into a UInt, returns an error if the input
// is invalid.
func parseUInt(s string) (uint, error) {
	n, err := strconv.Atoi(s)
	return uint(n), err
}

// parseNumeral is a wrapper around numerus.Parse, returning an error if the input string
// is invalid (when empty, or if an error is returned from numerus.Parse).
func parseNumeral(x string) (uint, error) {
	if x == "" {
		// numerus will return 0 (non-error) here, but we want an error
		return 0, errors.New("invalid numeral: empty string")
	}
	n, err := numerus.Parse(strings.ToUpper(x))
	return uint(n), err
}

// parser is a pairing of parse functions with "style" description.
type parser struct {
	Fn    parseFn
	Style string
}

// parsers is an internal list of numerical representations that can be identified by this
// package.
var parsers = []parser{
	{
		Fn:    parseUInt,
		Style: "decimal",
	},
	{
		Fn:    parseNumeral,
		Style: "upper-roman",
	},
}

type isNext struct {
	fn parseFn
	n  uint
}

func (i *isNext) IsNext(s string) bool {
	x, err := i.fn(s)
	if err != nil {
		return false
	}
	if x == (i.n + 1) {
		i.n++
		return true
	}
	return false
}

// enumFieldSuffixes is a list of string suffixes which will be trimmed from "enumeration
// fields"
var enumFieldSuffixes = []string{".", ":", "-"}

// trimFieldEnumSuffix removes all enumSuffixes from the given string and returns the
// result along with the number of characters removed.
func trimEnumFieldSuffixes(s string) string {
	for {
		var trim bool
		for _, x := range enumFieldSuffixes {
			if strings.HasSuffix(s, x) {
				s = strings.TrimSuffix(s, x)
				trim = true
				break
			}
		}
		if !trim {
			return s
		}
	}
}

// enumWordSuffixes is a list of strings which should be removed from string resulting
// from the prefix removal.
var enumWordSuffixes = []string{"-"}

// prefixes of an enumeration element which might need to be removed.
var enumWordPrefixes = []string{"No."}

// trimPrefix removes at most one prefix from prefixes from the given string, returns the
// the result.
func trimPrefix(s string, prefixes []string) (string, int) {
	l := len(s)
	for _, x := range prefixes {
		if strings.HasPrefix(strings.ToLower(s), strings.ToLower(x)) {
			s = strings.TrimSpace(s[len(x):])
			return s, l - len(s)
		}
	}
	return s, 0
}

// enumPrefix returns the expected enumeration prefix for the given string, and the
// number of characters that were removed from the beginning of the string to retrieve
// the value.  If the resulting prefix is the whole string then we don't do anything.
func enumPrefix(s string) (string, int) {
	x := s
	x, r := trimPrefix(x, enumWordPrefixes)
	words := strings.SplitN(x, " ", 2)
	if len(words) == 2 {
		_, n := trimPrefix(words[1], enumWordSuffixes)
		r += n
	}
	w := words[0]
	r += len(w)

	w = trimEnumFieldSuffixes(w)
	if r == len(s) {
		return s, 0
	}
	return w, r
}

func trimEnumPrefix(field string, tracks []Track) ([]Track, string) {
	if len(tracks) == 0 {
		return tracks, ""
	}

	t0 := tracks[0]
	name := t0.GetString(field)
	ep, epl := enumPrefix(name)

	var p parser
	var n uint
	for _, x := range parsers {
		var err error
		n, err = x.Fn(ep)
		if err == nil {
			p = x
			break
		}
	}

	if p.Fn == nil {
		return tracks, ""
	}

	nexter := &isNext{
		fn: p.Fn,
		n:  n,
	}

	epls := make([]int, len(tracks))
	epls[0] = epl
	i := 1
	for i < len(tracks) {
		ep, epls[i] = enumPrefix(tracks[i].GetString(field))
		if !nexter.IsNext(ep) {
			break
		}
		i++
	}

	if i < len(tracks) {
		return tracks, ""
	}

	result := make([]Track, 0, len(tracks))
	for i, t := range tracks {
		result = append(result, pfxTrack{
			Track: t,
			field: field,
			pfx:   epls[i],
		})
	}
	return result, p.Style
}

// TrimEnumPrefix removes enumeratative prefixes from the Tracks in the Group. The
// resulting group will have a ListStyle field which will indicate the style of enumeration
// (if any).
func TrimEnumPrefix(g Group) Group {
	nt, style := trimEnumPrefix("Name", g.Tracks())
	return subGrpFlds{
		Group:  g,
		tracks: nt,
		flds: map[string]interface{}{
			"ListStyle": style,
		},
	}
}

// TransformFn is a type which represents a function which Transforms a Group into
// another.
type TransformFn func(Group) Group

// Transform applies the TransformFn to the Group and returns the result.
func Transform(g Group, fn TransformFn) Group { return fn(g) }

// SubTransform recursively applies the TransformFn to each Group in the Collection.
func SubTransform(c Collection, fn TransformFn) Collection {
	m := make(map[Key]Group, len(c.Keys()))
	col := subCol{c, m, make(map[string]interface{})}
	for _, k := range c.Keys() {
		g := c.Get(k)
		if cg, ok := g.(Collection); ok {
			m[k] = SubTransform(cg, fn)
			continue
		}
		m[k] = fn(g)
	}
	return col
}
