// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"github.com/dhowden/httpauth"

	"github.com/tchaik/tchaik/index/history"
	"github.com/tchaik/tchaik/store"
)

// traceFS is a type which implements http.FileSystem and is used at the top-level to
// intialise a trace which can be passed through to FileSystem implementations.
type traceFS struct {
	store.FileSystem
	family string
}

// Open implements http.FileSystem.
func (t *traceFS) Open(path string) (http.File, error) {
	tr := trace.New(t.family, path)
	ctx := trace.NewContext(context.Background(), tr)
	f, err := t.FileSystem.Open(ctx, path)

	// TODO: Decide where this should be in general (requests can be on-going).
	tr.Finish()
	return f, err
}

type fsServeMux struct {
	httpauth.ServeMux
}

// HandleFileSystem is a convenience method for adding an http.FileServer handler to an
// http.ServeMux.
func (fsm *fsServeMux) HandleFileSystem(pattern string, fs store.FileSystem) {
	fsm.ServeMux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(&traceFS{fs, pattern})))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	http.ServeFile(w, r, path.Join(staticDir, "index.html"))
}

// NewHandler creates the root http.Handler.
func NewHandler(l Library, hs history.Store, mediaFileSystem, artworkFileSystem store.FileSystem) http.Handler {
	var c httpauth.Checker = httpauth.None{}
	if authUser != "" {
		c = httpauth.Creds(map[string]string{
			authUser: authPassword,
		})
	}
	h := fsServeMux{
		httpauth.NewServeMux(c, http.NewServeMux()),
	}

	h.HandleFunc("/", rootHandler)
	h.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

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
