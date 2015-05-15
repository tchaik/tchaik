// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itl

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	rawitl "github.com/dhowden/itl"
	"github.com/tchaik/tchaik/index"
)

// ReadFrom creates a Tchaik Library implementation from an iTunes Music Library passed through
// an io.Reader.
func ReadFrom(r io.Reader) (index.Library, error) {
	l, err := rawitl.ReadFromXML(r)
	if err != nil {
		return nil, err
	}
	return &itlLibrary{&l}, nil
}

type itlLibrary struct {
	*rawitl.Library
}

// Implements Library.
func (l *itlLibrary) Tracks() []index.Track {
	tracks := make([]index.Track, 0, len(l.Library.Tracks))
	for _, t := range l.Library.Tracks {
		if strings.HasSuffix(t.Kind, "audio file") {
			x := t
			tracks = append(tracks, &itlTrack{&x})
		}
	}
	return tracks
}

// Implements Library.
func (l *itlLibrary) Track(id string) (index.Track, bool) {
	t, ok := l.Library.Tracks[id]
	if ok {
		return &itlTrack{&t}, true
	}
	return nil, false
}

// itlTrack is a wrapper type which implements Track for an rawitl.Track.
type itlTrack struct {
	*rawitl.Track
}

func decodeLocation(l string) (string, error) {
	u, err := url.ParseRequestURI(l)
	if err != nil {
		return "", err
	}
	// Annoyingly this doesn't replace &#38; (&)
	path := strings.Replace(u.Path, "&#38;", "&", -1)
	return path, nil
}

// GetString fetches the given string field in the Track, panics if field doesn't
// exist.
func (t *itlTrack) GetString(name string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name)
			panic(r)
		}
	}()

	switch name {
	case "Location":
		loc, err := decodeLocation(html.UnescapeString(t.Location))
		if err != nil {
			panic(fmt.Sprintf("error in decodeLocation: %v", err))
		}
		return loc
	case "TrackID":
		return strconv.Itoa(t.TrackID)
	case "Name":
		return html.UnescapeString(t.Name)
	case "Artist":
		return html.UnescapeString(t.Artist)
	case "Album":
		return html.UnescapeString(t.Album)
	case "AlbumArtist":
		return html.UnescapeString(t.AlbumArtist)
	case "Composer":
		return html.UnescapeString(t.Composer)
	case "Kind":
		return html.UnescapeString(t.Kind)
	}

	tt := reflect.TypeOf(t)
	ft, ok := tt.FieldByName(name)
	if !ok {
		panic(fmt.Sprintf("invalid field '%v'", name))
	}
	if ft.Type.Kind() != reflect.String {
		panic(fmt.Sprintf("field '%v' is not a string", name))
	}

	v := reflect.ValueOf(t)
	f := v.FieldByName(name)
	return html.UnescapeString(f.String())
}

// GetInt fetches the given int field in the Track, panics if field doesn't exist.
func (t *itlTrack) GetInt(name string) int {
	switch name {
	case "TrackID": // NB: This should really be read as a string
		return t.TrackID
	case "DiscNumber":
		return t.DiscNumber
	case "DiscCount":
		return t.DiscCount
	case "TrackNumber":
		return t.TrackNumber
	case "TrackCount":
		return t.TrackCount
	case "Year":
		return t.Year
	case "TotalTime":
		return t.TotalTime
	}

	tt := reflect.TypeOf(t)
	ft, ok := tt.FieldByName(name)
	if !ok {
		panic(fmt.Sprintf("invalid field '%v'", name))
	}
	if ft.Type.Kind() != reflect.Int {
		panic(fmt.Sprintf("field '%v' is not an int", name))
	}

	v := reflect.ValueOf(t)
	f := v.FieldByName(name)
	return int(f.Int())
}

// GetTime fetches the given time field in the Track, panics if field doesn't exist.
func (t *itlTrack) GetTime(name string) time.Time {
	switch name {
	case "DateAdded":
		return t.DateAdded
	case "DateModified":
		return t.DateModified
	}
	panic(fmt.Sprintf("field '%v' is not a time value", name))
}
