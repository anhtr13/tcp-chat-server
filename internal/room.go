package server

import (
	"sync"
)

type room struct {
	room_name string
	members   map[string]*client
	mtx       sync.Mutex
}

func newRoom(room_name string) *room {
	return &room{
		room_name: room_name,
		members:   map[string]*client{},
		mtx:       sync.Mutex{},
	}
}

// Must lock mutex from outside
func (r *room) broadcast(msg string) {
	for _, c := range r.members {
		c.write(MESSAGE, msg)
	}
}
