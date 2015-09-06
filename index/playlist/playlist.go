// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package playlist defines functionality for creating, manipulating and persisting playlists.
package playlist

import (
	"encoding/json"
	"fmt"

	"tchaik.com/index"
)

// Transformer is an interface which defines the Transform method for making changes to
// playlist items.
type Transformer interface {
	// Transform returns the action associated with the transformation.
	Transform() string
}

// RemoveItem is a type which defines a RemovePath action for a specific index.Path, which
// removes a path from a playlist item.
type RemovePath index.Path

// Transform implements Transformer
func (RemovePath) Transform() string {
	return "remove"
}

func (r RemovePath) MarshalJSON() ([]byte, error) {
	exp := struct {
		Action string     `json:"action"`
		Path   index.Path `json:"path"`
	}{
		Action: r.Transform(),
		Path:   index.Path(r),
	}
	return json.Marshal(exp)
}

// Item is a type which defines a playlist item.
type Item struct {
	path       index.Path
	transforms []Transformer
}

func newItem(path index.Path) *Item {
	return &Item{
		path: path,
	}
}

// AddTransform adds the Transformer to the Item.
func (i *Item) AddTransform(t Transformer) {
	i.transforms = append(i.transforms, t)
}

// MarshalJSON implements json.Marshaler.
func (i *Item) MarshalJSON() ([]byte, error) {
	exp := struct {
		Path       index.Path    `json:"path"`
		Transforms []Transformer `json:"transforms",omitempty`
	}{
		Path:       i.path,
		Transforms: i.transforms,
	}
	return json.Marshal(exp)
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Item) UnmarshalJSON(b []byte) error {
	exp := struct {
		Path       index.Path               `json:"path"`
		Transforms []map[string]interface{} `json:"transforms"`
	}{}
	err := json.Unmarshal(b, &exp)
	if err != nil {
		return err
	}

	transforms := make([]Transformer, 0, len(exp.Transforms))
	for _, t := range exp.Transforms {
		if t["action"] == "remove" {
			p, err := index.PathFromJSONInterface(t["path"])
			if err != nil {
				return fmt.Errorf("invalid format for path:", t["path"])
			}
			transforms = append(transforms, RemovePath(p))
		}
	}

	i.path = exp.Path
	i.transforms = transforms
	return nil
}

// Playlist is a basic implementation of a playlist
type Playlist struct {
	items []*Item
}

// MarshalJSON implements json.Marshaler.
func (p *Playlist) MarshalJSON() ([]byte, error) {
	exp := struct {
		Items []*Item `json:"items"`
	}{
		p.items,
	}
	return json.Marshal(exp)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Playlist) UnmarshalJSON(b []byte) error {
	exp := struct {
		Items []*Item `json:"items"`
	}{}
	err := json.Unmarshal(b, &exp)
	if err != nil {
		return err
	}
	p.items = exp.Items
	return nil
}

// Add adds a new with the path to the Playlist.
func (p *Playlist) Add(path index.Path) {
	p.items = append(p.items, newItem(path))
}

// Remove removes the item with index `n` and path `path` from the Playlist.
func (p *Playlist) Remove(n int, path index.Path) error {
	if n >= len(p.items) {
		return fmt.Errorf("invalid item index (items: %d): %d", len(p.items), n)
	}

	if path.Equal(p.items[n].path) {
		if n > 0 {
			p.items = append(p.items[:n-1], p.items[n-1:]...)
		} else {
			p.items = p.items[1:]
		}
		return nil
	}
	p.items[n].AddTransform(RemovePath(path))
	return nil
}

// Items returns a slice of *Item instances which represent each item in the playlist.
func (p *Playlist) Items() []*Item {
	items := make([]*Item, len(p.items))
	for i, item := range p.items {
		items[i] = item
	}
	return items
}

// Paths returns the list of paths for the tracks within the Item, using Collection
// as the data source.
func Paths(item *Item, c index.Collection) ([]index.Path, error) {
	g, err := index.GroupFromPath(c, item.path[1:]) // Trim "Root" prefix
	if err != nil {
		return nil, err
	}

	removePaths := make([]index.Path, len(item.transforms))
	for _, transform := range item.transforms {
		if path, ok := transform.(RemovePath); ok {
			removePaths = append(removePaths, index.Path(path))
		}
	}

	var paths []index.Path
	walkFn := func(t index.Track, p index.Path) error {
		for _, rp := range removePaths {
			if rp.Contains(p) {
				return nil
			}
		}
		paths = append(paths, p)
		return nil
	}
	index.Walk(g, item.path, walkFn)
	return paths, nil
}
