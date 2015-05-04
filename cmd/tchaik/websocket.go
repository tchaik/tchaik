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

	"github.com/dhowden/tchaik/index"
)

type Command struct {
	Action string
	Data   interface{}
}

const (
	KeyAction         string = "KEY"
	CtrlAction        string = "CTRL"
	FetchAction       string = "FETCH"
	SearchAction      string = "SEARCH"
	PlayerAction      string = "PLAYER"
	FilterListAction  string = "FILTER_LIST"
	FilterPathsAction string = "FILTER_PATHS"
	FetchRecentAction string = "FETCH_RECENT"
)

func (l LibraryAPI) WebsocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var key string
		defer l.players.remove(key)

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
				resp, err = handleSearch(l, c)
			case FilterListAction:
				resp, err = handleFilterList(l, c)
			case FilterPathsAction:
				resp, err = handleFilterPaths(l, c)
			case KeyAction:
				key, err = handleKey(l, c, ws, key)
			case PlayerAction:
				err = handlePlayer(l, c)
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

func handlePlayer(l LibraryAPI, c Command) error {
	return nil
}

func handleKey(l LibraryAPI, c Command, ws *websocket.Conn, key string) (string, error) {
	key, ok := c.Data.(string)
	if !ok {
		return "", fmt.Errorf("data property should be a 'string', got '%T'", c.Data)
	}

	l.players.remove(key)
	if key != "" {
		l.players.add(key, ValidatedPlayer(websocketPlayer{ws}))
	}
	return key, nil
}

func extractPath(data map[string]interface{}) ([]string, error) {
	rawPath, ok := data["path"]
	if !ok {
		return nil, fmt.Errorf("expected 'path' in data map")
	}

	rawPathSlice, ok := rawPath.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected path to be a list of strings, got '%T'", rawPath)
	}

	path := make([]string, len(rawPathSlice))
	for i, x := range rawPathSlice {
		s, ok := x.(string)
		if !ok {
			return nil, fmt.Errorf("expected path component to be a 'string', got '%T'", x)
		}
		path[i] = s
	}
	return path, nil
}

func handleCollectionList(l LibraryAPI, c Command) (interface{}, error) {
	data, ok := c.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data property should be a 'map[string]interface{}', got '%T'", c.Data)
	}

	path, err := extractPath(data)
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
	filterName, ok := c.Data.(string)
	if !ok {
		return nil, fmt.Errorf("expected data to be a 'string', got '%T'", c.Data)
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
	data, ok := c.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data property should be a 'map[string]interface{}', got '%T'", c.Data)
	}

	path, err := extractPath(data)
	if err != nil {
		return nil, err
	}

	rawName, ok := data["name"]
	if !ok {
		return nil, fmt.Errorf("data map should contain a filter 'name' field")
	}

	filterName, ok := rawName.(string)
	if !ok {
		return nil, fmt.Errorf("expected filer 'name' to be a 'string', got '%T'", filterName)
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

func handleSearch(l LibraryAPI, c Command) (interface{}, error) {
	input, ok := c.Data.(string)
	if !ok {
		return nil, fmt.Errorf("expected 'data' to be a 'string', got '%T'", c.Data)
	}

	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   l.searcher.Search(input),
	}, nil
}
