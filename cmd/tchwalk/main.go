// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
tchwalk is a tool which transverses a directory tree and reads all supported audio files
(.mp3 and m4a - ID3.v1,2.{2,3,4} and ID4) and uses the metadata to create a Tchaik library.

Only tracks which have readable metadata will be added to the library.  Any errors are
logged to stdout.

As no other unique identifying data is know, the SHA1 sum of the file path is used as the
track's TrackID.
*/
package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/dhowden/tchaik/index"
)

var fileExtensions = []string{".mp3", ".m4a", ".flac"}

// Library is an implementation of index.Library.
type Library struct {
	tracks map[string]*Track
}

func (l *Library) Track(id string) (index.Track, bool) {
	t, ok := l.tracks[id]
	return t, ok
}

func (l *Library) Tracks() []index.Track {
	tracks := make([]index.Track, 0, len(l.tracks))
	for _, t := range l.tracks {
		tracks = append(tracks, t)
	}
	return tracks
}

// Track is a wrapper around tag.Metadata which implements index.Track
type Track struct {
	tag.Metadata
	Location string
}

func (m *Track) GetString(name string) string {
	switch name {
	case "Name":
		return m.Title()
	case "Album":
		return m.Album()
	case "Artist":
		return m.Artist()
	case "Composer":
		return m.Composer()
	case "Location":
		return m.Location
	case "TrackID":
		sum := sha1.Sum([]byte(m.Location))
		return string(fmt.Sprintf("%x", sum))
	}
	return ""
}

func (m *Track) GetInt(name string) int {
	switch name {
	case "Year":
		return m.Year()
	case "TrackNumber":
		x, _ := m.Track()
		return x
	case "TrackCount":
		_, n := m.Track()
		return n
	case "DiscNumber":
		x, _ := m.Disc()
		return x
	case "DiscCount":
		_, n := m.Disc()
		return n
	}
	return 0
}

func validExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(filepath.Base(path)))
	for _, x := range fileExtensions {
		if ext == x {
			return true
		}
	}
	return false
}

func walk(root string) <-chan string {
	ch := make(chan string)
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ch <- path
		return nil
	}

	go func() {
		err := filepath.Walk(root, fn)
		if err != nil {
			log.Println(err)
		}
		close(ch)
	}()
	return ch
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %v [root-directory] [outputfile]\n", os.Args[0])
		os.Exit(1)
	}

	tracks := make(map[string]*Track)
	files := walk(os.Args[1])
	for path := range files {
		if validExtension(path) {
			track, err := processPath(path)
			if err != nil {
				log.Printf("error processing '%v': %v\n", path, err)
				continue
			}
			tracks[path] = track
		}
	}

	l := &Library{
		tracks: tracks,
	}
	tchLibrary := index.Convert(l, "TrackID")

	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	index.WriteTo(tchLibrary, f)
}

func processPath(path string) (*Track, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return &Track{
		Metadata: m,
		Location: path,
	}, nil
}
