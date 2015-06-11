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

	"github.com/tchaik/tchaik/index"
)

type Command struct {
	Action string
	Data   map[string]interface{}
}

func (c Command) getString(f string) (string, error) {
	raw, ok := c.Data[f]
	if !ok {
		return "", fmt.Errorf("expected '%s' in data map", f)
	}

	value, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("expected '%s' to be of type 'string', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getFloat(f string) (float64, error) {
	raw, ok := c.Data[f]
	if !ok {
		return 0.0, fmt.Errorf("expected '%s' in data map", f)
	}

	value, ok := raw.(float64)
	if !ok {
		return 0.0, fmt.Errorf("expected '%s' to be of type 'float64', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getBool(f string) (bool, error) {
	raw, ok := c.Data[f]
	if !ok {
		return false, fmt.Errorf("expected '%s' in data map", f)
	}

	value, ok := raw.(bool)
	if !ok {
		return false, fmt.Errorf("expected '%s' to be of type 'bool', got '%T'", f, raw)
	}
	return value, nil
}

func (c Command) getStringSlice(f string) ([]string, error) {
	raw, ok := c.Data[f]
	if !ok {
		return nil, fmt.Errorf("expected '%s' in data map", f)
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
	KeyAction         string = "KEY"
	CtrlAction               = "CTRL"
	FetchAction              = "FETCH"
	SearchAction             = "SEARCH"
	PlayerAction             = "PLAYER"
	FilterListAction         = "FILTER_LIST"
	FilterPathsAction        = "FILTER_PATHS"
	FetchRecentAction        = "FETCH_RECENT"
)

func (l LibraryAPI) WebsocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var key string
		defer l.players.remove(key)

		searcher := &sameSearcher{
			Searcher: l.searcher,
		}

		var err error
		for {
			var c Command
			err = websocket.JSON.Receive(ws, &c)
			if err != nil {
				if err != io.EOF {
					err = fmt.Errorf("receive: %v", err)
				}
				break
			}

			var resp interface{}
			switch c.Action {
			case FetchAction:
				resp, err = handleCollectionList(l, c)
			case SearchAction:
				resp, err = handleSearch(l, searcher, c)
			case FilterListAction:
				resp, err = handleFilterList(l, c)
			case FilterPathsAction:
				resp, err = handleFilterPaths(l, c)
			case KeyAction:
				key, err = handleKey(l, c, ws, key)
			case PlayerAction:
				resp, err = handlePlayer(l, c)
			case FetchRecentAction:
				resp = handleFetchRecent(l, c)
			default:
				err = fmt.Errorf("unknown action: %v", c.Action)
			}

			if err != nil {
				break
			}

			if resp == nil {
				continue
			}

			err = websocket.JSON.Send(ws, resp)
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
	})
}

func handlePlayer(l LibraryAPI, c Command) (interface{}, error) {
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
			Data:   l.players.list(),
		}, nil
	}

	key, err := c.getString("key")
	if err != nil {
		return nil, err
	}

	p := l.players.get(key)
	if p == nil {
		return nil, fmt.Errorf("invalid player key: %v", key)
	}

	switch action {
	case "PLAY":
		err = p.Play()

	case "PAUSE":
		err = p.Pause()

	case "NEXT":
		err = p.NextTrack()

	case "PREV":
		err = p.PreviousTrack()

	case "TOGGLE_PLAY_PAUSE":
		err = p.TogglePlayPause()

	case "TOGGLE_MUTE":
		err = p.ToggleMute()

	case "SET_VOLUME":
		var f float64
		f, err = c.getFloat("value")
		if err == nil {
			err = p.SetVolume(f)
		}

	case "SET_MUTE":
		var b bool
		b, err = c.getBool("value")
		if err == nil {
			err = p.SetMute(b)
		}

	case "SET_TIME":
		var f float64
		f, err = c.getFloat("value")
		if err == nil {
			err = p.SetTime(f)
		}
	}

	return nil, err
}

func handleKey(l LibraryAPI, c Command, ws *websocket.Conn, key string) (string, error) {
	key, err := c.getString("key")
	if err != nil {
		return "", err
	}

	l.players.remove(key)
	if key != "" {
		l.players.add(ValidatedPlayer(WebsocketPlayer(key, ws)))
	}
	return key, nil
}

func handleCollectionList(l LibraryAPI, c Command) (interface{}, error) {
	path, err := c.getStringSlice("path")
	if err != nil {
		return nil, err
	}

	if len(path) < 1 {
		return nil, fmt.Errorf("invalid path: %v\n", path)
	}

	root := l.collections[path[0]]
	if root == nil {
		return nil, fmt.Errorf("unknown collection: %#v", path[0])
	}
	g, err := l.Fetch(root, path[1:])
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

func handleFilterList(l LibraryAPI, c Command) (interface{}, error) {
	filterName, err := c.getString("name")
	if err != nil {
		return nil, err
	}

	filterItems, ok := l.filters[filterName]
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

func handleFilterPaths(l LibraryAPI, c Command) (interface{}, error) {
	path, err := c.getStringSlice("path")
	if err != nil {
		return nil, err
	}

	filterName, err := c.getString("name")
	if err != nil {
		return nil, err
	}

	filterItems, ok := l.filters[filterName]
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
			Paths []index.Path
		}{
			Path:  []string{filterName, name},
			Paths: item.Paths(),
		},
	}, nil
}

func handleFetchRecent(l LibraryAPI, c Command) interface{} {
	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   l.recent,
	}
}

func handleSearch(l LibraryAPI, r *sameSearcher, c Command) (interface{}, error) {
	input, err := c.getString("input")
	if err != nil {
		return nil, err
	}

	paths := r.Search(input)
	if r.same {
		return nil, nil
	}
	data := build(index.NewPathsCollection(l.collections["Root"], paths), index.Key("Root"))

	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   data,
	}, nil
}
