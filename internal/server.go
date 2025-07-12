package server

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Server struct {
	rooms map[string]*room
	mtx   sync.Mutex
}

func NewServer() *Server {
	return &Server{
		rooms: map[string]*room{},
		mtx:   sync.Mutex{},
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	client := newClient("Anonymous", conn)
	fmt.Printf("%s has connected.\n", client.client_id)
	for {
		dec := gob.NewDecoder(conn)
		msg := message{}
		err := dec.Decode(&msg)
		if err == io.EOF {
			s.handleDisconnect(client)
			fmt.Printf("%s has disconnected.\n", client.client_id)
			return
		}
		if err != nil {
			fmt.Printf("Error when handle connection %s: %s.\n", client.client_id, err.Error())
			return
		}

		// fmt.Printf("Received payload: %s\n", msg)
		event, data := msg.Event, msg.Data
		var respMsg string = ""

		switch event {
		case RENAME:
			s.handleRename(client, data)
			respMsg = fmt.Sprintf("Your new name is %s", data)
		case JOIN_ROOM:
			err = s.handleJoinRoom(client, data)
		case MESSAGE:
			err = s.handleSendMessage(client, data)
		case GET_USER_ROOMS:
			rooms := s.handleGetClientRooms(client)
			respMsg = strings.Join(rooms, ", ")
		case GET_ALL_ROOMS:
			rooms := s.handleGetAllRooms()
			respMsg = strings.Join(rooms, ", ")
		default:
			err = fmt.Errorf("Unknow command: %s", event)
		}

		if err != nil {
			client.write(ERROR, err.Error())
		}
		if respMsg != "" {
			client.write(MESSAGE, respMsg)
		}
	}
}

func (s *Server) handleRename(c *client, newName string) {
	c.mtx.Lock()
	c.client_name = newName
	c.mtx.Unlock()
}

func (s *Server) handleJoinRoom(c *client, roomName string) error {
	if roomName == "" {
		return fmt.Errorf("Invalid room.")
	}

	s.mtx.Lock()
	room := s.rooms[roomName]
	if s.rooms[roomName] == nil {
		room = newRoom(roomName)
		s.rooms[roomName] = room
	}
	s.mtx.Unlock()

	room.mtx.Lock()
	if room.members[c.client_id] != nil {
		room.mtx.Unlock()
		return fmt.Errorf("You're already in the room.")
	}
	room.members[c.client_id] = c
	room.broadcast(fmt.Sprintf("%s has joined room.", c.client_name))
	room.mtx.Unlock()

	c.mtx.Lock()
	prev_room := c.current_room
	if prev_room != nil {
		prev_room.mtx.Lock()
		delete(prev_room.members, c.client_id)
		if len(prev_room.members) == 0 {
			s.mtx.Lock()
			delete(s.rooms, prev_room.room_name)
			s.mtx.Unlock()
		} else {
			prev_room.broadcast(fmt.Sprintf("%s has left the room.", c.client_name))
		}
		prev_room.mtx.Unlock()
	}
	c.current_room = room
	c.mtx.Unlock()

	return nil
}

func (s *Server) handleSendMessage(c *client, msg string) error {
	c.mtx.Lock()
	room := c.current_room
	if room == nil {
		c.mtx.Unlock()
		return fmt.Errorf("You're not in any room, join a room to send message.")
	}
	c.mtx.Unlock()

	msg = fmt.Sprintf("%s: %s", c.client_name, msg)

	room.mtx.Lock()
	room.broadcast(msg)
	room.mtx.Unlock()

	return nil
}

func (s *Server) handleGetClientRooms(c *client) []string {
	client_rooms := []string{}
	for _, room := range s.rooms {
		if room.members[c.client_id] != nil {
			client_rooms = append(client_rooms, room.room_name)
		}
	}
	return client_rooms
}

func (s *Server) handleGetAllRooms() []string {
	all_rooms := []string{}
	for _, room := range s.rooms {
		all_rooms = append(all_rooms, room.room_name)
	}
	return all_rooms
}

func (s *Server) handleDisconnect(c *client) {
	c.mtx.Lock()
	room := c.current_room
	if room == nil {
		c.mtx.Unlock()
		return
	}
	c.mtx.Unlock()

	room.mtx.Lock()
	defer room.mtx.Unlock()

	delete(room.members, c.client_id)
	if len(room.members) == 0 {
		s.mtx.Lock()
		delete(s.rooms, room.room_name)
		s.mtx.Unlock()
	} else {
		room.broadcast(fmt.Sprintf("%s has left the room.", c.client_name))
	}
}
