// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func TestPathEqual(t *testing.T) {
	tests := []struct {
		p, q  Path
		equal bool
	}{
		{
			Path(nil), Path(nil),
			true,
		},
		{
			Path([]Key{}), Path([]Key{}),
			true,
		},
		{
			NewPath("a"), NewPath("a"),
			true,
		},
		{
			NewPath("a:b"), NewPath("a"),
			false,
		},
		{
			NewPath("a:b"), NewPath("a:c"),
			false,
		},
		{
			Path(nil), NewPath("a"),
			false,
		},
	}

	for ii, tt := range tests {
		if tt.p.Equal(tt.q) != tt.equal {
			t.Errorf("[%d] (%#v).Equal(%#v) = %v, expected %v", ii, tt.p, tt.q, !tt.equal, tt.equal)
		}
	}
}

func TestPathContains(t *testing.T) {
	tests := []struct {
		p, q     Path
		contains bool
	}{
		{
			Path(nil), Path(nil),
			false,
		},
		{
			Path([]Key{}), Path([]Key{}),
			false,
		},
		{
			NewPath("a"), NewPath("a"),
			true,
		},
		{
			NewPath("a:b"), NewPath("a"),
			false,
		},
		{
			NewPath("a:b"), NewPath("a:c"),
			false,
		},
		{
			Path(nil), NewPath("a"),
			false,
		},
		{
			NewPath("a:b"), NewPath("a:b:c"),
			true,
		},
	}

	for ii, tt := range tests {
		if tt.p.Contains(tt.q) != tt.contains {
			t.Errorf("[%d] (%#v).Contains(%#v) = %v, expected %v", ii, tt.p, tt.q, !tt.contains, tt.contains)
		}
	}
}

func TestOrderedIntersection(t *testing.T) {
	tests := []struct {
		in  [][]Path
		out []Path
	}{
		{
			in:  nil,
			out: []Path{},
		},

		{
			in: [][]Path{
				{NewPath("A")},
			},
			out: []Path{NewPath("A")},
		},

		{
			in: [][]Path{
				{NewPath("A")},
				{NewPath("B")},
			},
			out: []Path{},
		},

		{
			in: [][]Path{
				{NewPath("A")},
				{NewPath("B"), NewPath("A")},
			},
			out: []Path{NewPath("A")},
		},

		{
			in: [][]Path{
				{NewPath("A"), NewPath("B")},
				{NewPath("B"), NewPath("A")},
				{NewPath("A"), NewPath("B"), NewPath("C")},
				{NewPath("C"), NewPath("A"), NewPath("B"), NewPath("B")},
			},
			out: []Path{NewPath("B"), NewPath("A")},
		},
	}

	for ii, tt := range tests {
		got := OrderedIntersection(tt.in...)
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("[%d] got %#v, expected: %#v", ii, got, tt.out)
		}
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		in  [][]Path
		out []Path
	}{
		{
			in:  nil,
			out: []Path{},
		},

		{
			in: [][]Path{
				{NewPath("A")},
			},
			out: []Path{NewPath("A")},
		},

		{
			in: [][]Path{
				{NewPath("A")},
				{NewPath("B")},
			},
			out: []Path{NewPath("A"), NewPath("B")},
		},

		{
			in: [][]Path{
				{NewPath("A")},
				{NewPath("B"), NewPath("A")},
			},
			out: []Path{NewPath("A"), NewPath("B")},
		},

		{
			in: [][]Path{
				{NewPath("A"), NewPath("B")},
				{NewPath("B"), NewPath("A")},
				{NewPath("A"), NewPath("B"), NewPath("C")},
				{NewPath("C"), NewPath("A"), NewPath("B"), NewPath("B")},
			},
			out: []Path{NewPath("A"), NewPath("B"), NewPath("C")},
		},
	}

	for ii, tt := range tests {
		got := Union(tt.in...)
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("[%d] got %#v, expected: %#v", ii, got, tt.out)
		}
	}
}

func TestIndexOfPath(t *testing.T) {
	tests := []struct {
		haystack []Path
		needle   Path
		idx      int
	}{
		{
			[]Path{},
			Path{},
			-1,
		},
		{
			[]Path{
				Path{"Root"},
			},
			Path{"Root"},
			0,
		}, {
			[]Path{
				Path{"Root"}, Path{"Root", "One"},
			},
			Path{"Root", "One"},
			1,
		},
	}

	for ii, tt := range tests {
		i := IndexOfPath(tt.haystack, tt.needle)
		if i != tt.idx {
			t.Errorf("[%d] IndexOfPath(%v, %v) = %d, expected %d", ii, tt.haystack, tt.needle, i, tt.idx)
		}
	}
}

func TestPathFromJSONInterface(t *testing.T) {
	tests := []struct {
		in  interface{}
		out Path
	}{
		{
			interface{}(nil),
			Path(nil),
		},

		{
			[]interface{}{nil},
			Path(nil),
		},

		{
			[]interface{}{""},
			Path{""},
		},

		{
			[]interface{}{"", 123.},
			Path{"", "123"},
		},
	}

	for ii, tt := range tests {
		got, _ := PathFromJSONInterface(tt.in)
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("[%d] got: %#v, expected %#v", ii, got, tt.out)
		}
	}
}
