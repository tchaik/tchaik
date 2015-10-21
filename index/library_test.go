package index

import (
	"testing"
	"time"
)

var tr = track{
	ID:          "ID",
	Name:        "Name",
	Album:       "Album",
	AlbumArtist: "AlbumArtist",
	Artist:      "Artist",
	Composer:    "Composer",
	Genre:       "Genre",
	Location:    "Location",
	Kind:        "Kind",

	TotalTime:   1,
	Year:        2,
	DiscNumber:  3,
	TrackNumber: 4,
	TrackCount:  5,
	DiscCount:   6,
	BitRate:     7,

	DateAdded:    time.Time{},
	DateModified: time.Time{},
}

func TestTrack(t *testing.T) {

	stringFields := []string{"ID", "Name", "Album", "AlbumArtist", "Artist", "Composer", "Genre", "Location", "Kind"}
	for _, f := range stringFields {
		got := tr.GetString(f)
		if got != f {
			t.Errorf("tr.GetString(%#v) = %#v, expected %#v", f, got, f)
		}
	}

	intFields := []string{"TotalTime", "Year", "DiscNumber", "TrackNumber", "TrackCount", "DiscCount", "BitRate"}
	for i, f := range intFields {
		got := tr.GetInt(f)
		expected := i + 1
		if got != expected {
			t.Errorf("tr.GetInt(%#v) = %d, expected %d", f, got, expected)
		}
	}
}
