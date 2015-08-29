// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"tchaik.com/index"
	"tchaik.com/index/checklist"
	"tchaik.com/index/favourite"
	"tchaik.com/index/history"
)

// Meta is a container for extra metadata which wraps the central media library.
// TODO: Refactor meta to be user-specific.
type Meta struct {
	history    history.Store
	favourites favourite.Store
	checklist  checklist.Store
}

func (m *Meta) annotateFavourites(p index.Path, g group) group {
	if len(g.Tracks) > 0 || len(g.Groups) > 0 {
		g.Favourite = m.favourites.Get(p)
	}
	return g
}

func (m *Meta) annotateChecklist(p index.Path, g group) group {
	if len(g.Tracks) > 0 || len(g.Groups) > 0 {
		g.Checklist = m.checklist.Get(p)
	}
	return g
}
