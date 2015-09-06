// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cursor

import (
	"fmt"

	"tchaik.com/index"
	"tchaik.com/index/playlist"
)

type Action string

const (
	ActionSet      Action = "set"
	ActionNext            = "next"
	ActionPrevious        = "previous"
)

type RepAction struct {
	Name   string     `json:"name"`
	Action Action     `json:"action"`
	Path   index.Path `json:"path"`
	Index  int        `json:index`
}

var actionToAction = map[string]Action{
	"SET":  ActionSet,
	"NEXT": ActionNext,
	"PREV": ActionPrevious,
}

func (a RepAction) Apply(s Store, ps playlist.Store, collection index.Collection) error {
	action, ok := actionToAction[string(a.Action)]
	if !ok {
		return fmt.Errorf("unknown action: %v", a.Action)
	}

	if action == ActionSet {
		p := ps.Get(a.Name)
		if p == nil {
			return fmt.Errorf("cannot set cursor for invalid playlist name: %v", a.Name)
		}

		c := NewCursor(p, collection)
		c.Set(a.Index, a.Path)
		return s.Set(a.Name, c)
	}

	c := s.Get(a.Name)
	if c == nil {
		return fmt.Errorf("invalid cursor name: %v", a.Name)
	}

	var err error
	switch action {
	case ActionPrevious:
		err = c.Backward()
	case ActionNext:
		err = c.Forward()
	}
	err1 := s.Set(a.Name, c)
	if err == nil {
		err = err1
	}
	return err
}
