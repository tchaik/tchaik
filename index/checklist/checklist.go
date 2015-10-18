// Package checklist defines a collection of types for managing index checklists.
package checklist

import (
	"fmt"
	"sync"

	"tchaik.com/index"
)

// Store is an interface which defines methods necessary for setting and getting checklist items for
// index paths.
type Store interface {
	// Set the rating for the path.
	Set(index.Path, bool) error

	// Get the rating for the path.
	Get(index.Path) bool

	// List retuns a list of paths in the favourite Store.
	List() []index.Path
}

// NewStore creates a basic implementation of a checklist store, using the given path as the
// source of data. Note: we do not enforce any locking on the underlying file, which is read
// once to initialise the store, and then overwritten after each call to Set.
func NewStore(path string) (Store, error) {
	m := make(map[string]bool)
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

	m     map[string]bool
	store index.PersistStore
}

// Set implements Store.
func (s *store) Set(p index.Path, v bool) error {
	s.Lock()
	defer s.Unlock()

	k := fmt.Sprintf("%v", p)
	if v {
		s.m[k] = true
	} else {
		delete(s.m, k)
	}
	return s.store.Persist(&s.m)
}

// Get implements Store.
func (s *store) Get(p index.Path) bool {
	s.RLock()
	defer s.RUnlock()

	return s.m[fmt.Sprintf("%v", p)]
}

// List implements Store.
func (s *store) List() []index.Path {
	s.RLock()
	defer s.RUnlock()

	result := make([]index.Path, 0, len(s.m))
	for k := range s.m {
		result = append(result, index.NewPath(k))
	}
	return result
}
