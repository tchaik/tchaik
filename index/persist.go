// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"encoding/json"
	"io"
	"os"
)

// PersistStore is a type which defines a simple persistence store.
type PersistStore string

// NewPersistStore creates a new PersistStore.  By default PersistStore uses
// JSON to persist data.
func NewPersistStore(path string, data interface{}) (PersistStore, error) {
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		f, err = os.Create(path)
		if err != nil {
			return "", err
		}
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(data)
	if err != nil && err != io.EOF {
		return "", err
	}
	return PersistStore(path), nil
}

func (p PersistStore) Persist(data interface{}) error {
	f, err := os.Create(string(p))
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}
