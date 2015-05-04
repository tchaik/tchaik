// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "golang.org/x/net/websocket"

// Player is an interface which defines methods for controlling a player.
type Player interface {
	// Play the current track.
	Play() error
	// Pause the current track.
	Pause() error
	// NextTrack jumps to the next track.
	NextTrack() error
	// PreviousTrack jumps to the previous track.
	PreviousTrack() error

	// Toggle play/pause.
	TogglePlayPause() error
	// Toggle mute on/off.
	ToggleMute() error

	// SetMute enabled/disables mute.
	SetMute(bool) error
	// SetVolume sets the volume (value should be between 0.0 and 1.0).
	SetVolume(float64) error
	// SetTime sets the current play position
	SetTime(float64) error
}

type multiPlayer struct {
	players []Player
}

// MultiPlayer returns a player that will apply calls to all provided Players
// in sequence.  If an error is returning by a Player then it is returned
// immediately.
func MultiPlayer(players ...Player) Player {
	return multiPlayer{
		players: players,
	}
}

type cmdFn func(Player) error

func (m multiPlayer) applyCmdFn(fn cmdFn) error {
	for _, p := range m.players {
		err := fn(p)
		if err != nil {
			return err
		}
	}
	return nil
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

func (m multiPlayer) Play() error            { return m.applyCmdFn(Player.Play) }
func (m multiPlayer) Pause() error           { return m.applyCmdFn(Player.Pause) }
func (m multiPlayer) NextTrack() error       { return m.applyCmdFn(Player.NextTrack) }
func (m multiPlayer) PreviousTrack() error   { return m.applyCmdFn(Player.PreviousTrack) }
func (m multiPlayer) TogglePlayPause() error { return m.applyCmdFn(Player.TogglePlayPause) }
func (m multiPlayer) ToggleMute() error      { return m.applyCmdFn(Player.ToggleMute) }

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

// Validated wraps a player with validation checks for value-setting methods.
func ValidatedPlayer(p Player) Player {
	return validator{
		Player: p,
	}
}

type validator struct {
	Player
}

// InvalidValueError is an error returned by value-setting methods.
type InvalidValueError string

// Error implements error.
func (v InvalidValueError) Error() string { return string(v) }

// SetVolume implements Player.
func (v validator) SetVolume(f float64) error {
	if f < 0.0 || f > 1.0 {
		return InvalidValueError("invalid volume value: must be between 0.0 and 1.0")
	}
	return v.Player.SetVolume(f)
}

// SetTime implements Player.
func (v validator) SetTime(f float64) error {
	if f < 0.0 {
		return InvalidValueError("invalid time value: must be greater than 0.0")
	}
	return v.Player.SetTime(f)
}

type websocketPlayer struct {
	*websocket.Conn
}

func (w *websocketPlayer) sendCtrlAction(data interface{}) error {
	return websocket.JSON.Send(w.Conn, struct {
		Action string
		Data   interface{}
	}{
		Action: CtrlAction,
		Data:   data,
	})
}

func (w *websocketPlayer) sendCtrlValue(key string, value interface{}) error {
	return w.sendCtrlAction(struct {
		Key   string
		Value interface{}
	}{
		Key:   key,
		Value: value,
	})
}

func (w websocketPlayer) Play() error            { return w.sendCtrlAction("PLAY") }
func (w websocketPlayer) Pause() error           { return w.sendCtrlAction("PAUSE") }
func (w websocketPlayer) NextTrack() error       { return w.sendCtrlAction("NEXT") }
func (w websocketPlayer) PreviousTrack() error   { return w.sendCtrlAction("PREV") }
func (w websocketPlayer) TogglePlayPause() error { return w.sendCtrlAction("TOGGLE_PLAY_PAUSE") }
func (w websocketPlayer) ToggleMute() error      { return w.sendCtrlAction("TOGGLE_MUTE") }

func (w websocketPlayer) SetMute(b bool) error      { return w.sendCtrlValue("mute", b) }
func (w websocketPlayer) SetVolume(f float64) error { return w.sendCtrlValue("volume", f) }
func (w websocketPlayer) SetTime(f float64) error   { return w.sendCtrlValue("time", f) }
