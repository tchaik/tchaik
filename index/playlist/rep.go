// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playlist

import (
	"fmt"

	"tchaik.com/index"
)

// Action is a type used to represent a playlist mutation action.
type Action string

const (
	ActionCreate Action = "create"
	ActionDelete        = "delete"

	ActionAddItem    = "addItem"
	ActionRemoveItem = "deleteItem"
)

var actionToAction = map[string]Action{
	"ADD_ITEM": ActionAddItem,
	"REMOVE":   ActionRemoveItem,
}

type RepAction struct {
	Name   string     `json:"name"`
	Action Action     `json:"action"`
	Path   index.Path `json:"path"`
	Index  int        `json:"index"`
}

func (a RepAction) Apply(s Store) error {
	if a.Action == ActionCreate {
		s.Set(a.Name, &Playlist{})
		return nil
	}

	action, ok := actionToAction[string(a.Action)]
	if !ok {
		return fmt.Errorf("unknown action: %v", a.Action)
	}

	p := s.Get(a.Name)
	if p == nil {
		return fmt.Errorf("invalid playlist name: '%v'", a.Name)
	}

	switch action {
	case ActionDelete:
		s.Delete(a.Name)
	case ActionAddItem:
		p.Add(a.Path)
	case ActionRemoveItem:
		p.Remove(a.Index, a.Path)
	}

	s.Set(a.Name, p)
	return nil
}
