// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"tchaik.com/index"
	"tchaik.com/index/checklist"
	"tchaik.com/index/cursor"
	"tchaik.com/index/favourite"
	"tchaik.com/index/history"
	"tchaik.com/index/playlist"
)

// Meta is a container for extra metadata which wraps the central media library.
// TODO: Refactor meta to be user-specific.
type Meta struct {
	history    history.Store
	favourites favourite.Store
	checklist  checklist.Store
	playlists  playlist.Store
	cursors    cursor.Store
}

type metaFieldGrp struct {
	index.Group

	field string
	value interface{}
}

func (mfg metaFieldGrp) Field(f string) interface{} {
	if f == mfg.field {
		return mfg.value
	}
	return mfg.Group.Field(f)
}

type metaFieldCol struct {
	index.Collection

	field string
	value interface{}
}

func (mfc metaFieldCol) Field(f string) interface{} {
	if f == mfc.field {
		return mfc.value
	}
	return mfc.Collection.Field(f)
}

func newMetaField(g index.Group, field string, value bool) index.Group {
	if !value {
		return g
	}
	if c, ok := g.(index.Collection); ok {
		// TODO(dhowden): currently need to maintain the underlying interface type
		// so that it can be correctly transmitted, need a better way of doing this.
		return metaFieldCol{
			Collection: c,
			field:      field,
			value:      value,
		}
	}
	return metaFieldGrp{
		Group: g,
		field: field,
		value: value,
	}
}

func (m *Meta) annotateFavourites(p index.Path, g index.Group) index.Group {
	return newMetaField(g, "Favourite", m.favourites.Get(p))
}

func (m *Meta) annotateChecklist(p index.Path, g index.Group) index.Group {
	return newMetaField(g, "Checklist", m.checklist.Get(p))
}
