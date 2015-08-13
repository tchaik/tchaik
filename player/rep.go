// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package player

import (
	"encoding/json"
	"fmt"
)

// RepFn is a type which represents a command function which will be called
// with a standardised command structure when Player methods are used.
type RepFn func(interface{})

// RepAction is a representation of a player action as it would be transmitted.
type RepAction struct {
	Action string      `json:"action"`
	Value  interface{} `json:",omitempty"`
}

// Apply applies the action in RepAction to the Player.
func (r RepAction) Apply(p Player) (err error) {
	a := Action(r.Action)
	switch a {
	case ActionPlay, ActionPause, ActionNext, ActionPrev, ActionTogglePlayPause, ActionToggleMute, ActionToggleRepeat:
		err = p.Do(a)

	case ActionSetVolume, ActionSetMute, ActionSetTime:
		if r.Value == nil {
			err = InvalidValueError("value required")
			break
		}

		switch a {
		case ActionSetVolume:
			f, ok := r.Value.(float64)
			if !ok {
				err = InvalidValueError("invalid volume value: expected float")
				break
			}
			err = p.SetVolume(f)

		case ActionSetMute:
			b, ok := r.Value.(bool)
			if !ok {
				err = InvalidValueError("invalid mute value: expected boolean")
				break
			}
			err = p.SetMute(b)

		case ActionSetTime:
			f, ok := r.Value.(float64)
			if !ok {
				err = InvalidValueError("invalid time value: expected float")
				break
			}
			err = p.SetTime(f)
		}

	default:
		err = InvalidActionError(a)
	}
	return
}

type rep struct {
	key string
	fn  RepFn
}

// NewRep creates a Player p which will call fn with a standardised representation for each
// call to methods of p (eee RepActions for the Action mapping).
// FIXME: this is possibly a unhelpful and overly generic mess.
func NewRep(key string, fn RepFn) Player {
	return &rep{
		key: key,
		fn:  fn,
	}
}

func (r rep) sendAction(action string) error {
	r.fn(RepAction{
		Action: action,
	})
	return nil
}

func (r rep) sendActionValue(action string, value interface{}) error {
	r.fn(RepAction{
		Action: action,
		Value:  value,
	})
	return nil
}

func (r rep) Key() string { return r.key }

var RepActions = map[Action]string{
	ActionPlay:            "PLAY",
	ActionPause:           "PAUSE",
	ActionNext:            "NEXT",
	ActionPrev:            "PREV",
	ActionTogglePlayPause: "TOGGLE_PLAY_PAUSE",
	ActionToggleRepeat:    "TOGGLE_REPEAT",
	ActionToggleMute:      "TOGGLE_MUTE",

	ActionSetVolume: "SET_VOLUME",
	ActionSetMute:   "SET_MUTE",
	ActionSetTime:   "SET_TIME",
}

func RepActionToAction(a string) (Action, bool) {
	for k, v := range RepActions {
		if v == a {
			return k, true
		}
	}
	return Action(""), false
}

type InvalidActionError string

func (i InvalidActionError) Error() string {
	return fmt.Sprintf("invalid player action: '%s'", string(i))
}

func (r rep) Do(a Action) error {
	s, ok := RepActions[a]
	if !ok {
		return InvalidActionError(a)
	}
	return r.sendAction(s)
}

func (r rep) SetMute(b bool) error      { return r.sendActionValue("mute", b) }
func (r rep) SetVolume(f float64) error { return r.sendActionValue("volume", f) }
func (r rep) SetTime(f float64) error   { return r.sendActionValue("time", f) }

func (r rep) MarshalJSON() ([]byte, error) {
	rep := struct {
		Key string `json:"key"`
	}{
		Key: r.key,
	}
	return json.Marshal(rep)
}
