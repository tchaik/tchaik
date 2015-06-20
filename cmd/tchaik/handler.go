// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"path"

	"github.com/dhowden/httpauth"
	"github.com/tchaik/tchaik/index/history"
	"github.com/tchaik/tchaik/store"
)

type fsServeMux struct {
	httpauth.ServeMux
}

// HandleFileSystem is a convenience method for adding an http.FileServer handler to an
// http.ServeMux.
func (fsm *fsServeMux) HandleFileSystem(pattern string, fs http.FileSystem) {
	fsm.ServeMux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(fs)))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, path.Join(staticDir, "index.html"))
}

// NewHandler creates the root http.Handler.
func NewHandler(l Library, hs history.Store, mediaFileSystem, artworkFileSystem http.FileSystem) http.Handler {
	var c httpauth.Checker = httpauth.None{}
	if auth {
		c = creds
	}
	h := fsServeMux{
		httpauth.NewServeMux(c, http.NewServeMux()),
	}

	h.HandleFunc("/", rootHandler)
	h.HandleFileSystem("/static/", http.Dir(staticDir))

	mediaFileSystem = l.FileSystem(mediaFileSystem)
	artworkFileSystem = l.FileSystem(artworkFileSystem)
	h.HandleFileSystem("/track/", mediaFileSystem)
	h.HandleFileSystem("/artwork/", artworkFileSystem)
	h.HandleFileSystem("/icon/", store.FaviconFileSystem(artworkFileSystem))

	p := newPlayers()
	h.Handle("/socket", NewWebsocketHandler(l, p, hs))
	h.Handle("/api/players/", http.StripPrefix("/api/players/", &playersHandler{p}))

	return h
}
