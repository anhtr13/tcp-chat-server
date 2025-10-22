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
	mtx   sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		rooms: map[string]*room{},
		mtx:   sync.RWMutex{},
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	client := new_client("Anonymous", conn)
	fmt.Printf("%s has connected.\n", client.client_id)
	for {
		dec := gob.NewDecoder(conn)
		msg := message{}
		err := dec.Decode(&msg)

		if err == io.EOF {
			room := client.get_current_room()
			if room != nil {
				room.remove_member(client)
				room.broadcast(fmt.Sprintf("%s has left the room.", client.client_name))
			}
			fmt.Printf("%s has disconnected.\n", client.client_id)
			return
		}
		if err != nil {
			fmt.Printf("Error when handle connection %s: %s.\n", client.client_id, err.Error())
			return
		}

		event := EVENT(strings.TrimSpace(string(msg.Event)))
		data := strings.TrimSpace(msg.Data)

		switch event {
		case RENAME:
			client.rename(data)
			client.write(MESSAGE, fmt.Sprintf("Your new name is %s", data))
		case JOIN_ROOM:
			if data == "" {
				client.write(ERROR, "Invalid room name.")
				continue
			}
			prev_room := client.get_current_room()
			if prev_room != nil {
				prev_room.remove_member(client)
				prev_room.broadcast(fmt.Sprintf("%s has left the room.", client.client_name))
			}
			room := s.get_room(data)
			room.add_member(client)
			client.change_room(room)
			room.broadcast(fmt.Sprintf("%s has joined room.", client.client_name))
		case MESSAGE:
			room := client.get_current_room()
			if room == nil {
				client.write(ERROR, "You're not in any room, join a room to send message.")
				continue
			}
			room.broadcast(fmt.Sprintf("%s: %s", client.client_name, data))
		case GET_ROOMS:
			rooms := s.get_all_rooms()
			resp := fmt.Sprintf("[%s]", strings.Join(rooms, ", "))
			client.write(MESSAGE, resp)
		default:
			client.write(ERROR, fmt.Sprintf("Unknown command: %s", event))
		}
	}
}

func (s *Server) get_room(room_name string) *room {
	s.mtx.Lock()
	r := s.rooms[room_name]
	if r == nil {
		r = new_room(room_name)
		s.rooms[room_name] = r
	}
	s.mtx.Unlock()
	return r
}

func (s *Server) get_all_rooms() []string {
	all_rooms := []string{}
	s.mtx.RLock()
	for _, room := range s.rooms {
		all_rooms = append(all_rooms, room.room_name)
	}
	s.mtx.RUnlock()
	return all_rooms
}
