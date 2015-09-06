// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cursor defines types for constructing play cursors.
package cursor

import (
	"fmt"

	"tchaik.com/index"
	"tchaik.com/index/playlist"
)

// Position is a type which stores a position in a playlist.
type Position struct {
	Path  index.Path `json:"path"`
	Index int        `json:"index"`
}

// Cursor is a moveable marker on a playlist.
type Cursor struct {
	Current  Position `json:"current"`
	Next     Position `json:"next"`
	Previous Position `json:"previous"`

	p *playlist.Playlist
	c index.Collection
}

// NewCursor creates a new Cursor for the playlist.Playlist using the index.Collection
// as the source of tracks.
func NewCursor(p *playlist.Playlist, c index.Collection) *Cursor {
	return &Cursor{
		p: p,
		c: c,
	}
}

// Set sets the value of the playlist cursor to the current position and index.Path.
func (c *Cursor) Set(i int, p index.Path) {
	c.Current = Position{Index: i, Path: p}
	c.Next, _ = c.next(c.Current)
	c.Previous, _ = c.prev(c.Current)
}

// Forward moves the cursor forwards.  Returns an error if the next track could not be found,
// and sets the Next item to be empty.
func (c *Cursor) Forward() (err error) {
	c.Previous = c.Current
	c.Current = c.Next
	c.Next, err = c.next(c.Current)
	return
}

// Backward moves the cursor backwards.  Returns an error if the previous track could not be found,
// and sets the Previous item to be empty.
func (c *Cursor) Backward() (err error) {
	c.Next = c.Current
	c.Current = c.Previous
	c.Previous, err = c.prev(c.Current)
	return
}

func (c *Cursor) paths(n int) ([]index.Path, error) {
	items := c.p.Items()
	item := items[n]
	return playlist.Paths(item, c.c)
}

func indexOfPath(paths []index.Path, p index.Path) int {
	n := -1
	for i, x := range paths {
		if p.Equal(x) {
			n = i
			break
		}
	}
	return n
}

func (c *Cursor) pathIndex(p Position) ([]index.Path, int, error) {
	paths, err := c.paths(p.Index)
	if err != nil {
		return nil, 0, err
	}
	return paths, indexOfPath(paths, p.Path), nil
}

func (c *Cursor) next(p Position) (Position, error) {
	paths, i, err := c.pathIndex(p)
	if err != nil {
		return Position{}, err
	}

	if i == -1 {
		return Position{}, fmt.Errorf("didn't find path: %v", p.Path)
	}

	if i < len(paths)-1 {
		return Position{
			Path:  paths[i+1],
			Index: p.Index,
		}, nil
	}

	if p.Index < len(c.p.Items())-1 {
		paths, err = c.paths(p.Index + 1)
		if err != nil {
			return Position{}, err
		}

		return Position{
			Path:  paths[0],
			Index: p.Index + 1,
		}, nil
	}
	return Position{}, nil
}

func (c *Cursor) prev(p Position) (Position, error) {
	paths, i, err := c.pathIndex(p)
	if err != nil {
		return Position{}, err
	}

	if i == -1 {
		return Position{}, fmt.Errorf("didn't find path: %v", p.Path)
	}

	if i > 0 {
		return Position{
			Path:  paths[i-1],
			Index: p.Index,
		}, nil
	}

	if p.Index > 0 {
		paths, err = c.paths(p.Index - 1)
		if err != nil {
			return Position{}, err
		}

		return Position{
			Path:  paths[len(paths)-1],
			Index: p.Index - 1,
		}, nil
	}
	return Position{}, nil
}
