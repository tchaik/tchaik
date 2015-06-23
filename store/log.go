// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"io"
	"log"
	"net/http"
	"os"
)

// LogFileSystem returns a wrapper around an http.FileSystem which logs calls
// to Open.
func LogFileSystem(prefix string, fs http.FileSystem) http.FileSystem {
	if prefix != "" {
		prefix += " "
	}
	return logFileSystem{
		FileSystem: fs,
		Logger:     log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

type logFileSystem struct {
	http.FileSystem
	*log.Logger
}

// Open implements http.FileSystem.
func (l logFileSystem) Open(path string) (http.File, error) {
	f, err := l.FileSystem.Open(path)
	if err != nil {
		l.Printf("open error: %v", err)
		return nil, err
	}
	l.Printf("open: %v", path)
	return f, err
}

// LogRWFileSystem returns a wrapper around a RWFileSystem which logs calls
// to Open and Create.
func LogRWFileSystem(prefix string, fs RWFileSystem) RWFileSystem {
	if prefix != "" {
		prefix += " "
	}
	return logRWFileSystem{
		RWFileSystem: fs,
		Logger:       log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

type logRWFileSystem struct {
	RWFileSystem
	*log.Logger
}

// Open implements RWFileSystem.
func (l logRWFileSystem) Open(path string) (http.File, error) {
	f, err := l.RWFileSystem.Open(path)
	if err != nil {
		l.Printf("open error: %v", err)
		return nil, err
	}
	l.Printf("open: %v", path)
	return f, err
}

// Create implements RWFileSystem.
func (l logRWFileSystem) Create(path string) (io.WriteCloser, error) {
	f, err := l.RWFileSystem.Create(path)
	if err != nil {
		l.Printf("create error: %v", err)
	}
	l.Printf("create: %v", path)
	return f, err
}
