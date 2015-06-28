// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"github.com/dhowden/tag"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/nfnt/resize"
)

// ArtworkFileSystem wraps a FileSystem, reworking file system operations
// to refer to artwork from the underlying file.
func ArtworkFileSystem(fs FileSystem) FileSystem {
	return artworkFileSystem{
		FileSystem: fs,
	}
}

type artworkFileSystem struct {
	FileSystem
}

// Open the given file and return an http.File which contains the artwork, and hence
// the Name() of the returned file will have an extention for the artwork, not the
// media file.
func (afs artworkFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	f, err := afs.FileSystem.Open(ctx, path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var m tag.Metadata
	m, err = tag.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("error extracting picture from '%v': %v", path, err)
	}

	p := m.Picture()
	if p == nil {
		return nil, fmt.Errorf("no picture attached to '%v'", path)
	}

	name := stat.Name()
	if p.Ext != "" {
		name += "." + p.Ext
	}

	return &file{
		ReadSeeker: bytes.NewReader(p.Data),
		stat: &fileInfo{
			name:    name,
			size:    int64(len(p.Data)),
			modTime: stat.ModTime(),
		},
	}, nil
}

// FaviconFileSystem wraps another FileSystem assumed to contain only images, which are then
// resized to 48px x 48px and returned in .ico format.
func FaviconFileSystem(fs FileSystem) FileSystem {
	return faviconFileSystem{
		FileSystem: fs,
	}
}

type faviconFileSystem struct {
	FileSystem
}

// Open the given path (assumed to contain an image) and then resize to 48px x 48px and
// return in .ico format.
func (ffs faviconFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	f, err := ffs.FileSystem.Open(ctx, path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	filename := stat.Name()
	ext := filepath.Ext(filename)
	var img image.Image
	switch ext {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	default:
		err = fmt.Errorf("unsupported favicon image source: %v", stat.Name())
	}

	if err != nil {
		return nil, err
	}

	img = resize.Thumbnail(48, 48, img, resize.NearestNeighbor)
	buf := &bytes.Buffer{}
	err = ico.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	icoFilename := strings.TrimSuffix(filename, ext) + ".ico"

	return &file{
		ReadSeeker: bytes.NewReader(buf.Bytes()),
		stat: &fileInfo{
			name:    icoFilename,
			size:    int64(buf.Len()),
			modTime: stat.ModTime(),
		},
	}, nil
}
