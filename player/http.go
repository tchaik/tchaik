// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package player

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

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

// NewHTTPHandler returns an http.Handler which defines a REST API for interacting with
// Players.
func NewHTTPHandler(p *Players) http.Handler {
	return &httpHandler{
		players: p,
	}
}

type httpHandler struct {
	players *Players
}

// ServeHTTP implements http.Handler.
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		switch r.Method {
		case "GET":
			h.writeJSON(w, r, h.players)
		case "POST":
			h.createPlayer(w, r)
		}
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := h.players.Get(paths[0])
	if p == nil {
		http.Error(w, "invalid player key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "DELETE":
		h.players.Remove(paths[0])
		w.WriteHeader(http.StatusNoContent)
		return

	case "PUT":
		h.playerAction(p, w, r)

	case "GET":
		h.writeJSON(w, r, p)
	}
}

func (h *httpHandler) writeJSON(w http.ResponseWriter, r *http.Request, x interface{}) {
	b, err := json.Marshal(x)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func (h *httpHandler) createPlayer(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	postData := struct {
		Key        string
		PlayerKeys []string
	}{}
	err := dec.Decode(&postData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	if p := h.players.Get(postData.Key); p != nil {
		http.Error(w, "player key already exists", http.StatusBadRequest)
		return
	}

	if postData.PlayerKeys == nil || len(postData.PlayerKeys) == 0 {
		http.Error(w, "no player keys specified", http.StatusBadRequest)
		return
	}

	var players []Player
	for _, pk := range postData.PlayerKeys {
		p := h.players.Get(pk)
		if p == nil {
			http.Error(w, fmt.Sprintf("invalid player key: %v", pk), http.StatusBadRequest)
			return
		}
		players = append(players, p)
	}
	h.players.Add(MultiPlayer(postData.Key, players...))
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) playerAction(p Player, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	data := struct {
		Action string
		Value  interface{}
	}{}
	err := dec.Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	a := Action(data.Action)
	switch a {
	case ActionPlay, ActionPause, ActionNext, ActionPrev, ActionTogglePlayPause, ActionToggleMute:
		err = p.Do(a)

	case ActionSetVolume, ActionSetMute, ActionSetTime:
		if data.Value == nil {
			err = InvalidValueError("value required")
			break
		}

		switch a {
		case ActionSetVolume:
			f, ok := data.Value.(float64)
			if !ok {
				err = InvalidValueError("invalid volume value: expected float")
				break
			}
			err = p.SetVolume(f)

		case ActionSetMute:
			b, ok := data.Value.(bool)
			if !ok {
				err = InvalidValueError("invalid mute value: expected boolean")
				break
			}
			err = p.SetMute(b)

		case ActionSetTime:
			f, ok := data.Value.(float64)
			if !ok {
				err = InvalidValueError("invalid time value: expected float")
				break
			}
			err = p.SetTime(f)
		}

	default:
		err = InvalidActionError(a)
	}

	if err != nil {
		if err, ok := err.(InvalidActionError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err, ok := err.(InvalidValueError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("error sending player command: %v", err), http.StatusInternalServerError)
	}
}
