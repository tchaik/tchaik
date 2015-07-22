// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package index defines functionality for creating and manipulating a Tchaik music index.
package index

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Library is an interface which defines methods for listing tracks.
type Library interface {
	// Tracks returns a slice of all the tracks in the library.
	Tracks() []Track

	// Track returns the track from the given identifier and true if successful,
	// false otherwise.
	Track(identifier string) (Track, bool)
}

// Track is an interface which defines methods for retrieving track metadata.  Methods return
// zero values if attributes are unset and panic on undefined attributes and type mismatches.
type Track interface {
	// GetString returns a string value for the given attribute name.  Panics
	// if no such string attribute exists.
	GetString(string) string

	// GetStrings returns a list of strings for the given attribute name.
	// Panics if no such attribute exists.
	GetStrings(string) []string

	// GetInt returns an int value for the given attribute name.  Panics if no such
	// attribute exists.
	GetInt(string) int

	// GetTime returns a time.Time value for the given attribute name. Panics if no
	// such attribute exists.
	GetTime(string) time.Time
}

// Convert reads all the data exported by the Library and writes into the standard
// tchaik Library implementation.
// NB: The identifier field is set to be the value of ID on every track, regardless
// of whether this value has already been set in the input Library.
func Convert(l Library, id string) *library {
	allTracks := l.Tracks()
	tracks := make(map[string]*track, len(allTracks))

	for _, t := range allTracks {
		identifier := t.GetString(id)
		tracks[identifier] = &track{
			// string fields
			ID:          identifier,
			Name:        t.GetString("Name"),
			Album:       t.GetString("Album"),
			AlbumArtist: t.GetString("AlbumArtist"),
			Artist:      t.GetString("Artist"),
			Composer:    t.GetString("Composer"),
			Genre:       t.GetString("Genre"),
			Location:    t.GetString("Location"),

			// integer fields
			TotalTime:   t.GetInt("TotalTime"),
			Year:        t.GetInt("Year"),
			DiscNumber:  t.GetInt("DiscNumber"),
			TrackNumber: t.GetInt("TrackNumber"),
			TrackCount:  t.GetInt("TrackCount"),
			DiscCount:   t.GetInt("DiscCount"),
			BitRate:     t.GetInt("BitRate"),

			// date fields
			DateAdded:    t.GetTime("DateAdded"),
			DateModified: t.GetTime("DateModified"),
		}
	}
	return &library{
		tracks,
	}
}

// library is the default internal implementation Library which acts as the data
// source for all media tracks.
type library struct {
	trks map[string]*track
}

// Tracks implements Library.
func (l *library) Tracks() []Track {
	tracks := make([]Track, 0, len(l.trks))
	for _, t := range l.trks {
		tracks = append(tracks, t)
	}
	return tracks
}

// Track implements Library.
func (l *library) Track(id string) (Track, bool) {
	t, ok := l.trks[id]
	return t, ok
}

func (l *library) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.trks)
}

func (l *library) UnmarshalJSON(b []byte) error {
	l.trks = make(map[string]*track)
	return json.Unmarshal(b, &l.trks)
}

// WriteTo writes the Library data to the writer, currently using gzipped-JSON.
func WriteTo(l Library, w io.Writer) error {
	gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer gzw.Close()
	enc := json.NewEncoder(gzw)
	return enc.Encode(l)
}

// ReadFrom reads the gzipped-JSON representation of a Library.
func ReadFrom(r io.Reader) (Library, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()
	dec := json.NewDecoder(gzr)
	l := &library{}
	err = dec.Decode(l)
	return l, err
}

// track is the default implementation of the Track interface.
type track struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Album       string `json:"album,omitempty"`
	AlbumArtist string `json:"albumArtist,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Composer    string `json:"composer,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Location    string `json:"location,omitempty"`

	TotalTime   int `json:"totalTime,omitempty"`
	Year        int `json:"year,omitempty"`
	DiscNumber  int `json:"discNumber,omitempty"`
	TrackNumber int `json:"trackNumber,omitempty"`
	TrackCount  int `json:"trackCount,omitempty"`
	DiscCount   int `json:"discCount,omitempty"`
	BitRate     int `json:"bitRate,omitempty"`

	DateAdded    time.Time `json:"dateAdded,omitempty"`
	DateModified time.Time `json:"dateModified,omitempty"`
}

// GetString implements Track.
func (t *track) GetString(name string) string {
	switch name {
	case "ID":
		return t.ID
	case "Name":
		return t.Name
	case "Album":
		return t.Album
	case "AlbumArtist":
		return t.AlbumArtist
	case "Artist":
		return t.Artist
	case "Composer":
		return t.Composer
	case "Genre":
		return t.Genre
	case "Location":
		return t.Location
	}
	panic(fmt.Sprintf("unknown string field '%v'", name))
}

// DefaultGetStrings is a function which returns the default value for a GetStrings
// attribute which is based on an existing GetString attribute.  In particular, we handle
// the case where an empty 'GetString' attribute would be "", whereas the corresponding
// 'GetStrings' method should return 'nil' and not '[]string{""}'.
func DefaultGetStrings(t Track, f string) []string {
	v := t.GetString(f)
	if v != "" {
		return []string{v}
	}
	return nil
}

// GetStrings implements Track.
func (t *track) GetStrings(name string) []string {
	switch name {
	case "Artist", "AlbumArtist", "Composer":
		return DefaultGetStrings(t, name)
	}
	panic(fmt.Sprintf("unknown strings field '%v", name))
}

// GetInt implements Track.
func (t *track) GetInt(name string) int {
	switch name {
	case "TotalTime":
		return t.TotalTime
	case "Year":
		return t.Year
	case "DiscNumber":
		return t.DiscNumber
	case "TrackNumber":
		return t.TrackNumber
	case "TrackCount":
		return t.TrackCount
	case "DiscCount":
		return t.DiscCount
	case "BitRate":
		return t.BitRate
	}
	panic(fmt.Sprintf("unknown int field '%v'", name))
}

// GetTime implements Track.
func (t *track) GetTime(name string) time.Time {
	switch name {
	case "DateAdded":
		return t.DateAdded
	case "DateModified":
		return t.DateModified
	}
	panic(fmt.Sprintf("unknown time field '%v'", name))
}
