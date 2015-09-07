// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playlist

import (
	"testing"

	"tchaik.com/index"
)

func TestPlaylistAdd(t *testing.T) {
	pathA := index.NewPath("Root:a")
	pathB := index.NewPath("Root:b")

	p := &Playlist{}
	p.Add(pathA)
	if len(p.Items()) != 1 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 1)
	}

	p.Add(pathB)
	if len(p.Items()) != 2 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 2)
	}
}

func TestPlaylistRemoveItem(t *testing.T) {
	pathA := index.NewPath("Root:a")

	p := &Playlist{}
	p.Add(pathA)
	if len(p.Items()) != 1 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 1)
	}

	err := p.Remove(0, pathA)
	if err != nil {
		t.Errorf("unexpected error removing item: %v", err)
	}

	if len(p.Items()) != 0 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 0)
	}
}

func TestPlaylistRemoveMultipleItems(t *testing.T) {
	pathA := index.NewPath("Root:a")
	pathB := index.NewPath("Root:b")

	p := &Playlist{}
	p.Add(pathA)
	p.Add(pathB)

	err := p.Remove(1, pathB)
	if err != nil {
		t.Errorf("unexpected error removing item: %v", err)
	}

	if len(p.Items()) != 1 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 1)
	}

	err = p.Remove(0, pathA)
	if err != nil {
		t.Errorf("unexpected error removing item: %v", err)
	}

	if len(p.Items()) != 0 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 0)
	}
}

func TestPlaylistItemRemove(t *testing.T) {
	pathA := index.NewPath("Root:a")
	subPathA := index.NewPath("Root:a:b")
	pathB := index.NewPath("Root:b")

	p := &Playlist{}
	p.Add(pathA)
	err := p.Remove(0, subPathA)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(p.Items()) != 1 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 1)
	}

	p.Add(pathB)
	err = p.Remove(1, pathB)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(p.Items()) != 1 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 1)
	}

	err = p.Remove(0, pathA)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(p.Items()) != 0 {
		t.Errorf("len(p.Items()) = %d, expected: %d", len(p.Items()), 0)
	}

	err = p.Remove(0, pathA)
	if err == nil {
		t.Errorf("expected error for removing invalid item (items: %v)", p.Items())
	}
}
