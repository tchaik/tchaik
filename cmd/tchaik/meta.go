// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

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

func loadLocalMeta() (*Meta, error) {
	fmt.Printf("Loading play history...")
	playHistoryStore, err := history.NewStore(playHistoryPath)
	if err != nil {
		return nil, fmt.Errorf("error loading play history: %v", err)
	}
	fmt.Println("done.")

	fmt.Printf("Loading favourites...")
	favouriteStore, err := favourite.NewStore(favouritesPath)
	if err != nil {
		return nil, fmt.Errorf("\nerror loading favourites: %v", err)
	}
	fmt.Println("done.")

	fmt.Printf("Loading checklist...")
	checklistStore, err := checklist.NewStore(checklistPath)
	if err != nil {
		return nil, fmt.Errorf("\nerror loading checklist: %v", err)
	}
	fmt.Println("done.")

	fmt.Printf("Loading playlists...")
	playlistStore, err := playlist.NewStore(playlistPath)
	if err != nil {
		return nil, fmt.Errorf("\nerror loading playlists: %v", err)
	}
	// TODO(dhowden): remove this once we can better intialise the "Default" playlist
	if p := playlistStore.Get("Default"); p == nil {
		playlistStore.Set("Default", &playlist.Playlist{})
	}
	fmt.Println("done")

	fmt.Printf("Loading cursors...")
	cursorStore, err := cursor.NewStore(cursorPath)
	if err != nil {
		return nil, fmt.Errorf("\nerror loading cursor: %v", err)
	}
	fmt.Println("done")

	return &Meta{
		history:    playHistoryStore,
		favourites: favouriteStore,
		checklist:  checklistStore,
		playlists:  playlistStore,
		cursors:    cursorStore,
	}, nil
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

// Annotate adds any meta information to the Group (identified by Path).
func (m *Meta) Annotate(p index.Path, g index.Group) index.Group {
	g = newMetaField(g, "Favourite", m.favourites.Get(p))
	return newMetaField(g, "Checklist", m.checklist.Get(p))
}
