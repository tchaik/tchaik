package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/dhowden/tchaik/index"
)

type LibraryAPI struct {
	index.Library

	collections map[string]index.Collection
	filters     map[string][]index.FilterItem
	recent      []index.Path
	searcher    index.Searcher
	players     *players
}

type libraryFileSystem struct {
	http.FileSystem
	index.Library
}

// Open implements http.FileSystem and rewrites TrackID values to their corresponding Location
// values using the index.Library
func (l *libraryFileSystem) Open(path string) (http.File, error) {
	t, ok := l.Library.Track(strings.Trim(path, "/")) // TrackIDs arrive with leading slash
	if !ok {
		return nil, fmt.Errorf("could not find track: %v", path)
	}

	loc := t.GetString("Location")
	if loc == "" {
		return nil, fmt.Errorf("invalid (empty) location for track: %v", path)
	}
	return l.FileSystem.Open(loc)
}

type group struct {
	Name        string
	Key         index.Key
	TotalTime   interface{} `json:",omitempty"`
	Artist      interface{} `json:",omitempty"`
	AlbumArtist interface{} `json:",omitempty"`
	Composer    interface{} `json:",omitempty"`
	ListStyle   interface{} `json:",omitempty"`
	TrackID     interface{} `json:",omitempty"`
	Year        interface{} `json:",omitempty"`
	Groups      []group     `json:",omitempty"`
	Tracks      []track     `json:",omitempty"`
}

type track struct {
	TrackID     string `json:",omitempty"`
	Name        string `json:",omitempty"`
	Album       string `json:",omitempty"`
	Artist      string `json:",omitempty"`
	AlbumArtist string `json:",omitempty"`
	Composer    string `json:",omitempty"`
	Year        int    `json:",omitempty"`
	DiscNumber  int    `json:",omitempty"`
	TotalTime   int    `json:",omitempty"`
}

func buildCollection(h group, c index.Collection) group {
	for _, k := range c.Keys() {
		g := c.Get(k)
		g = index.FirstTrackAttr(index.StringAttr("AlbumArtist"), g)
		h.Groups = append(h.Groups, group{
			Name:        g.Name(),
			Key:         k,
			AlbumArtist: g.Field("AlbumArtist"),
		})
	}
	return h
}

func build(g index.Group, key index.Key) group {
	h := group{
		Name:        g.Name(),
		Key:         key,
		TotalTime:   g.Field("TotalTime"),
		Artist:      g.Field("Artist"),
		AlbumArtist: g.Field("AlbumArtist"),
		// AlbumArtist: fmt.Sprintf("%v", g.Field("AlbumArtist")),
		Composer:  g.Field("Composer"),
		Year:      g.Field("Year"),
		ListStyle: g.Field("ListStyle"),
		TrackID:   g.Field("TrackID"),
	}

	if c, ok := g.(index.Collection); ok {
		return buildCollection(h, c)
	}

	getString := func(t index.Track, field string) string {
		if g.Field(field) != "" {
			return ""
		}
		return t.GetString(field)
	}

	getInt := func(t index.Track, field string) int {
		if g.Field(field) != 0 {
			return 0
		}
		return t.GetInt(field)
	}

	for _, t := range g.Tracks() {
		h.Tracks = append(h.Tracks, track{
			TrackID:    t.GetString("TrackID"),
			Name:       t.GetString("Name"),
			TotalTime:  t.GetInt("TotalTime"),
			DiscNumber: t.GetInt("DiscNumber"),
			// Potentially common fields (don't want to re-transmit everything)
			Album:       getString(t, "Album"),
			Artist:      getString(t, "Artist"),
			AlbumArtist: getString(t, "AlbumArtist"),
			Composer:    getString(t, "Composer"),
			Year:        getInt(t, "Year"),
		})
	}
	return h
}

func (l *LibraryAPI) Fetch(c index.Collection, path []string) (group, error) {
	if len(path) == 0 {
		return build(c, index.Key("Root")), nil
	}

	var g index.Group = c
	k := index.Key(path[0])
	g = c.Get(k)

	if g == nil {
		return group{}, fmt.Errorf("invalid path: near '%v'", path[0])
	}

	index.Sort(g.Tracks(), index.MultiSort(index.SortByInt("DiscNumber"), index.SortByInt("TrackNumber")))
	c = index.Collect(g, index.ByPrefix("Name"))
	g = index.SubTransform(c, index.TrimEnumPrefix)
	g = index.SumGroupIntAttr("TotalTime", g)
	commonFields := []index.Attr{
		index.StringAttr("Album"),
		index.StringAttr("Artist"),
		index.StringAttr("AlbumArtist"),
		index.StringAttr("Composer"),
		index.IntAttr("Year"),
	}
	g = index.CommonGroupAttr(commonFields, g)
	g = index.RemoveEmptyCollections(g)

	for i, p := range path[1:] {
		var ok bool
		c, ok = g.(index.Collection)
		if !ok {
			return group{}, fmt.Errorf("retrieved Group is not a Collection")
		}
		k = index.Key(p)
		g = c.Get(k)

		if g == nil {
			return group{}, fmt.Errorf("invalid path near '%v'", path[1:][i])
		}

		if _, ok = g.(index.Collection); !ok {
			if i == len(path[1:])-1 {
				break
			}
			return group{}, fmt.Errorf("retrieved Group isn't a Collection: %v", p)
		}
	}
	if g == nil {
		return group{}, fmt.Errorf("could not find group")
	}
	g = index.FirstTrackAttr(index.StringAttr("TrackID"), g)

	return build(g, k), nil
}

func (l *LibraryAPI) FileSystem(fs http.FileSystem) http.FileSystem {
	return &libraryFileSystem{fs, l.Library}
}

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

type players struct {
	sync.RWMutex
	m map[string]Player
}

func newPlayers() *players {
	return &players{m: make(map[string]Player)}
}

func (s *players) add(id string, p Player) {
	s.Lock()
	defer s.Unlock()

	s.m[id] = p
}

func (s *players) remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)
}

func (s *players) get(id string) Player {
	s.RLock()
	defer s.RUnlock()

	return s.m[id]
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

func createPlayer(l LibraryAPI, w http.ResponseWriter, r *http.Request) {
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

	if p := l.players.get(postData.Key); p != nil {
		http.Error(w, "player key already exists", http.StatusBadRequest)
		return
	}

	if postData.PlayerKeys == nil || len(postData.PlayerKeys) == 0 {
		http.Error(w, "no player keys specified", http.StatusBadRequest)
		return
	}

	var players []Player
	for _, pk := range postData.PlayerKeys {
		p := l.players.get(pk)
		if p == nil {
			http.Error(w, fmt.Sprintf("invalid player key: %v", pk), http.StatusBadRequest)
			return
		}
		players = append(players, p)
	}
	l.players.add(postData.Key, MultiPlayer(players...))
	w.WriteHeader(http.StatusCreated)
}

func playerAction(p Player, w http.ResponseWriter, r *http.Request) {
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
		return
	}

	if err != nil {
		if err, ok := err.(InvalidValueError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("error sending player command: %v", err), http.StatusInternalServerError)
	}
}

func playersHandler(l LibraryAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" && r.Method == "POST" {
			createPlayer(l, w, r)
			return
		}

		paths := strings.Split(r.URL.Path, "/")

		if len(paths) != 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		p := l.players.get(paths[0])
		if p == nil {
			http.Error(w, "invalid player key", http.StatusBadRequest)
			return
		}

		if r.Method == "DELETE" {
			l.players.remove(paths[0])
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == "PUT" {
			playerAction(p, w, r)
		}
	})
}
