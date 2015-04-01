// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
)

func TestSortTracks(t *testing.T) {
	tracks := []Track{
		testTrack{Name: "A", Album: "X"},
		testTrack{Name: "B", Album: "X"},
		testTrack{Name: "C", Album: "X"},
		testTrack{Name: "A", Album: "Y"},
		testTrack{Name: "B", Album: "Y"},
		testTrack{Name: "C", Album: "Y"},
	}

	expectedNames := []string{"A", "A", "B", "B", "C", "C"}
	Sort(tracks, SortByString("Name"))

	for ii, tt := range tracks {
		if expectedNames[ii] != tt.GetString("Name") {
			t.Fatalf("Sort(tracks, ByString(Name)): %v", tt.GetString("Name"))
		}
	}
}

func TestMultiSort(t *testing.T) {
	expectedTracks := []Track{
		testTrack{Name: "A", Album: "X"},
		testTrack{Name: "B", Album: "X"},
		testTrack{Name: "C", Album: "X"},
		testTrack{Name: "A", Album: "Y"},
		testTrack{Name: "B", Album: "Y"},
		testTrack{Name: "C", Album: "Y"},
	}

	tracks := []Track{
		expectedTracks[4],
		expectedTracks[5],
		expectedTracks[0],
		expectedTracks[3],
		expectedTracks[1],
		expectedTracks[2],
	}

	Sort(tracks, MultiSort(SortByString("Album"), SortByString("Name")))
	if !reflect.DeepEqual(tracks, expectedTracks) {
		t.Errorf("Sort(...) = %v, expected %v", tracks, expectedTracks)
	}
}
