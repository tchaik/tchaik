// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
)

// LogFileSystem returns a wrapper around an http.FileSystem which logs calls
// to Open.
func LogFileSystem(prefix string, fs FileSystem) FileSystem {
	if prefix != "" {
		prefix += " "
	}
	return logFileSystem{
		FileSystem: fs,
		Logger:     log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

type logFileSystem struct {
	FileSystem
	*log.Logger
}

// Open implements FileSystem.
func (l logFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	f, err := l.FileSystem.Open(ctx, path)
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
func (l logRWFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	f, err := l.RWFileSystem.Open(ctx, path)
	if err != nil {
		l.Printf("open error: %v", err)
		return nil, err
	}
	l.Printf("open: %v", path)
	return f, err
}

// Create implements RWFileSystem.
func (l logRWFileSystem) Create(ctx context.Context, path string) (io.WriteCloser, error) {
	f, err := l.RWFileSystem.Create(ctx, path)
	if err != nil {
		l.Printf("create error: %v", err)
	}
	l.Printf("create: %v", path)
	return f, err
}
