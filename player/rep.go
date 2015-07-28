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

func (r rep) sendAction(data interface{}) error {
	r.fn(data)
	return nil
}

func (r rep) sendCtrlValue(key string, value interface{}) error {
	r.sendAction(struct {
		Key   string
		Value interface{}
	}{
		Key:   key,
		Value: value,
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

func (r rep) SetMute(b bool) error      { return r.sendCtrlValue("mute", b) }
func (r rep) SetVolume(f float64) error { return r.sendCtrlValue("volume", f) }
func (r rep) SetTime(f float64) error   { return r.sendCtrlValue("time", f) }

func (r rep) MarshalJSON() ([]byte, error) {
	rep := struct {
		Key string `json:"key"`
	}{
		Key: r.key,
	}
	return json.Marshal(rep)
}
