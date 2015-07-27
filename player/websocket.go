// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package player

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

// TODO: Remove this.
const (
	CtrlAction string = "CTRL"
)

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

func WebsocketToAction(a string) (Action, bool) {
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
