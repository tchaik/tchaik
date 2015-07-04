// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

// Action is a type which represents an enumeration of available player actions.
type Action string

// Actions which don't need values.
const (
	ActionPlay            Action = "play"
	ActionPause                  = "pause"
	ActionNext                   = "next"
	ActionPrev                   = "prev"
	ActionTogglePlayPause        = "togglePlayPause"
	ActionToggleMute             = "toggleMute"
)

// Player actions which require values.
const (
	ActionSetVolume Action = "setVolume"
	ActionSetMute          = "setMute"
	ActionSetTime          = "setTime"
)

// Player is an interface which defines methods for controlling a player.
type Player interface {
	// Key returns a unique identifier for this Player.
	Key() string

	// Do sends a Action to the player.
	Do(a Action) error

	// SetMute enabled/disables mute.
	SetMute(bool) error
	// SetVolume sets the volume (value should be between 0.0 and 1.0).
	SetVolume(float64) error
	// SetTime sets the current play position
	SetTime(float64) error
}

type multiPlayer struct {
	key     string
	players []Player
}

// MultiPlayer returns a player that will apply calls to all provided Players
// in sequence.  If an error is returning by a Player then it is returned
// immediately.
func MultiPlayer(key string, players ...Player) Player {
	return multiPlayer{
		key:     key,
		players: players,
	}
}

type setFloatFn func(Player, float64) error

func (m multiPlayer) applySetFloatFn(fn setFloatFn, f float64) error {
	for _, p := range m.players {
		err := fn(p, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m multiPlayer) Key() string { return m.key }

func (m multiPlayer) Do(a Action) error {
	for _, p := range m.players {
		err := p.Do(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m multiPlayer) SetVolume(f float64) error { return m.applySetFloatFn(Player.SetVolume, f) }
func (m multiPlayer) SetTime(f float64) error   { return m.applySetFloatFn(Player.SetTime, f) }

func (m multiPlayer) SetMute(b bool) error {
	for _, p := range m.players {
		err := p.SetMute(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m multiPlayer) MarshalJSON() ([]byte, error) {
	playerKeys := make([]string, len(m.players))
	for i, p := range m.players {
		playerKeys[i] = p.Key()
	}

	rep := struct {
		Key        string   `json:"key"`
		PlayerKeys []string `json:"playerKeys"`
	}{
		m.key,
		playerKeys,
	}
	return json.Marshal(rep)
}

// Validated wraps a player with validation checks for value-setting methods.
func ValidatedPlayer(p Player) Player {
	return validatedPlayer{
		Player: p,
	}
}

type validatedPlayer struct {
	Player
}

// InvalidValueError is an error returned by value-setting methods.
type InvalidValueError string

// Error implements error.
func (v InvalidValueError) Error() string { return string(v) }

// SetVolume implements Player.
func (v validatedPlayer) SetVolume(f float64) error {
	if f < 0.0 || f > 1.0 {
		return InvalidValueError(fmt.Sprintf("invalid volume value '%v': must be between 0.0 and 1.0", f))
	}
	return v.Player.SetVolume(f)
}

// SetTime implements Player.
func (v validatedPlayer) SetTime(f float64) error {
	if f < 0.0 {
		return InvalidValueError(fmt.Sprintf("invalid time value '%v': must be greater than 0.0", f))
	}
	return v.Player.SetTime(f)
}

func (v validatedPlayer) MarshalJSON() ([]byte, error) {
	if m, ok := v.Player.(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	return nil, nil
}

// WebsocketPlayer returns a player which will transmit commands on the websocket connection.
func WebsocketPlayer(key string, ws *websocket.Conn) Player {
	return &websocketPlayer{
		Conn: ws,
		key:  key,
	}
}

type websocketPlayer struct {
	*websocket.Conn
	key string
}

func (w *websocketPlayer) sendAction(data interface{}) error {
	return websocket.JSON.Send(w.Conn, struct {
		Action string
		Data   interface{}
	}{
		Action: CtrlAction,
		Data:   data,
	})
}

func (w *websocketPlayer) sendCtrlValue(key string, value interface{}) error {
	return w.sendAction(struct {
		Key   string
		Value interface{}
	}{
		Key:   key,
		Value: value,
	})
}

func (w websocketPlayer) Key() string { return w.key }

var websocketActions = map[Action]string{
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

func websocketToAction(a string) (Action, bool) {
	for k, v := range websocketActions {
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

func (w websocketPlayer) Do(a Action) error {
	s, ok := websocketActions[a]
	if !ok {
		return InvalidActionError(a)
	}
	return w.sendAction(s)
}

func (w websocketPlayer) SetMute(b bool) error      { return w.sendCtrlValue("mute", b) }
func (w websocketPlayer) SetVolume(f float64) error { return w.sendCtrlValue("volume", f) }
func (w websocketPlayer) SetTime(f float64) error   { return w.sendCtrlValue("time", f) }

func (w websocketPlayer) MarshalJSON() ([]byte, error) {
	rep := struct {
		Key string `json:"key"`
	}{
		Key: w.key,
	}
	return json.Marshal(rep)
}
