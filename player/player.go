// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package player defines types and methods for defining a Tchaik player.
package player

import (
	"encoding/json"
	"fmt"
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
func Validated(p Player) Player {
	return validated{
		Player: p,
	}
}

type validated struct {
	Player
}

// InvalidValueError is an error returned by value-setting methods.
type InvalidValueError string

// Error implements error.
func (v InvalidValueError) Error() string { return string(v) }

// SetVolume implements Player.
func (v validated) SetVolume(f float64) error {
	if f < 0.0 || f > 1.0 {
		return InvalidValueError(fmt.Sprintf("invalid volume value '%v': must be between 0.0 and 1.0", f))
	}
	return v.Player.SetVolume(f)
}

// SetTime implements Player.
func (v validated) SetTime(f float64) error {
	if f < 0.0 {
		return InvalidValueError(fmt.Sprintf("invalid time value '%v': must be greater than 0.0", f))
	}
	return v.Player.SetTime(f)
}

func (v validated) MarshalJSON() ([]byte, error) {
	if m, ok := v.Player.(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	return nil, nil
}

// Players is a collection of players which are identified by key.
type Players struct {
	sync.RWMutex
	m map[string]Player
}

// NewPlayers creates a Players.
func NewPlayers() *Players {
	return &Players{m: make(map[string]Player)}
}

// Add the Player to the Players.
func (s *Players) Add(p Player) {
	s.Lock()
	defer s.Unlock()

	s.m[p.Key()] = p
}

// Remove the Player from Players (by key).
func (s *Players) Remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)
}

// Get the Player identified by the key.
func (s *Players) Get(key string) Player {
	s.RLock()
	defer s.RUnlock()

	return s.m[key]
}

// List all Player keys in Players.
func (s *Players) List() []string {
	s.RLock()
	defer s.RUnlock()

	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

// MarshalJSON implements json.Marshaler
func (s *Players) MarshalJSON() ([]byte, error) {
	keys := s.List()
	return json.Marshal(struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	})
}
