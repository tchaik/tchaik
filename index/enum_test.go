// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func TestEnumParseUInt(t *testing.T) {
	table := []struct {
		str  string
		v    uint
		fail bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"2", 2, false},

		{"01", 1, false},
		{"02", 2, false},

		{"", 0, true},
		{"i", 0, true},
		{"I", 0, true},
	}

	for ii, tt := range table {
		got, err := parseUInt(tt.str)
		if tt.fail && err == nil {
			t.Errorf("[%d] parseUInt(%s) = %d, expected error", ii, tt.str, got)
			continue
		}
		if !tt.fail && err != nil {
			t.Errorf("[%d] parseUInt(%s) gave error, expected %d", ii, tt.str, tt.v)
			continue
		}
		if got != tt.v {
			t.Errorf("[%d] parseUInt(%s) = %d, expected %d", ii, tt.str, got, tt.v)
		}
	}
}

func TestEnumParseNumeral(t *testing.T) {
	table := []struct {
		str  string
		v    uint
		fail bool
	}{
		{"I", 1, false},
		{"II", 2, false},
		{"III", 3, false},

		{"i", 1, false},
		{"ii", 2, false},
		{"iii", 3, false},

		{"", 0, true},
	}

	for ii, tt := range table {
		got, err := parseNumeral(tt.str)
		if tt.fail && err == nil {
			t.Errorf("[%d] parseNumeral(%s) = %d, expected error", ii, tt.str, got)
			continue
		}
		if !tt.fail && err != nil {
			t.Errorf("[%d] parseNumeral(%s) gave error, expected %d", ii, tt.str, tt.v)
			continue
		}
		if got != tt.v {
			t.Errorf("[%d] parseNumeral(%s) = %d, expected %d", ii, tt.str, got, tt.v)
		}
	}
}

func TestEnumIsNext(t *testing.T) {
	table := []struct {
		in   isNext
		strs []string
		n    uint
	}{
		// praser returns an error
		{isNext{fn: parseUInt}, []string{"A"}, 0},

		// parser returns 0
		{isNext{fn: parseUInt}, []string{"0"}, 0},

		// Starts at 2
		{isNext{fn: parseUInt}, []string{"2", "3", "4"}, 0},

		// Correct
		{isNext{fn: parseUInt}, []string{"1"}, 1},
		{isNext{fn: parseUInt}, []string{"1", "2"}, 2},
		{isNext{fn: parseUInt}, []string{"1", "2", "3"}, 3},

		// First correct, then non-matching number
		{isNext{fn: parseUInt}, []string{"1", "3"}, 1},

		// First correct, then error
		{isNext{fn: parseUInt}, []string{"1", "A"}, 1},
	}

	for ii, tt := range table {
		var count uint
		for _, s := range tt.strs {
			if tt.in.IsNext(s) {
				count++
			}
		}
		if tt.n != count {
			t.Errorf("[%d] count = %d, expected %d", ii, count, tt.n)
		}
	}
}

func TestTrimEnumSuffixes(t *testing.T) {
	table := []struct {
		in, out string
	}{
		{"1", "1"},
		{"1.", "1"},
		{"1..", "1"},
		{"1:", "1"},
		{"1-", "1"},
	}

	for ii, tt := range table {
		got := trimEnumFieldSuffixes(tt.in)
		if got != tt.out {
			t.Errorf("[%d] trimEnumSuffixes(%#v) = %#v, expected: %#v", ii, tt.in, got, tt.out)
		}
	}
}

func TestTrimPrefix(t *testing.T) {
	table := []struct {
		in, out string
		n       int
	}{
		{"1", "1", 0},
		{"ii", "ii", 0},
		{"No.1", "1", 3},
		{"No. 1", "1", 4},
	}

	for ii, tt := range table {
		got_x, got_n := trimPrefix(tt.in, enumWordPrefixes)
		if got_x != tt.out || got_n != tt.n {
			t.Errorf("[%d] trimEnumPrefix(%#v) = %#v, %d, expected: %#v, %d", ii, tt.in, got_x, got_n, tt.out, tt.n)
		}
	}
}

func TestEnumPrefix(t *testing.T) {
	table := []struct {
		in, out string
		outN    int
	}{
		// Edge cases
		{"", "", 0},
		{"1", "1", 0},

		// Strip the prefix and any suffix.
		{"1. Track Name", "1", 2},
		{"i. Track Name", "i", 2},
		{"I. Track Name", "I", 2},
		{"II. Track Name", "II", 3},

		// The enumeration prefix is the full input, so return everything.
		{"1.", "1.", 0},
		{"No. 1", "No. 1", 0},
	}

	for ii, tt := range table {
		gotStr, gotN := enumPrefix(tt.in)
		if gotStr != tt.out || gotN != tt.outN {
			t.Errorf("[%d] enumPrefix(%#v) = %#v, %d, expected: %#v, %d", ii, tt.in, gotStr, gotN, tt.out, tt.outN)
		}
	}
}

func getStrings(tracks []Track, name string) []string {
	res := make([]string, len(tracks))
	for i, t := range tracks {
		res[i] = t.GetString(name)
	}
	return res
}

func getTracks(names []string) []Track {
	res := make([]Track, len(names))
	for i, n := range names {
		res[i] = testTrack{Name: n}
	}
	return res
}

func TestStripEnumPrefix(t *testing.T) {
	table := []struct {
		in, out   []string
		listStyle string
	}{
		// Edge cases
		// Empty input
		{
			[]string{},
			[]string{},
			"",
		},
		// Invalid enumeration values
		{
			[]string{"One", "Two"},
			[]string{"One", "Two"},
			"",
		},
		// Invalid enumeration progression (1 -> 3)
		{
			[]string{"i", "iii"},
			[]string{"i", "iii"},
			"",
		},

		// Expected behaviour
		// Lower numeral parsing
		{
			[]string{"i First", "ii Second"},
			[]string{"First", "Second"},
			"upper-roman",
		},

		// Lower numeral parsing
		{
			[]string{"ii Second", "iii Third"},
			[]string{"Second", "Third"},
			"upper-roman",
		},

		// Upper-roman, no suffix
		{
			[]string{"I First", "II Second", "III Third", "IV Fourth"},
			[]string{"First", "Second", "Third", "Fourth"},
			"upper-roman",
		},

		// Lower numeral parsing, "." suffix
		{
			[]string{"i. First", "ii. Second"},
			[]string{"First", "Second"},
			"upper-roman",
		},

		// Lower numeral parsing, " - " suffix
		{
			[]string{"i - First", "ii - Second"},
			[]string{"First", "Second"},
			"upper-roman",
		},

		// Roman numeral, and remove suffix 1, 2, 3, 4
		{
			[]string{"I. First", "II. Second", "III. Third", "IV. Fourth"},
			[]string{"First", "Second", "Third", "Fourth"},
			"upper-roman",
		},

		// Simple 1, 2, 3, 4
		{
			[]string{"1. First", "2. Second", "3. Third", "4. Fourth"},
			[]string{"First", "Second", "Third", "Fourth"},
			"decimal",
		},

		// 01 -> 1, 02 -> 2, etc...
		{
			[]string{"01 First", "02 Second", "03 Third", "04 Fourth"},
			[]string{"First", "Second", "Third", "Fourth"},
			"decimal",
		},

		// Remove "No." prefix, remove ":" suffix
		{
			[]string{"No. 1: First", "No. 2: Second", "No. 3: Third", "No. 4: Fourth"},
			[]string{"First", "Second", "Third", "Fourth"},
			"decimal",
		},
	}

	for _, tt := range table {
		in := group{tracks: getTracks(tt.in)}
		got := TrimEnumPrefix(in)
		gotStr := getStrings(got.Tracks(), "Name")

		if !reflect.DeepEqual(gotStr, tt.out) {
			t.Errorf("Transform(%#v) = %#v, expected %#v", tt.in, gotStr, tt.out)
		}

		gotListStyle := got.Field("ListStyle")
		if tt.listStyle != gotListStyle {
			t.Errorf("got.Fields(\"ListStyle\") = %#v, expected %#v", gotListStyle, tt.listStyle)
		}
	}
}

type tracks []Track

func (t tracks) Tracks() []Track {
	return t
}

func TestTrimTrackNumPrefix(t *testing.T) {
	tests := []struct {
		tracks []testTrack
		titles []string
	}{
		// Empty input
		{
			[]testTrack{},
			[]string{},
		},

		// One track on one disc.
		{
			[]testTrack{
				testTrack{
					Name:        "01 One",
					TrackNumber: 1,
					DiscNumber:  1,
				},
			},
			[]string{
				"One",
			},
		},

		// Two tracks, incorrect prefix on second one.
		{
			[]testTrack{
				testTrack{
					Name:        "01 One",
					TrackNumber: 1,
					DiscNumber:  1,
				},
				testTrack{
					Name:        "03 Two",
					TrackNumber: 2,
					DiscNumber:  1,
				},
			},
			[]string{
				"01 One", "03 Two",
			},
		},

		// Two tracks, not consistent on the disc.
		{
			[]testTrack{
				testTrack{
					Name:        "01 One",
					TrackNumber: 1,
					DiscNumber:  1,
				},
				testTrack{
					Name:        "Two",
					TrackNumber: 2,
					DiscNumber:  1,
				},
			},
			[]string{
				"01 One", "Two",
			},
		},

		// Two tracks, inconsistent together, but consistent on the discs.
		{
			[]testTrack{
				testTrack{
					Name:        "01 One",
					TrackNumber: 1,
					DiscNumber:  1,
				},
				testTrack{
					Name:        "Two",
					TrackNumber: 2,
					DiscNumber:  2,
				},
			},
			[]string{
				"One", "Two",
			},
		},
	}

	for ii, tt := range tests {
		out := trimTrackNumPrefix("Name", "TrackNumber", "DiscNumber", testTracker(tt.tracks).Tracks())
		got := trackStrings(tracks(out), "Name")
		if !reflect.DeepEqual(got, tt.titles) {
			t.Errorf("[%d] got %#v, expected %#v", ii, got, tt.titles)
		}
	}
}
