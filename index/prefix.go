// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"strconv"
	"strings"
)

// pfxCol is the collection implementation for prefix-grouped collections
type pfxCol struct {
	col
	field string

	last string
	n    int
}

// pfxTrack is a track which truncates the specified field with by the given
// pfx length
type pfxTrack struct {
	Track
	field string
	pfx   int
}

// GetString returns field value for the given field name.
func (p pfxTrack) GetString(f string) string {
	if f == p.field && p.pfx > 0 {
		return strings.TrimSpace(p.Track.GetString(f)[p.pfx:])
	}
	return p.Track.GetString(f)
}

var prefixGroupSplit = []string{": ", " - ", "-"}

func (c *pfxCol) add(name string, t Track) {
	pfxLen := len(name)
	for _, x := range prefixGroupSplit {
		if strings.HasSuffix(name, x) {
			name = strings.TrimSuffix(name, x)
			break
		}
	}

	if c.last != name {
		c.n++
		c.last = name
	}
	c.col.addTrack(name, Key(strconv.Itoa(c.n)), pfxTrack{t, c.field, pfxLen})
}

// extension of strings.SplitAfter to split string multiple times using multiple
// split strings
func splitAfterMultiple(x string, s []string) []string {
	if len(s) == 0 {
		return []string{}
	}
	r := strings.SplitAfter(x, s[0])
	for _, sx := range s[1:] {
		t := make([]string, 0, len(r))
		for _, ry := range r {
			t = append(t, strings.SplitAfter(ry, sx)...)
		}
		r = t
	}
	return r
}

type item struct {
	before, after int // the number of words shared
}

type ByPrefix string

func (p ByPrefix) Collect(t Tracker) Collection {
	tracks := t.Tracks()
	field := string(p)

	newName := "Prefix collection"
	if tg, ok := t.(Group); ok {
		newName = tg.Name()
	}

	gg := pfxCol{col: newCol(newName), field: field}
	if len(tracks) == 1 {
		gg.add("", tracks[0])
		return gg
	}

	words := make([][]string, len(tracks))
	for i, t := range tracks {
		words[i] = splitAfterMultiple(t.GetString(field), prefixGroupSplit)
	}

	items := make([]item, len(tracks))
	items[0] = item{
		before: -1, // this will never be an option
		after:  largestPrefixWords(words[0], words[1]),
	}
	for i := 1; i < len(words)-1; i++ {
		items[i] = item{
			before: items[i-1].after,
			after:  largestPrefixWords(words[i], words[i+1]),
		}
	}
	items[len(tracks)-1] = item{
		before: items[len(tracks)-2].after,
		after:  -1, // this will never be an option
	}

	var name string
	var curr int
	for i, item := range items {
		if item.before >= item.after {
			// on the current one, so fine...
			if item.before == curr {
				gg.add(name, tracks[i])
				continue
			}
		}

		name = ""
		curr = item.after
		if item.after > 0 {
			name = strings.Join(words[i][:item.after], "")
		}
		gg.add(name, tracks[i])
	}

	return gg
}

func largestPrefixWords(s, t []string) int {
	if len(s) > len(t) {
		s, t = t, s
	}
	n := 0
	for i := range s {
		if s[i] != t[i] {
			break
		}
		n = i + 1
	}
	if len(s) == len(t) && len(s) == n {
		return 0
	}
	return n
}
