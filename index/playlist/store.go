// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playlist

import (
	"sync"

	"tchaik.com/index"
)

// Store is an interface which defines methods for implementing a playlist store.
type Store interface {
	// Names returns the name of playlists in the store.
	Names() []string

	// Get returns the playlist for the given name.
	Get(name string) *Playlist

	// Set sets the playlist for the given name.
	Set(name string, p *Playlist) error

	// Delete removes the playlist with the given name.
	Delete(name string) error
}

// NewStore creates a basic implementation of a playlist store, using the given path as the
// source of data. If the file does not exist it will be created.
func NewStore(path string) (Store, error) {
	m := make(map[string]*Playlist)
	s, err := index.NewPersistStore(path, &m)
	if err != nil {
		return nil, err
	}

	return &store{
		m:     m,
		store: s,
	}, nil
}

type store struct {
	sync.RWMutex

	m     map[string]*Playlist
	store index.PersistStore
}

// Names implements Store.
func (s *store) Names() []string {
	n := make([]string, 0, len(s.m))
	for k := range s.m {
		n = append(n, k)
	}
	return n
}

// Add implements Store.
func (s *store) Get(name string) *Playlist {
	s.Lock()
	defer s.Unlock()

	return s.m[name]
}

// Put implements Store.
func (s *store) Set(name string, p *Playlist) error {
	s.Lock()
	defer s.Unlock()

	s.m[name] = p
	return s.store.Persist(&s.m)
}

// Delete implements Store.
func (s *store) Delete(name string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.m, name)
	return s.store.Persist(&s.m)
}
