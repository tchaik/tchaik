// Package walk implements a path walker which reads audio files under a path and
// constructs an index.Library from supported metadata tags (see github.com/dhowden/tag).
package walk

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/tchaik/tchaik/index"
)

var fileExtensions = []string{".mp3", ".m4a", ".flac", ".ogg"}

// NewLibrary constructs an index.Library by walking through the directory tree under
// the given path.  Any errors are logged to stdout (TODO: fix this!)
func NewLibrary(path string) index.Library {
	tracks := make(map[string]*track)
	files := walk(path)
	for p := range files {
		if validExtension(p) {
			t, err := processPath(p)
			if err != nil {
				log.Printf("error processing '%v': %v\n", p, err)
				continue
			}
			tracks[p] = t
		}
	}

	return &library{
		tracks: tracks,
	}
}

// library is an implementation of index.library.
type library struct {
	tracks map[string]*track
}

func (l *library) Track(id string) (index.Track, bool) {
	t, ok := l.tracks[id]
	return t, ok
}

func (l *library) Tracks() []index.Track {
	tracks := make([]index.Track, 0, len(l.tracks))
	for _, t := range l.tracks {
		tracks = append(tracks, t)
	}
	return tracks
}

// track is a wrapper around tag.Metadata which implements index.Track
type track struct {
	tag.Metadata
	Location    string
	FileInfo    os.FileInfo
	CreatedTime time.Time
}

func (m *track) GetString(name string) string {
	switch name {
	case "Name":
		title := m.Title()
		if title == "" {
			fileName := m.FileInfo.Name()
			ext := filepath.Ext(fileName)
			title = strings.TrimSuffix(fileName, ext)
		}
		return title
	case "Album":
		return m.Album()
	case "Artist":
		return m.Artist()
	case "AlbumArtist":
		return m.AlbumArtist()
	case "Composer":
		return m.Composer()
	case "Genre":
		return m.Genre()
	case "Location":
		return m.Location
	case "TrackID":
		sum := sha1.Sum([]byte(m.Location))
		return string(fmt.Sprintf("%x", sum))
	}
	return ""
}

func (m *track) GetInt(name string) int {
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

func (m *track) GetTime(name string) time.Time {
	switch name {
	case "DateModified":
		return m.FileInfo.ModTime()
	case "DateAdded":
		return m.CreatedTime
	}
	return time.Time{}
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

func processPath(path string) (*track, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	createdTime, err := getCreatedTime(path)
	if err != nil {
		return nil, err
	}

	return &track{
		Metadata:    m,
		Location:    path,
		FileInfo:    fileInfo,
		CreatedTime: createdTime,
	}, nil
}
