// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type players struct {
	sync.RWMutex
	m map[string]Player
}

func newPlayers() *players {
	return &players{m: make(map[string]Player)}
}

func (s *players) add(p Player) {
	s.Lock()
	defer s.Unlock()

	s.m[p.Key()] = p
}

func (s *players) remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)
}

func (s *players) get(key string) Player {
	s.RLock()
	defer s.RUnlock()

	return s.m[key]
}

func (s *players) list() []string {
	s.RLock()
	defer s.RUnlock()

	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *players) MarshalJSON() ([]byte, error) {
	keys := s.list()
	return json.Marshal(struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	})
}

type playersHandler struct {
	players *players
}

func (h *playersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	p := h.players.get(paths[0])
	if p == nil {
		http.Error(w, "invalid player key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "DELETE":
		h.players.remove(paths[0])
		w.WriteHeader(http.StatusNoContent)
		return

	case "PUT":
		h.playerAction(p, w, r)

	case "GET":
		h.writeJSON(w, r, p)
	}
}

func (playersHandler) writeJSON(w http.ResponseWriter, r *http.Request, x interface{}) {
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

func (h *playersHandler) createPlayer(w http.ResponseWriter, r *http.Request) {
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

	if p := h.players.get(postData.Key); p != nil {
		http.Error(w, "player key already exists", http.StatusBadRequest)
		return
	}

	if postData.PlayerKeys == nil || len(postData.PlayerKeys) == 0 {
		http.Error(w, "no player keys specified", http.StatusBadRequest)
		return
	}

	var players []Player
	for _, pk := range postData.PlayerKeys {
		p := h.players.get(pk)
		if p == nil {
			http.Error(w, fmt.Sprintf("invalid player key: %v", pk), http.StatusBadRequest)
			return
		}
		players = append(players, p)
	}
	h.players.add(MultiPlayer(postData.Key, players...))
	w.WriteHeader(http.StatusCreated)
}

func (playersHandler) playerAction(p Player, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	putData := struct {
		Action string
		Value  interface{}
	}{}
	err := dec.Decode(&putData)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	switch putData.Action {
	case "play":
		err = p.Play()

	case "pause":
		err = p.Pause()

	case "next":
		err = p.NextTrack()

	case "prev":
		err = p.PreviousTrack()

	case "togglePlayPause":
		err = p.TogglePlayPause()

	case "toggleMute":
		err = p.ToggleMute()

	case "setVolume":
		f, ok := putData.Value.(float64)
		if !ok {
			err = InvalidValueError("invalid volume value: expected float")
			break
		}
		err = p.SetVolume(f)

	case "setMute":
		b, ok := putData.Value.(bool)
		if !ok {
			err = InvalidValueError("invalid mute value: expected boolean")
			break
		}
		err = p.SetMute(b)

	case "setTime":
		f, ok := putData.Value.(float64)
		if !ok {
			err = InvalidValueError("invalid time value: expected float")
			break
		}
		err = p.SetTime(f)

	default:
		err = InvalidValueError("invalid action")
	}

	if err != nil {
		if err, ok := err.(InvalidValueError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("error sending player command: %v", err), http.StatusInternalServerError)
	}
}
