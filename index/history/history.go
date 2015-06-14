// Package history implements functionality for fetching/adding to play history.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tchaik/tchaik/index"
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
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		f, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()

	m := make(map[string][]time.Time)
	dec := json.NewDecoder(f)
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &basicStore{
		m:    m,
		path: path,
	}, nil
}

type basicStore struct {
	sync.RWMutex

	m    map[string][]time.Time
	path string
}

func (s *basicStore) persist() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(s.m)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

// Add implements Store.
func (s *basicStore) Add(p index.Path) error {
	s.Lock()
	defer s.Unlock()

	k := fmt.Sprintf("%v", p)
	s.m[k] = append(s.m[k], time.Now().UTC())
	return s.persist()
}

// Get implements Store.
func (s *basicStore) Get(p index.Path) []time.Time {
	return s.m[fmt.Sprintf("%v", p)]
}
