// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"tchaik.com/index"
	"tchaik.com/index/attr"
	"tchaik.com/store"
)

// Library is a type which encompases the components which form a full library.
type Library struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string]index.Filter
	recent      []index.Path
	searcher    index.Searcher
}

type libraryFileSystem struct {
	store.FileSystem
	index.Library
}

// Open implements store.FileSystem and rewrites ID values to their corresponding Location
// values using the index.Library.
func (l *libraryFileSystem) Open(ctx context.Context, path string) (http.File, error) {
	t, ok := l.Library.Track(strings.Trim(path, "/")) // IDs arrive with leading slash
	if !ok {
		return nil, fmt.Errorf("could not find track: %v", path)
	}

	loc := t.GetString("Location")
	if loc == "" {
		return nil, fmt.Errorf("invalid (empty) location for track: %v", path)
	}
	loc = filepath.ToSlash(loc)
	return l.FileSystem.Open(ctx, loc)
}

// Fetch fetches a Group and its corresponding Key given a path.  Returns an error if the path
// is invalid.
func (l *Library) Fetch(p index.Path) (index.Group, index.Key, error) {
	if len(p) == 0 {
		return nil, "", fmt.Errorf("invalid path: %v\n", p)
	}

	root := l.collections[string(p[0])]
	if root == nil {
		return nil, "", fmt.Errorf("unknown collection: %#v", p[0])
	}

	if len(p) == 1 {
		return root, p[0], nil
	}

	g, err := l.Build(root, p[1:])
	if err != nil {
		return nil, "", fmt.Errorf("error in Fetch: %v (path: %#v)", err, p[1:])
	}
	return g, p[1], nil
}

// Build fetches a Group from the index.Collection given by the Path.
func (l *Library) Build(c index.Collection, p index.Path) (index.Group, error) {
	if len(p) == 0 {
		return c, nil
	}

	g, err := index.GroupFromPath(&rootCollection{c}, p)
	if err != nil {
		return nil, err
	}
	g = index.FirstTrackAttr(attr.String("ID"), g)
	return g, nil
}

// FileSystem wraps the http.FileSystem in a library lookup which will translate /ID
// requests into their corresponding track paths.
func (l *Library) FileSystem(fs store.FileSystem) store.FileSystem {
	return store.Trace(&libraryFileSystem{fs, l.Library}, "libraryFileSystem")
}

// ExpandPaths constructs a collection (group) whose sub-groups are taken from the "Root"
// collection.
func (l *Library) ExpandPaths(paths []index.Path) index.Group {
	return &Group{
		Group: index.NewPathsCollection(l.collections["Root"], paths),
		Key:   index.Key("Root"),
	}
}
