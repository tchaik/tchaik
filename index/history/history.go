// Package history implements functionality for fetching/adding to play history.
package history

import (
	"fmt"
	"sync"
	"time"

	"tchaik.com/index"
)

// Store is an interface which defines methods necessary for fetching/adding to play history for
// index paths.  All times are stored in UTC.
type Store interface {
	// Add a play event to the store.
	Add(index.Path) error
	// Get the play events associated to a path.
	Get(index.Path) []time.Time
}

// NewStore creates a basic implementation of a play history store, using the given path as the
// source of data. If the file does not exist it will be created.
func NewStore(path string) (Store, error) {
	m := make(map[string][]time.Time)
	s, err := index.NewPersistStore(path, &m)
	if err != nil {
		return nil, err
	}

	return &basicStore{
		m:     m,
		store: s,
	}, nil
}

type basicStore struct {
	sync.RWMutex

	m     map[string][]time.Time
	store index.PersistStore
}

// Add implements Store.
func (s *basicStore) Add(p index.Path) error {
	s.Lock()
	defer s.Unlock()

	k := fmt.Sprintf("%v", p)
	s.m[k] = append(s.m[k], time.Now().UTC())
	return s.store.Persist(&s.m)
}

// Get implements Store.
func (s *basicStore) Get(p index.Path) []time.Time {
	s.RLock()
	defer s.RUnlock()

	return s.m[fmt.Sprintf("%v", p)]
}
