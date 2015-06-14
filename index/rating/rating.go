// Package rating defines types and methods for setting/getting ratings for paths and
// persisting this data.
package rating

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/tchaik/tchaik/index"
)

// Value is a type which represents a rating value.
type Value uint

// None is the Value used to mark a path as having no rating.
const None Value = 0

// IsValid returns true if the Value is valid.
func (v Value) IsValid() bool {
	return 0 <= v && v <= 5
}

// Store is an interface which defines methods necessary for setting and getting ratings for
// index paths.
type Store interface {
	// Set the rating for the path.
	Set(index.Path, Value) error
	// Get the rating for the path.
	Get(index.Path) Value
}

// NewStore creates a basic implementation of a ratings store, using the given path as the
// source of data. Note: we do not enforce any locking on the underlying file, which is read
// once to initialise the store, and then overwritten after each call to Set.
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

	m := make(map[string]Value)
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

	m    map[string]Value
	path string
}

func (s *basicStore) persist() error {
	f, err := os.Open(s.path)
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

// Set implements Store.
func (s *basicStore) Set(p index.Path, v Value) error {
	s.Lock()
	defer s.Unlock()

	s.m[fmt.Sprintf("%v", p)] = v
	return s.persist()
}

// Get implements Store.
func (s *basicStore) Get(p index.Path) Value {
	return s.m[fmt.Sprintf("%v", p)]
}
