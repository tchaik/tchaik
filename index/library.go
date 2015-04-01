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
	"strconv"
)

// Library represents the tchaik music library.  Currently we don't have anything
// better than wrapping around the iTunes library (oh the shame!)
type Library interface {
	// Tracks returns a slice of all the tracks in the library.
	Tracks() []Track

	// Track returns the track from the given identifier, second return value true
	// if successful.
	Track(identifier string) (Track, bool)
}

// Track is an interface which defines a music file.
type Track interface {
	// GetString returns the string attribute with given name.
	GetString(string) string
	// GetInt returns the int attribute with given name.
	GetInt(string) int
}

var trackStringFields = []string{
	"TrackID", // unique identifier for the track.
	"Name",
	"Album",
	"AlbumArtist",
	"Artist",
	"Composer",
	"Location", // location of the associated file
}

var trackIntFields = []string{
	"TotalTime",
	"Year",
	"DiscNumber",
	"TrackNumber",
	"TrackCount",
	"DiscCount",
}

// Convert reads all the data exported by the Library and writes into the standard
// tchaik Library implementation.
// NB: The identifier field is set to be the value of TrackID on every track, regardless
// of whether this value has already been set in the input Library.
func Convert(l Library, id string) *library {
	allTracks := l.Tracks()
	tracks := make(map[string]*track, len(allTracks))

	n := len(trackStringFields) + len(trackIntFields)
	for _, t := range allTracks {
		m := make(map[string]interface{}, n)
		for _, f := range trackStringFields {
			m[f] = t.GetString(f)
		}
		for _, f := range trackIntFields {
			m[f] = t.GetInt(f)
		}

		identifier := t.GetString(id)
		m["TrackID"] = identifier
		tracks[identifier] = &track{
			flds: m,
		}
	}
	return &library{
		tracks,
	}
}

// library is the default internal implementation Library which acts as the data
// source for all media tracks.
type library struct {
	Trks map[string]*track // NB: Exported so that we can easily encode
}

// Tracks implements Library.
func (l *library) Tracks() []Track {
	tracks := make([]Track, 0, len(l.Trks))
	for _, t := range l.Trks {
		tracks = append(tracks, t)
	}
	return tracks
}

// Track implements Library.
func (l *library) Track(id string) (Track, bool) {
	t, ok := l.Trks[id]
	return t, ok
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
	flds map[string]interface{} // NB: Exported so that we can easily encode.
}

func (t *track) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.flds)
}

func (t *track) UnmarshalJSON(b []byte) error {
	t.flds = make(map[string]interface{})
	return json.Unmarshal(b, &t.flds)
}

// GetString implements Track.
func (t *track) GetString(name string) string {
	x, ok := t.flds[name]
	if !ok {
		panic(fmt.Sprintf("unknown string field '%v'", name))
	}
	if x == nil {
		panic(fmt.Sprintf("<nil> string field '%v'", name))
	}
	s, ok := x.(string)
	if !ok {
		panic(fmt.Sprintf("field '%v': expected string, got %#v (%T)", name, x, x))
	}
	return s
}

// GetInt implements Track.
func (t *track) GetInt(name string) int {
	x, ok := t.flds[name]
	if !ok {
		panic(fmt.Sprintf("unknown int field '%v'", name))
	}
	if x == nil {
		panic(fmt.Sprintf("<nil> int field '%v'", name))
	}
	switch x := x.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		n, err := strconv.Atoi(x)
		if err != nil {
			panic(fmt.Sprintf("error converting string to int: %v", err))
		}
		return n
	}
	panic(fmt.Sprintf("unknown type '%T' for field '%v'", name, name))
}
