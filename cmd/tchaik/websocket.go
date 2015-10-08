// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"

	"tchaik.com/index"
	"tchaik.com/index/attr"
	"tchaik.com/index/cursor"
	"tchaik.com/index/playlist"
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

func (c Command) getInt(f string) (int, error) {
	raw, err := c.getFloat(f)
	if err != nil {
		return 0, err
	}
	return int(raw), nil
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

func (c Command) getPath(f string) (index.Path, error) {
	raw, err := c.get(f)
	if err != nil {
		return nil, err
	}

	return index.PathFromJSONInterface(raw)
}

func newBootstrapSearcher(root index.Collection) index.Searcher {
	return &bootstrapSearcher{
		root: root,
	}
}

type bootstrapSearcher struct {
	once sync.Once
	root index.Collection

	index.Searcher
}

func (b *bootstrapSearcher) bootstrap() {
	wi := index.BuildCollectionWordIndex(b.root, []string{"Composer", "Artist", "Album", "Name"})
	b.Searcher = index.FlatSearcher{
		Searcher: index.WordsIntersectSearcher(index.BuildPrefixExpandSearcher(wi, wi, 10)),
	}
}

func (b *bootstrapSearcher) Search(input string) []index.Path {
	b.once.Do(b.bootstrap)
	return b.Searcher.Search(input)
}

func newBootstrapFilter(root index.Collection, field attr.Interface) index.Filter {
	return &bootstrapFilter{
		root:  root,
		field: field,
	}
}

type bootstrapFilter struct {
	once  sync.Once
	root  index.Collection
	field attr.Interface

	index.Filter
}

func (b *bootstrapFilter) bootstrap() {
	b.Filter = index.FilterCollection(b.root, b.field)
}

func (b *bootstrapFilter) Items() []index.FilterItem {
	b.once.Do(b.bootstrap)
	return b.Filter.Items()
}

type bootstrapRecent struct {
	once sync.Once
	root index.Collection
	n    int

	list []index.Path
}

func (b *bootstrapRecent) bootstrap() {
	b.list = index.Recent(b.root, b.n)
}

func (b *bootstrapRecent) List() []index.Path {
	b.once.Do(b.bootstrap)
	return b.list
}

// sameSearcher is a light wrapper around a index.Searher which caches the path
// slice returned by Search and sets the attribute `same` to true when subsequent
// searches return the same result (and hence does not need to be re-transmitted).
type sameSearcher struct {
	index.Searcher
	paths []index.Path
	same  bool
}

// Search implements index.Searcher.
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
	ActionKey    string = "KEY"
	ActionPlayer        = "PLAYER"

	// Path Actions
	ActionRecordPlay   = "RECORD_PLAY"
	ActionSetFavourite = "SET_FAVOURITE"
	ActionSetChecklist = "SET_CHECKLIST"

	// Playlist Actions
	ActionPlaylist = "PLAYLIST"

	// Cursor Actions
	ActionCursor = "CURSOR"

	// Library Actions
	ActionCtrl          = "CTRL"
	ActionFetch         = "FETCH"
	ActionSearch        = "SEARCH"
	ActionFilterList    = "FILTER_LIST"
	ActionFilterPaths   = "FILTER_PATHS"
	ActionFetchPathList = "FETCH_PATHLIST"
)

type websocketHandlerFunc func(c Command, r *Response) error

type websocketMux struct {
	m map[string]websocketHandlerFunc
}

func (w *websocketMux) HandleFunc(a string, fn websocketHandlerFunc) {
	w.m[a] = fn
}

func (w *websocketMux) Handle(c Command, r *Response) error {
	fn, ok := w.m[c.Action]
	if !ok {
		return fmt.Errorf("unknown action: %v", c.Action)
	}
	return fn(c, r)
}

// NewWebsocketHandler creates a websocket handler for the library, players and history.
func NewWebsocketHandler(l Library, m *Meta, p *player.Players) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		mux := &websocketMux{
			m: make(map[string]websocketHandlerFunc),
		}

		h := &websocketHandler{
			Conn:    ws,
			mux:     mux,
			lib:     l,
			meta:    m,
			players: p,
			searcher: &sameSearcher{
				Searcher: l.searcher,
			},
		}

		mux.HandleFunc(ActionKey, h.key)
		mux.HandleFunc(ActionPlayer, h.player)
		mux.HandleFunc(ActionRecordPlay, h.recordPlay)
		mux.HandleFunc(ActionSetFavourite, h.setFavourite)
		mux.HandleFunc(ActionSetChecklist, h.setChecklist)
		mux.HandleFunc(ActionPlaylist, h.playlist)
		mux.HandleFunc(ActionCursor, h.cursor)
		mux.HandleFunc(ActionFetch, h.collectionList)
		mux.HandleFunc(ActionSearch, h.search)
		mux.HandleFunc(ActionFilterList, h.filterList)
		mux.HandleFunc(ActionFilterPaths, h.filterPaths)
		mux.HandleFunc(ActionFetchPathList, h.fetchPathList)

		h.handle()
	})
}

type websocketHandler struct {
	*websocket.Conn
	mux      *websocketMux
	players  *player.Players
	lib      Library
	searcher *sameSearcher
	meta     *Meta

	playerKey string
}

func (h *websocketHandler) handle() {
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

		resp := &Response{
			Action: c.Action,
		}
		err = h.mux.Handle(c, resp)
		if err != nil {
			break
		}
		if resp.Data == nil {
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

// Response is a type which represnets a response to a Websocket Command.
type Response struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (h *websocketHandler) player(c Command, resp *Response) error {
	action, err := c.getString("action")
	if err != nil {
		return err
	}

	if action == "LIST" {
		resp.Data = h.players.List()
		return nil
	}

	key, err := c.getString("key")
	if err != nil {
		return err
	}

	p := h.players.Get(key)
	if p == nil {
		return fmt.Errorf("invalid player key: %v", key)
	}

	r := player.RepAction{
		Action: action,
		Value:  c.Data["value"],
	}
	return r.Apply(p)
}

func (h *websocketHandler) key(c Command, resp *Response) error {
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

func (h *websocketHandler) recordPlay(c Command, resp *Response) error {
	p, err := c.getPath("path")
	if err != nil {
		return err
	}
	return h.meta.history.Add(p)
}

func (h *websocketHandler) setFavourite(c Command, resp *Response) error {
	p, err := c.getPath("path")
	if err != nil {
		return err
	}
	value, err := c.getBool("value")
	if err != nil {
		return err
	}
	return h.meta.favourites.Set(p, value)
}

func (h *websocketHandler) setChecklist(c Command, resp *Response) error {
	p, err := c.getPath("path")
	if err != nil {
		return err
	}
	value, err := c.getBool("value")
	if err != nil {
		return err
	}
	return h.meta.checklist.Set(p, value)
}

func (h *websocketHandler) cursor(c Command, resp *Response) error {
	name, err := c.getString("name")
	if err != nil {
		return err
	}

	action, err := c.getString("action")
	if err != nil {
		return err
	}

	if action != "FETCH" {
		path, _ := c.getPath("path")
		index, _ := c.getInt("index")

		ra := cursor.RepAction{
			Name:   name,
			Action: cursor.Action(action),
			Path:   path,
			Index:  index,
		}

		root := &rootCollection{h.lib.collections["Root"]}
		err = ra.Apply(h.meta.cursors, h.meta.playlists, root)
		if err != nil {
			return err
		}
	}

	resp.Data = h.meta.cursors.Get(name)
	return nil
}

func (h *websocketHandler) playlist(c Command, resp *Response) error {
	name, err := c.getString("name")
	if err != nil {
		return err
	}

	action, err := c.getString("action")
	if err != nil {
		return err
	}

	if action != "FETCH" {
		path, err := c.getPath("path")
		if err != nil {
			return err
		}
		index, _ := c.getInt("index")

		ra := playlist.RepAction{
			Name:   name,
			Action: playlist.Action(action),
			Path:   path,
			Index:  index,
		}

		err = ra.Apply(h.meta.playlists)
		if err != nil {
			return err
		}
	}

	resp.Data = h.meta.playlists.Get(name)
	return nil
}

func (h *websocketHandler) collectionList(c Command, resp *Response) error {
	p, err := c.getPath("path")
	if err != nil {
		return err
	}

	g, k, err := h.lib.Fetch(p)
	if err != nil {
		return err
	}
	g = h.meta.Annotate(p, g)

	resp.Data = struct {
		Path index.Path  `json:"path"`
		Item index.Group `json:"item"`
	}{
		p,
		&Group{
			Group: g,
			Key:   k,
		},
	}
	return nil
}

func (h *websocketHandler) filterList(c Command, resp *Response) error {
	filterName, err := c.getString("name")
	if err != nil {
		return err
	}

	filter, ok := h.lib.filters[filterName]
	if !ok {
		return fmt.Errorf("invalid filter name: %#v", filterName)
	}

	filterNames := make([]string, len(filter.Items()))
	for i, x := range filter.Items() {
		filterNames[i] = x.Name()
	}

	resp.Data = struct {
		Name  string   `json:"name"`
		Items []string `json:"items"`
	}{
		Name:  filterName,
		Items: filterNames,
	}
	return nil
}

func (h *websocketHandler) filterPaths(c Command, resp *Response) error {
	path, err := c.getPath("path")
	if err != nil {
		return err
	}

	filterName, err := c.getString("name")
	if err != nil {
		return err
	}

	filter, ok := h.lib.filters[filterName]
	if !ok {
		return fmt.Errorf("invalid filter name: %#v", filterName)
	}

	if len(path) != 1 {
		return fmt.Errorf("invalid path: %#v", path)
	}
	name := string(path[0])

	var item index.FilterItem
	for _, x := range filter.Items() {
		if x.Name() == name {
			item = x
			break
		}
	}
	if item == nil {
		return fmt.Errorf("invalid filter item: %#v", name)
	}

	resp.Data = struct {
		Path  index.Path  `json:"path"`
		Paths index.Group `json:"paths"`
	}{
		Path:  index.PathFromStringSlice([]string{filterName, name}),
		Paths: h.lib.ExpandPaths(item.Paths()),
	}
	return nil
}

// Lister is an interface which defines the List method.
type Lister interface {
	// List returns a list of index.Paths.
	List() []index.Path
}

// NB: matches paths based only on the first two keys.
func filterByRootLister(l Lister, paths []index.Path) []index.Path {
	exp := make(map[string]bool)
	for _, p := range l.List() {
		if len(p) > 1 {
			exp[fmt.Sprintf("%v", p[:2])] = true
		}
	}

	result := make([]index.Path, 0, len(paths))
	for _, p := range paths {
		if exp[fmt.Sprintf("%v", p)] {
			result = append(result, p)
		}
	}
	return result
}

func (h *websocketHandler) fetchPathList(c Command, resp *Response) error {
	name, err := c.getString("name")
	if err != nil {
		return err
	}

	var paths []index.Path
	switch name {
	case "recent":
		paths = h.lib.recent.List()

	case "favourite":
		paths = index.CollectionPaths(h.lib.collections["Root"], []index.Key{"Root"})
		paths = filterByRootLister(h.meta.favourites, paths)

	case "checklist":
		paths = index.CollectionPaths(h.lib.collections["Root"], []index.Key{"Root"})
		paths = filterByRootLister(h.meta.checklist, paths)
	}

	resp.Data = struct {
		Name string      `json:"name"`
		Data index.Group `json:"data"`
	}{
		Name: name,
		Data: h.lib.ExpandPaths(paths),
	}
	return nil
}

func (h *websocketHandler) search(c Command, resp *Response) error {
	input, err := c.getString("input")
	if err != nil {
		return err
	}

	paths := h.searcher.Search(input)
	if h.searcher.same {
		return nil
	}

	resp.Data = h.lib.ExpandPaths(paths)
	return nil
}

// WebsocketPlayer creates a player.Player which sends commands down the websocket.Conn when
// player.Player methods are called.
func WebsocketPlayer(key string, ws *websocket.Conn) player.Player {
	repFn := func(data interface{}) {
		websocket.JSON.Send(ws, &Response{
			Action: ActionCtrl,
			Data:   data,
		})
	}
	return player.NewRep(key, repFn)
}
