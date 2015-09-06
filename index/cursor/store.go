// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cursor

import (
	"sync"

	"tchaik.com/index"
)

// Store is an interface which defines methods for implementing a cursor store.
type Store interface {
	// Get returns the cursor for the given name.
	Get(name string) *Cursor

	// Set sets the cursor for the given name.
	Set(name string, p *Cursor) error

	// Delete removes the cursor with the given name.
	Delete(name string) error
}

// NewStore creates a basic implementation of a cursor store, using the given path as the
// source of data. If the file does not exist it will be created.
func NewStore(path string) (Store, error) {
	m := make(map[string]*Cursor)
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

	m     map[string]*Cursor
	store index.PersistStore
}

// Add implements Store.
func (s *store) Get(name string) *Cursor {
	s.Lock()
	defer s.Unlock()

	return s.m[name]
}

// Put implements Store.
func (s *store) Set(name string, c *Cursor) error {
	s.Lock()
	defer s.Unlock()

	s.m[name] = c
	return s.store.Persist(&s.m)
}

// Delete implements Store.
func (s *store) Delete(name string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.m, name)
	return s.store.Persist(&s.m)
}
