// Package favourite defines methods for setting/getting favourites for paths and
// persisting this data.
package favourite

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/tchaik/tchaik/index"
)

// Store is an interface which defines methods necessary for setting and getting favourites for
// index paths.
type Store interface {
	// Set the rating for the path.
	Set(index.Path, bool) error
	// Get the rating for the path.
	Get(index.Path) bool
}

// NewStore creates a basic implementation of a favourites store, using the given path as the
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

	m := make(map[string]bool)
	dec := json.NewDecoder(f)
	err = dec.Decode(&m)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &basicStore{
		m:    m,
		path: path,
	}, nil
}

type basicStore struct {
	sync.RWMutex

	m    map[string]bool
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

// Set implements Store.
func (s *basicStore) Set(p index.Path, v bool) error {
	s.Lock()
	defer s.Unlock()

	s.m[fmt.Sprintf("%v", p)] = v
	return s.persist()
}

// Get implements Store.
func (s *basicStore) Get(p index.Path) bool {
	s.RLock()
	defer s.RUnlock()

	return s.m[fmt.Sprintf("%v", p)]
}
