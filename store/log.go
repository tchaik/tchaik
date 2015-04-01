// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"io"
	"log"
	"net/http"
)

// LogFileSystem is a wrapper around an http.FileSystem which logs requests
// to Open.
type LogFileSystem struct {
	Name string
	http.FileSystem
}

func (l LogFileSystem) Open(path string) (http.File, error) {
	f, err := l.FileSystem.Open(path)
	if err != nil {
		log.Printf("%v open error: %v (%v)", l.Name, path, err)
		return nil, err
	}
	log.Printf("%v open: %v", l.Name, path)
	return f, err
}

// LogRWFileSystem is a wrapper around a RWFileSystem which logs requests
// calls to Open and Create.
type LogRWFileSystem struct {
	RWFileSystem
}

func (l LogRWFileSystem) Open(path string) (http.File, error) {
	log.Printf("Open: %v", path)
	return l.RWFileSystem.Open(path)
}

func (l LogRWFileSystem) Create(path string) (io.WriteCloser, error) {
	log.Printf("Create: %v", path)
	return l.RWFileSystem.Create(path)
}
