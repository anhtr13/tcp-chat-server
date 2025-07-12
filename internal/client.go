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
	mtx          sync.Mutex
}

func newClient(client_name string, client_conn net.Conn) *client {
	id := uuid.New()
	return &client{
		client_id:    id.String(),
		client_name:  client_name,
		current_room: nil,
		conn:         client_conn,
		mtx:          sync.Mutex{},
	}
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
