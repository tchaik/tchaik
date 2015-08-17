// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"tchaik.com/index"
	"tchaik.com/index/favourite"
	"tchaik.com/index/history"
	"tchaik.com/player"
)

// Command is a type which is a container for data received from the websocket.
type Command struct {
	Action string
	Data   map[string]interface{}
}

func (c Command) get(f string) (interface{}, error) {
	raw, ok := c.Data[f]
	if !ok {
		return nil, fmt.Errorf("expected '%s' in data map", f)
	}
	return raw, nil
}

func (c Command) getString(f string) (string, error) {
	raw, err := c.get(f)
	if err != nil {
		return "", err
	}

	value, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("expected '%s' to be of type 'string', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getFloat(f string) (float64, error) {
	raw, err := c.get(f)
	if err != nil {
		return 0.0, err
	}

	value, ok := raw.(float64)
	if !ok {
		return 0.0, fmt.Errorf("expected '%s' to be of type 'float64', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getBool(f string) (bool, error) {
	raw, err := c.get(f)
	if err != nil {
		return false, err
	}

	value, ok := raw.(bool)
	if !ok {
		return false, fmt.Errorf("expected '%s' to be of type 'bool', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getStringSlice(f string) ([]string, error) {
	raw, err := c.get(f)
	if err != nil {
		return nil, err
	}

	rawSlice, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected '%s' to be a list of strings, got '%T'", f, raw)
	}

	result := make([]string, len(rawSlice))
	for i, x := range rawSlice {
		s, ok := x.(string)
		if !ok {
			return nil, fmt.Errorf("expected '%s' to contain objects of type 'string', got '%T'", f, x)
		}
		result[i] = s
	}
	return result, nil
}

type sameSearcher struct {
	index.Searcher
	paths []index.Path
	same  bool
}

func (r *sameSearcher) Search(input string) []index.Path {
	paths := r.Searcher.Search(input)
	r.same = false
	if len(r.paths) == len(paths) {
		r.same = true
		for i, path := range r.paths {
			if path[1] != paths[i][1] {
				r.same = false
				break
			}
		}
	}
	r.paths = paths
	return paths
}

const (
	// Player Actions
	KeyAction    string = "KEY"
	PlayerAction        = "PLAYER"

	// Path Actions
	RecordPlayAction   = "RECORD_PLAY"
	SetFavouriteAction = "SET_FAVOURITE"

	// Library Actions
	CtrlAction           = "CTRL"
	FetchAction          = "FETCH"
	SearchAction         = "SEARCH"
	FilterListAction     = "FILTER_LIST"
	FilterPathsAction    = "FILTER_PATHS"
	FetchRecentAction    = "FETCH_RECENT"
	FetchFavouriteAction = "FETCH_FAVOURITE"
)

// NewWebsocketHandler creates a websocket handler for the library, players and history.
func NewWebsocketHandler(l Library, p *player.Players, s history.Store) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		h := &websocketHandler{
			Conn:    ws,
			lib:     l,
			players: p,
			history: s,
			searcher: &sameSearcher{
				Searcher: l.searcher,
			},
		}
		h.Handle()
	})
}

type websocketHandler struct {
	*websocket.Conn
	players  *player.Players
	history  history.Store
	lib      Library
	searcher *sameSearcher

	playerKey string
}

func (h *websocketHandler) Handle() {
	defer h.players.Remove(h.playerKey)

	var err error
	for {
		var c Command
		err = websocket.JSON.Receive(h.Conn, &c)
		if err != nil {
			if err != io.EOF {
				err = fmt.Errorf("receive: %v", err)
			}
			break
		}

		var resp interface{}
		switch c.Action {

		// Player actions
		case KeyAction:
			err = h.key(c)
		case PlayerAction:
			resp, err = h.player(c)

		// Path Actions
		case RecordPlayAction:
			err = h.recordPlay(c)

		case SetFavouriteAction:
			err = h.setFavourite(c)

		// Library actions
		case FetchAction:
			resp, err = h.collectionList(c)
		case SearchAction:
			resp, err = h.search(c)
		case FilterListAction:
			resp, err = h.filterList(c)
		case FilterPathsAction:
			resp, err = h.filterPaths(c)
		case FetchRecentAction:
			resp = h.fetchRecent(c)
		case FetchFavouriteAction:
			resp = h.fetchFavourite(c)
		default:
			err = fmt.Errorf("unknown action: %v", c.Action)
		}

		if err != nil {
			break
		}

		if resp == nil {
			continue
		}

		err = websocket.JSON.Send(h.Conn, resp)
		if err != nil {
			if err != io.EOF {
				err = fmt.Errorf("send: %v", err)
			}
			break
		}
	}

	if err != nil && err != io.EOF {
		log.Printf("socket error: %v", err)
	}
}

func (h *websocketHandler) player(c Command) (interface{}, error) {
	action, err := c.getString("action")
	if err != nil {
		return nil, err
	}

	if action == "LIST" {
		return struct {
			Action string
			Data   interface{}
		}{
			Action: c.Action,
			Data:   h.players.List(),
		}, nil
	}

	key, err := c.getString("key")
	if err != nil {
		return nil, err
	}

	p := h.players.Get(key)
	if p == nil {
		return nil, fmt.Errorf("invalid player key: %v", key)
	}

	r := player.RepAction{
		Action: action,
		Value:  c.Data["value"],
	}
	return nil, r.Apply(p)
}

func (h *websocketHandler) key(c Command) error {
	key, err := c.getString("key")
	if err != nil {
		return err
	}

	h.players.Remove(h.playerKey)
	if key != "" {
		h.players.Add(player.Validated(WebsocketPlayer(key, h.Conn)))
	}
	h.playerKey = key
	return nil
}

func (h *websocketHandler) recordPlay(c Command) error {
	path, err := c.getStringSlice("path")
	if err != nil {
		return err
	}
	p := make([]index.Key, len(path))
	for i, x := range path {
		p[i] = index.Key(x)
	}
	return h.history.Add(p)
}

func (h *websocketHandler) setFavourite(c Command) error {
	path, err := c.getStringSlice("path")
	if err != nil {
		return err
	}
	value, err := c.getBool("value")
	if err != nil {
		return err
	}
	p := make([]index.Key, len(path))
	for i, x := range path {
		p[i] = index.Key(x)
	}
	return h.lib.favourites.Set(p, value)
}

func (h *websocketHandler) collectionList(c Command) (interface{}, error) {
	path, err := c.getStringSlice("path")
	if err != nil {
		return nil, err
	}

	if len(path) < 1 {
		return nil, fmt.Errorf("invalid path: %v\n", path)
	}

	root := h.lib.collections[path[0]]
	if root == nil {
		return nil, fmt.Errorf("unknown collection: %#v", path[0])
	}
	g, err := h.lib.Fetch(root, path[1:])
	if err != nil {
		return nil, fmt.Errorf("error in Fetch: %v (path: %#v)", err, path[1:])
	}

	return struct {
		Action string
		Data   interface{}
	}{
		c.Action,
		struct {
			Path []string
			Item group
		}{
			path,
			g,
		},
	}, nil
}

func (h *websocketHandler) filterList(c Command) (interface{}, error) {
	filterName, err := c.getString("name")
	if err != nil {
		return nil, err
	}

	filterItems, ok := h.lib.filters[filterName]
	if !ok {
		return nil, fmt.Errorf("invalid filter name: %#v", filterName)
	}

	filterNames := make([]string, len(filterItems))
	for i, x := range filterItems {
		filterNames[i] = x.Name()
	}
	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data: struct {
			Name  string
			Items []string
		}{
			Name:  filterName,
			Items: filterNames,
		},
	}, nil
}

func (h *websocketHandler) filterPaths(c Command) (interface{}, error) {
	path, err := c.getStringSlice("path")
	if err != nil {
		return nil, err
	}

	filterName, err := c.getString("name")
	if err != nil {
		return nil, err
	}

	filterItems, ok := h.lib.filters[filterName]
	if !ok {
		return nil, fmt.Errorf("invalid filter name: %#v", filterName)
	}

	if len(path) != 1 {
		return nil, fmt.Errorf("invalid path: %#v", path)
	}
	name := path[0]

	var item index.FilterItem
	for _, x := range filterItems {
		if x.Name() == name {
			item = x
			break
		}
	}
	if item == nil {
		return nil, fmt.Errorf("invalid filter item: %#v", name)
	}

	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data: struct {
			Path  []string
			Paths group
		}{
			Path:  []string{filterName, name},
			Paths: h.lib.ExpandPaths(item.Paths()),
		},
	}, nil
}

func (h *websocketHandler) fetchRecent(c Command) interface{} {
	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   h.lib.ExpandPaths(h.lib.recent),
	}
}

func (h *websocketHandler) fetchFavourite(c Command) interface{} {
	paths := index.CollectionPaths(h.lib.collections["Root"], []index.Key{"Root"})
	filter := favourite.RootFilter{Store: h.lib.favourites}

	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   h.lib.ExpandPaths(filter.Filter(paths)),
	}
}

func (h *websocketHandler) search(c Command) (interface{}, error) {
	input, err := c.getString("input")
	if err != nil {
		return nil, err
	}

	paths := h.searcher.Search(input)
	if h.searcher.same {
		return nil, nil
	}

	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   h.lib.ExpandPaths(paths),
	}, nil
}

// WebsocketPlayer creates a player.Player which sends commands down the websocket.Conn when
// player.Player methods are called.
func WebsocketPlayer(key string, ws *websocket.Conn) player.Player {
	repFn := func(data interface{}) {
		websocket.JSON.Send(ws, struct {
			Action string
			Data   interface{}
		}{
			Action: CtrlAction,
			Data:   data,
		})
	}
	return player.NewRep(key, repFn)
}
