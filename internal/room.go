package server

import (
	"sync"
)

type room struct {
	room_name string
	members   map[string]*client
	mtx       sync.RWMutex
}

func new_room(room_name string) *room {
	return &room{
		room_name: room_name,
		members:   map[string]*client{},
		mtx:       sync.RWMutex{},
	}
}

func (r *room) add_member(c *client) {
	r.mtx.Lock()
	r.members[c.client_id] = c
	r.mtx.Unlock()
}

func (r *room) remove_member(c *client) {
	r.mtx.Lock()
	delete(r.members, c.client_id)
	r.mtx.Unlock()
}

func (r *room) broadcast(msg string) {
	members := []*client{}

	r.mtx.RLock()
	for _, c := range r.members {
		members = append(members, c)
	}
	r.mtx.RUnlock()

	for _, c := range members {
		c.write(MESSAGE, msg)
	}
}
