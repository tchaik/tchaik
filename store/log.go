// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"io"
	"log"
	"net/http"
)

// LogFileSystem returns a wrapper around an http.FileSystem which logs calls
// to Open.
func LogFileSystem(prefix string, fs http.FileSystem) http.FileSystem {
	return logFileSystem{
		FileSystem: fs,
		prefix:     prefix,
	}
}

type logFileSystem struct {
	http.FileSystem
	prefix string
}

func (l logFileSystem) Open(path string) (http.File, error) {
	f, err := l.FileSystem.Open(path)
	if err != nil {
		log.Printf("%v open error: %v (%v)", l.prefix, path, err)
		return nil, err
	}
	log.Printf("%v open: %v", l.prefix, path)
	return f, err
}

// LogRWFileSystem returns a wrapper around a RWFileSystem which logs calls
// to Open and Create.
func LogRWFileSystem(prefix string, fs RWFileSystem) RWFileSystem {
	return logRWFileSystem{
		RWFileSystem: fs,
		prefix:       prefix,
	}
}

type logRWFileSystem struct {
	RWFileSystem
	prefix string
}

func (l logRWFileSystem) Open(path string) (http.File, error) {
	log.Printf("%v open: %v", l.prefix, path)
	return l.RWFileSystem.Open(path)
}

func (l logRWFileSystem) Create(path string) (io.WriteCloser, error) {
	log.Printf("%v create: %v", l.prefix, path)
	return l.RWFileSystem.Create(path)
}
