// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func TestSplitAfterMultiple(t *testing.T) {
	tests := []struct {
		in     string
		splits []string
		out    []string
	}{
		// Empty
		{
			"",
			[]string{},
			[]string{},
		},

		// Simple example (from docs for SplitAfter)
		{
			"a,b,c",
			[]string{","},
			[]string{"a,", "b,", "c"},
		},

		// Multiple splits, same input as previous
		{
			"a,b,c",
			[]string{",", "-"},
			[]string{"a,", "b,", "c"},
		},

		// Multiple splits, multiple splits to be made
		{
			"a,b:c,d:e",
			[]string{",", ":"},
			[]string{"a,", "b:", "c,", "d:", "e"},
		},
	}

	for ii, tt := range tests {
		got := splitAfterMultiple(tt.in, tt.splits)
		if !reflect.DeepEqual(tt.out, got) {
			t.Errorf("[%d] splitAfterMultiple(%#v, %#v) = %#v, expected %#v", ii, tt.in, tt.splits, got, tt.out)
		}
	}
}

func TestPrefixCollector(t *testing.T) {
	prefix1, prefix2, prefix3 := "Symphony No. 1", "Symphony No. 2", "Preludes, Book 1"
	table := []struct {
		in         []testTrack
		names      []string
		trackNames [][]string
		sizes      []int
	}{
		// One track
		{
			[]testTrack{
				{Name: "Fantasia on \"Greensleeves\""},
			},
			[]string{""},
			[][]string{
				[]string{
					"Fantasia on \"Greensleeves\"",
				},
			},
			[]int{1},
		},

		// No prefix matches
		{
			[]testTrack{
				{Name: "Speak to Me"},
				{Name: "Breathe"},
				{Name: "On the Run"},
				{Name: "Time"},
				{Name: "The Great Gig in the Sky"},
			},
			[]string{""},
			[][]string{
				[]string{
					"Speak to Me",
					"Breathe",
					"On the Run",
					"Time",
					"The Great Gig in the Sky",
				},
			},
			[]int{5},
		},

		// No prefix matches
		{
			[]testTrack{
				{Name: "Speak to Me"},
				{Name: "Speak to Me"},
			},
			[]string{""},
			[][]string{
				[]string{
					"Speak to Me",
					"Speak to Me",
				},
			},
			[]int{2},
		},

		// Prefix is a full word
		{
			[]testTrack{
				{Name: "Speak To Me (Live)"},
				{Name: "Speak To Me"},
			},
			[]string{""},
			[][]string{
				[]string{
					"Speak To Me (Live)",
					"Speak To Me",
				},
			},
			[]int{2},
		},

		// Simple 3 grouping, but no prefix splitter so nothing.
		{
			[]testTrack{
				{Name: prefix1 + " I. Andante"},
				{Name: prefix1 + " II. Adagio"},
				{Name: prefix1 + " III. Allegro"},
			},
			[]string{""},
			[][]string{
				[]string{
					prefix1 + " I. Andante",
					prefix1 + " II. Adagio",
					prefix1 + " III. Allegro",
				},
			},
			[]int{3},
		},

		// Simple 3 grouping, using ':' after prefix which should be stripped
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
			},
			[]string{prefix1},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{3},
		},

		// Simple 3 grouping, using '-' after prefix which should be stripped
		{
			[]testTrack{
				{Name: prefix1 + " - I. Andante"},
				{Name: prefix1 + " - II. Adagio"},
				{Name: prefix1 + " - III. Allegro"},
			},
			[]string{prefix1},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{3},
		},

		// Simple 3 grouping, prefix should not include No.
		{
			[]testTrack{
				{Name: prefix3 + ": No. 1: Danseuses de Delphes"},
				{Name: prefix3 + ": No. 2: Voiles"},
				{Name: prefix3 + ": No. 3: Le vent dans la plaine"},
			},
			[]string{prefix3},
			[][]string{
				[]string{
					"No. 1: Danseuses de Delphes",
					"No. 2: Voiles",
					"No. 3: Le vent dans la plaine",
				},
			},
			[]int{3},
		},

		// Simple 3, 3 grouping
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
			},
			[]string{prefix1, prefix2},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{3, 3},
		},

		// First track is different
		{
			[]testTrack{
				{Name: "Hebrides Overture"},
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
			},
			[]string{"", prefix1, prefix2},
			[][]string{
				[]string{
					"Hebrides Overture",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{1, 3, 3},
		},

		// First two tracks are different
		{
			[]testTrack{
				{Name: "Hebrides Overture"},
				{Name: "March of the Slaves"},
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
			},
			[]string{"", prefix1, prefix2},
			[][]string{
				[]string{
					"Hebrides Overture",
					"March of the Slaves",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{2, 3, 3},
		},

		// Middle track is different
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: "Hebrides Overture"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
			},
			[]string{prefix1, "", prefix2},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"Hebrides Overture",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{3, 1, 3},
		},

		// Middle two tracks are different
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: "Hebrides Overture"},
				{Name: "March of the Slaves"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
			},
			[]string{prefix1, "", prefix2},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"Hebrides Overture",
					"March of the Slaves",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
			},
			[]int{3, 2, 3},
		},

		// Last track is different
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
				{Name: "Hebrides Overture"},
			},
			[]string{prefix1, prefix2, ""},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"Hebrides Overture",
				},
			},
			[]int{3, 3, 1},
		},

		// Last two tracks are different
		{
			[]testTrack{
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
				{Name: "Hebrides Overture"},
				{Name: "March of the Slaves"},
			},
			[]string{prefix1, prefix2, ""},
			[][]string{
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"Hebrides Overture",
					"March of the Slaves",
				},
			},
			[]int{3, 3, 2},
		},

		// Different tracks a the beginning and the end
		{
			[]testTrack{
				{Name: "Hebrides Overture"},
				{Name: prefix1 + ": I. Andante"},
				{Name: prefix1 + ": II. Adagio"},
				{Name: prefix1 + ": III. Allegro"},
				{Name: "March of the Slaves"},
				{Name: prefix2 + ": I. Andante"},
				{Name: prefix2 + ": II. Adagio"},
				{Name: prefix2 + ": III. Allegro"},
				{Name: "Fantasy Overture"},
			},
			[]string{"", prefix1, "", prefix2, ""},
			[][]string{
				[]string{
					"Hebrides Overture",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"March of the Slaves",
				},
				[]string{
					"I. Andante",
					"II. Adagio",
					"III. Allegro",
				},
				[]string{
					"Fantasy Overture",
				},
			},
			[]int{1, 3, 1, 3, 1},
		},

		{
			[]testTrack{
				{Name: "Pictures At an Exhibition: Promenade"},
				{Name: "Pictures At an Exhibition: No. 1 The Gnome - Gnomus"},
				{Name: "Pictures At an Exhibition: Promenade"},
				{Name: "Pictures At an Exhibition: No. 2 The Old Castle - II Vecchio Castello"},
				{Name: "Pictures At an Exhibition: Promenade"},
				{Name: "Pictures At an Exhibition: No. 3 The Tuileries - Tuileries"},
				{Name: "Pictures At an Exhibition: No. 4 The Ox Card - Bydo"},
			},
			[]string{"Pictures At an Exhibition"},
			[][]string{
				[]string{
					"Promenade",
					"No. 1 The Gnome - Gnomus",
					"Promenade",
					"No. 2 The Old Castle - II Vecchio Castello",
					"Promenade",
					"No. 3 The Tuileries - Tuileries",
					"No. 4 The Ox Card - Bydo",
				},
			},
			[]int{7},
		},

		{
			[]testTrack{
				{Name: "Mussorgsky: Night on Bare Mountain"},
				{Name: "Mussorgsky: Pictures At an Exhibition: Promenade"},
				{Name: "Mussorgsky: Pictures At an Exhibition: No. 1 The Gnome - Gnomus"},
				{Name: "Mussorgsky: Pictures At an Exhibition: Promenade"},
				{Name: "Mussorgsky: Pictures At an Exhibition: No. 2 The Old Castle - II Vecchio Castello"},
				{Name: "Mussorgsky: Pictures At an Exhibition: Promenade"},
				{Name: "Mussorgsky: Pictures At an Exhibition: No. 3 The Tuileries - Tuileries"},
				{Name: "Mussorgsky: Pictures At an Exhibition: No. 4 The Ox Card - Bydo"},
			},
			[]string{"Mussorgsky", "Mussorgsky: Pictures At an Exhibition"},
			[][]string{
				[]string{
					"Night on Bare Mountain",
				},
				[]string{
					"Promenade",
					"No. 1 The Gnome - Gnomus",
					"Promenade",
					"No. 2 The Old Castle - II Vecchio Castello",
					"Promenade",
					"No. 3 The Tuileries - Tuileries",
					"No. 4 The Ox Card - Bydo",
				},
			},
			[]int{1, 7},
		},
	}

	for ii, tt := range table {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic from item: %d\n", ii)
					panic(r)
				}
			}()

			g := ByPrefix("Name").Collect(testTracker(tt.in))
			gotNames := names(g)
			gotSizes := tracksLen(g)
			gotTracks := collectionTrackStrings(g, "Name")

			if !reflect.DeepEqual(gotNames, tt.names) {
				t.Errorf("[%d] g.Names() = %#v, expected %#v", ii, gotNames, tt.names)
			}

			if !reflect.DeepEqual(gotSizes, tt.sizes) {
				t.Errorf("[%d] g.Sizes() = %#v, expected %#v", ii, gotSizes, tt.sizes)
			}

			if !reflect.DeepEqual(gotTracks, tt.trackNames) {
				t.Errorf("[%d] g.[Groups].[GetString(\"Name\")] = %#v, expected %#v", ii, gotTracks, tt.trackNames)
			}
		}()
	}
}

func tracksLen(g Collection) []int {
	result := make([]int, len(g.Keys()))
	for i, k := range g.Keys() {
		result[i] = len(g.Get(k).Tracks())
	}
	return result
}

func collectionTrackStrings(c Collection, field string) [][]string {
	result := make([][]string, len(c.Keys()))
	for i, k := range c.Keys() {
		result[i] = groupTrackStrings(c.Get(k), field)
	}
	return result
}

func groupTrackStrings(g Group, field string) []string {
	result := make([]string, len(g.Tracks()))
	for i, t := range g.Tracks() {
		result[i] = t.GetString(field)
	}
	return result
}
