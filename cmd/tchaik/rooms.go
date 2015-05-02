package main

import (
	"sync"

	"golang.org/x/net/websocket"
)

type rooms struct {
	sync.RWMutex
	m map[string][]*websocket.Conn
}

func newRooms() *rooms {
	return &rooms{m: make(map[string][]*websocket.Conn)}
}

func (r *rooms) join(name string, ws *websocket.Conn) {
	r.Lock()
	defer r.Unlock()

	r.m[name] = append(r.m[name], ws)
}

func (r *rooms) remove(ws *websocket.Conn) {
	r.Lock()
	defer r.Unlock()

	for k, v := range r.m {
		for i, e := range v {
			if e == ws {
				r.m[k] = append(v[:i], v[i+1:]...)
			}
		}
	}
}

func (r *rooms) get(name string) []*websocket.Conn {
	r.RLock()
	defer r.RUnlock()

	return r.m[name]
}

func (r *rooms) getRoomNames() []string {
	r.RLock()
	defer r.RUnlock()

	keys := make([]string, 0, len(r.m))
	for k, _ := range r.m {
		keys = append(keys, k)
	}
	return keys
}
