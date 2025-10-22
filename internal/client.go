package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
)

type client struct {
	client_id    string
	client_name  string
	current_room *room
	conn         net.Conn
	mtx          sync.RWMutex
}

func new_client(client_name string, client_conn net.Conn) *client {
	id := uuid.New()
	return &client{
		client_id:    id.String(),
		client_name:  client_name,
		current_room: nil,
		conn:         client_conn,
		mtx:          sync.RWMutex{},
	}
}

func (c *client) rename(new_name string) {
	c.mtx.Lock()
	c.client_name = new_name
	c.mtx.Unlock()
}

func (c *client) get_current_room() *room {
	c.mtx.RLock()
	r := c.current_room
	c.mtx.RUnlock()
	return r
}

func (c *client) change_room(r *room) {
	c.mtx.Lock()
	c.current_room = r
	c.mtx.Unlock()
}

func (c *client) write(event EVENT, data string) {
	enc := gob.NewEncoder(c.conn)
	err := enc.Encode(message{
		Event: event,
		Data:  data,
	})
	if err != nil {
		fmt.Println("Error when write to client: ", err.Error())
	}
}
