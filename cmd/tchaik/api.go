package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/dhowden/tchaik/index"
)

type LibraryAPI struct {
	index.Library

	root     index.Collection
	searcher index.Searcher
	sessions *sessions
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

func build(g index.Group, key index.Key) group {
	h := group{
		Name:        g.Name(),
		Key:         key,
		TotalTime:   g.Field("TotalTime"),
		Artist:      g.Field("Artist"),
		AlbumArtist: g.Field("AlbumArtist"),
		Composer:    g.Field("Composer"),
		Year:        g.Field("Year"),
		ListStyle:   g.Field("ListStyle"),
		TrackID:     g.Field("TrackID"),
	}

	if c, ok := g.(index.Collection); ok {
		for _, k := range c.Keys() {
			sg := c.Get(k)
			sg = index.FirstTrackAttr(index.StringAttr("TrackID"), sg)
			h.Groups = append(h.Groups, group{
				Name:    sg.Name(),
				Key:     k,
				TrackID: sg.Field("TrackID"),
			})
		}
		return h
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
	Data   string
	Path   []string
}

const (
	KeyAction    string = "KEY"
	CtrlAction   string = "CTRL"
	FetchAction  string = "FETCH"
	SearchAction string = "SEARCH"
)

var websocketSessions map[string]*websocket.Conn

func (l LibraryAPI) WebsocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		defer l.sessions.remove(ws)

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
				resp = handleSearch(l, c)
			case KeyAction:
				handleKey(l, c, ws)
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

type sessions struct {
	sync.RWMutex
	m map[string]*websocket.Conn
}

func newSessions() *sessions {
	return &sessions{m: make(map[string]*websocket.Conn)}
}

func (s *sessions) add(id string, ws *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	s.m[id] = ws
}

func (s *sessions) remove(ws *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.m {
		if v == ws {
			delete(s.m, k)
			return
		}
	}
}

func (s *sessions) get(id string) *websocket.Conn {
	s.RLock()
	defer s.RUnlock()

	return s.m[id]
}

func handleKey(l LibraryAPI, c Command, ws *websocket.Conn) {
	l.sessions.remove(ws)
	if c.Data == "" {
		return
	}
	l.sessions.add(c.Data, ws)
	return
}

func handleCollectionList(l LibraryAPI, c Command) (interface{}, error) {
	if len(c.Path) < 1 {
		return nil, fmt.Errorf("invalid path: %v\n", c.Path)
	}

	g, err := l.Fetch(l.root, c.Path[1:])
	if err != nil {
		return nil, fmt.Errorf("error in Fetch: %v (path: %#v)", err, c.Path[1:])
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
			c.Path,
			g,
		},
	}, nil
}

func handleSearch(l LibraryAPI, c Command) interface{} {
	return struct {
		Action string
		Data   interface{}
	}{
		Action: c.Action,
		Data:   l.searcher.Search(c.Data),
	}
}

func ctrlHandler(l LibraryAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "invalid request method", http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "error parsing parameters", http.StatusInternalServerError)
			return
		}
		k := r.Form.Get("key")
		ws := l.sessions.get(k)
		if ws == nil {
			http.Error(w, "invalid session key", http.StatusBadRequest)
			return
		}

		var data interface{}
		switch r.URL.Path {
		case "play", "pause", "next", "prev":
			data = strings.ToUpper(r.URL.Path)

		case "volume":
			v := r.Form.Get("value")
			f, err := strconv.ParseFloat(v, 32)
			if err != nil || f > 1.0 || f < 0.0 {
				http.Error(w, "invalid volume value (expected float between 0.0 and 1.0)", http.StatusBadRequest)
				return
			}
			data = struct {
				Key   string
				Value float64
			}{
				Key:   "volume",
				Value: f,
			}

		case "mute":
			v := r.Form.Get("value")
			b, err := strconv.ParseBool(v)
			if err != nil {
				http.Error(w, "invalid bool value", http.StatusBadRequest)
				return
			}
			data = struct {
				Key   string
				Value bool
			}{
				Key:   "mute",
				Value: b,
			}

		case "time":
			v := r.Form.Get("value")
			f, err := strconv.ParseFloat(v, 32)
			if err != nil || f < 0.0 {
				http.Error(w, "invalid time value (expected float greater than 0.0)", http.StatusBadRequest)
				return
			}
			data = struct {
				Key   string
				Value float64
			}{
				Key:   "time",
				Value: f,
			}

		default:
			http.NotFound(w, r)
			return
		}

		err = websocket.JSON.Send(ws, struct {
			Action string
			Data   interface{}
		}{
			Action: CtrlAction,
			Data:   data,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("error sending command: %v", err), http.StatusInternalServerError)
			return
		}
		return
	})
}
