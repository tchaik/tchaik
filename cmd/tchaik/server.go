// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"path"

	"github.com/dhowden/httpauth"
	"github.com/tchaik/tchaik/store"
)

type fsServeMux struct {
	httpauth.ServeMux
}

func (fsm *fsServeMux) HandleFileSystem(pattern string, fs http.FileSystem) {
	fsm.ServeMux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(fs)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, path.Join(staticDir, "index.html"))
}

type server struct {
	fsServeMux
	lib     Library
	players *players
}

func newServer(l Library, mediaFileSystem, artworkFileSystem http.FileSystem) *server {
	mediaFileSystem = l.FileSystem(mediaFileSystem)
	artworkFileSystem = l.FileSystem(artworkFileSystem)

	var c httpauth.Checker = httpauth.None{}
	if auth {
		c = creds
	}

	s := &server{
		fsServeMux: fsServeMux{httpauth.NewServeMux(c, http.NewServeMux())},
		lib:        l,
		players:    newPlayers(),
	}

	s.HandleFunc("/", rootHandler)
	s.HandleFileSystem("/static/", http.Dir(staticDir))
	s.HandleFileSystem("/track/", mediaFileSystem)
	s.HandleFileSystem("/artwork/", artworkFileSystem)
	s.HandleFileSystem("/icon/", store.FaviconFileSystem(artworkFileSystem))
	s.Handle("/socket", s.WebsocketHandler())
	s.Handle("/api/players/", http.StripPrefix("/api/players/", &playersHandler{s.players}))

	return s
}
