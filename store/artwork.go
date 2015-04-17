// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/dhowden/tag"
)

// ArtworkFileSystem wraps another FileSystem, reworking file system operations
// to refer to artwork from the underlying file.
type ArtworkFileSystem struct {
	http.FileSystem
}

// Open the given file and return an http.File which contains the artwork, and hence
// the Name() of the returned file will have an extention for the artwork, not the
// media file.
func (afs ArtworkFileSystem) Open(path string) (http.File, error) {
	f, err := afs.FileSystem.Open(path)
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
