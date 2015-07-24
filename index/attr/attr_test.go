// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package attr

import (
	"fmt"
	"reflect"
	"testing"
)

type testGetter struct {
	Name     string
	Duration int
	Artist   []string
}

func (t testGetter) GetString(f string) string {
	if f == "Name" {
		return t.Name
	}
	panic(fmt.Sprintf("invalid string field '%v'", f))
}

func (t testGetter) GetInt(f string) int {
	if f == "Duration" {
		return t.Duration
	}
	panic(fmt.Sprintf("invalid int field '%v'", f))
}

func (t testGetter) GetStrings(f string) []string {
	if f == "Artist" {
		return t.Artist
	}
	panic(fmt.Sprintf("invalid strings field '%v'", f))
}

func TestValue(t *testing.T) {
	name := "A"
	duration := 1
	artist := []string{"B"}

	g := testGetter{
		Name:     name,
		Duration: duration,
		Artist:   artist,
	}

	n := String("Name")
	got := n.Value(g)
	if got != name {
		t.Errorf("n.Value() = %#v, expected: %#v", got, name)
	}

	d := Int("Duration")
	got = d.Value(g)
	if got != duration {
		t.Errorf("d.Value() = %#v, expected: %#v", got, duration)
	}

	a := Strings("Artist")
	got = a.Value(g)
	if !reflect.DeepEqual(got, artist) {
		t.Errorf("d.Value() = %#v, expected: %#v", got, artist)
	}
}

func TestValueTypeIsEmpty(t *testing.T) {
	tests := []struct {
		a Interface
		v interface{}
		b bool
	}{
		// String attributes
		// Should be empty
		{
			String("String"),
			"",
			true,
		},

		// Should not be empty
		{
			String("String"),
			nil,
			false,
		},
		{
			String("String"),
			nil,
			false,
		},
		{
			String("String"),
			"A",
			false,
		},

		// Int attributes
		// Should be empty
		{
			Int("Int"),
			0,
			true,
		},

		// Should not be empty
		{
			Int("Int"),
			"",
			false,
		},
		{
			Int("Int"),
			nil,
			false,
		},
		{
			Int("Int"),
			1,
			false,
		},

		// Strings attributes
		// Should be empty
		{
			Strings("Strings"),
			nil,
			true,
		},
		{
			Strings("Strings"),
			[]string(nil),
			true,
		},

		// Should not be empty
		{
			Strings("Strings"),
			[]string{""},
			false,
		},
		{
			Strings("Strings"),
			[]string{"A"},
			false,
		},
	}

	for ii, tt := range tests {
		got := tt.a.IsEmpty(tt.v)
		if got != tt.b {
			t.Errorf("[%d] a.IsEmpty(%#v) = %#v, expected: %#v", ii, tt.v, got, tt.b)
		}
	}
}

func TestValueTypeIntersect(t *testing.T) {
	tests := []struct {
		a    Interface
		x, y interface{}
		z    interface{}
	}{
		// String attributes
		{
			String("String"),
			"",
			"",
			"",
		},
		{
			String("String"),
			"A",
			"A",
			"A",
		},
		{
			String("String"),
			"A",
			"B",
			"",
		},

		// Int attributes
		{
			Int("Int"),
			0,
			0,
			0,
		},
		{
			Int("Int"),
			1,
			1,
			1,
		},
		{
			Int("Int"),
			nil,
			false,
			0,
		},
		{
			Int("Int"),
			1,
			false,
			0,
		},

		// Strings attributes
		{
			Strings("Strings"),
			nil,
			nil,
			nil,
		},
		{
			Strings("Strings"),
			[]string{"A"},
			nil,
			nil,
		},
		{
			Strings("Strings"),
			nil,
			[]string{"A"},
			nil,
		},
		{
			Strings("Strings"),
			[]string{"A"},
			[]string{"A"},
			[]string{"A"},
		},
		{
			Strings("Strings"),
			[]string{"A", "B"},
			[]string{"C", "A"},
			[]string{"A"},
		},
	}

	for ii, tt := range tests {
		got := tt.a.Intersect(tt.x, tt.y)
		if !reflect.DeepEqual(got, tt.z) {
			t.Errorf("[%d] a.Intersect(%#v, %#v) = %#v, expected: %#v", ii, tt.x, tt.y, got, tt.z)
		}
	}
}

func TestStringSliceIntersect(t *testing.T) {
	tests := []struct {
		s, t []string
		r    []string
	}{
		{
			[]string(nil),
			[]string{"A"},
			[]string(nil),
		},

		{
			[]string{"A"},
			[]string{"B"},
			[]string(nil),
		},

		{
			[]string{""},
			[]string{""},
			[]string{""},
		},
	}

	for ii, tt := range tests {
		got := stringSliceIntersect(tt.s, tt.r)
		if !reflect.DeepEqual(got, tt.r) {
			t.Errorf("[%d] stringSliceIntersect(%#v, %#v) = %#v, expected: %#v", ii, tt.s, tt.t, got, tt.r)
		}
	}
}
