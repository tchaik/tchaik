// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"reflect"
	"testing"
	"time"
)

type testTrack struct {
	Name, Album, Artist, Composer           string
	TrackNumber, DiscNumber, Duration, Year int
}

func (f testTrack) GetString(k string) string {
	switch k {
	case "Name":
		return f.Name
	case "Album":
		return f.Album
	case "Artist":
		return f.Artist
	case "Composer":
		return f.Composer
	}
	return ""
}

func (f testTrack) GetStrings(k string) []string {
	switch k {
	case "Artist", "AlbumArtist", "Composer":
		return DefaultGetStrings(f, k)
	}
	return nil
}

func (f testTrack) GetInt(k string) int {
	switch k {
	case "TrackNumber":
		return f.TrackNumber
	case "DiscNumber":
		return f.DiscNumber
	case "Duration":
		return f.Duration
	case "Year":
		return f.Year
	}
	return 0
}

func (f testTrack) GetBool(k string) bool {
	return false
}

func (f testTrack) GetTime(string) time.Time {
	return time.Time{}
}

type testTracker []testTrack

func (d testTracker) Tracks() []Track {
	result := make([]Track, len(d))
	for i, x := range d {
		result[i] = x
	}
	return result
}

func TestByAttr(t *testing.T) {
	trackListing := []testTrack{
		{Name: "A", Album: "Album A"},
		{Name: "B", Album: "Album A"},
		{Name: "C", Album: "Album B"},
		{Name: "D", Album: "Album B"},
	}

	expected := map[string][]Track{
		"Album A": {trackListing[0], trackListing[1]},
		"Album B": {trackListing[2], trackListing[3]},
	}

	attrGroup := ByAttr(StringAttr("Album")).Collect(testTracker(trackListing[:]))
	SortKeysByGroupName(attrGroup)

	nkm := nameKeyMap(attrGroup)
	for n, v := range expected {
		k, ok := nkm[n]
		if !ok {
			t.Errorf("%v is not a key of nkm", n)
		}

		gr := attrGroup.Get(Key(k))
		if gr == nil {
			t.Errorf("attrGroup.Get(%v) = nil, expected non-nil!", k)
		}

		got := gr.Tracks()
		if !reflect.DeepEqual(v, got) {
			t.Errorf("gr.Tracks() = %#v, expected: %#v", got, v)
		}
	}
}

func TestSubCollect(t *testing.T) {
	album1 := "Mahler Symphonies"
	album2 := "Shostakovich Symphonies"

	prefix1 := "Symphony No. 1 in D"
	prefix2 := "A B C Y"

	trackListing := []testTrack{
		{Name: prefix1 + ": I. Langsam, schleppend - Immer sehr gemächlich", Album: album1},
		{Name: prefix1 + ": II. Kräftig bewegt, doch nicht zu schnell - Recht gemächlich", Album: album1},
		{Name: prefix1 + ": III. Feierlich und gemessen, ohne zu schleppen", Album: album1},
		{Name: prefix2 + " 1", Album: album2},
		{Name: prefix2 + " 2", Album: album2},
		{Name: prefix2 + " 2", Album: album2},
	}
	expectedAlbums := []string{album1, album2}

	albums := ByAttr(StringAttr("Album")).Collect(testTracker(trackListing[:]))
	SortKeysByGroupName(albums)
	albNames := names(albums)

	if !reflect.DeepEqual(albNames, expectedAlbums) {
		t.Errorf("albums.Names() = %v, expected %#v", names, expectedAlbums)
	}

	albPfx := SubCollect(albums, ByPrefix("Name"))
	albPfxNames := names(albPfx)

	if !reflect.DeepEqual(albPfxNames, expectedAlbums[:]) {
		t.Errorf("albPfx.Names() = %v, expected %#v", albPfxNames, expectedAlbums)
	}

	nkm := nameKeyMap(albPfx)
	prefixGroup := albPfx.Get(nkm[album1])
	pfxCol, ok := prefixGroup.(Collection)
	if !ok {
		t.Errorf("expected a Collection, but got %T", prefixGroup)
	}

	pfxColNames := names(pfxCol)
	expectedPrefixGroupNames := []string{prefix1}

	if !reflect.DeepEqual(pfxColNames, expectedPrefixGroupNames) {
		t.Errorf("prefixCollection.Names() = %#v, expected %#v", pfxColNames, expectedPrefixGroupNames)
	}
}
