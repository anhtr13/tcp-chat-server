package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

const BUFF_SIZE = 10

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
	defer conn.Close()

	fmt.Printf("%s has connected.\n", client.client_id)

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadBytes('\n')
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
			fmt.Println("Error when read file: ", err.Error())
			return
		}

		msg := message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("Cannot unmarshal payload: ", err.Error())
			return
		}

		event := EVENT(strings.TrimSpace(string(msg.Event)))
		payload := strings.TrimSpace(msg.Payload)

		switch event {
		case RENAME:
			client.rename(payload)
			client.write(MESSAGE, fmt.Sprintf("Your new name is %s", payload))
		case JOIN_ROOM:
			if payload == "" {
				client.write(ERROR, "Invalid room name.")
				continue
			}
			prev_room := client.get_current_room()
			if prev_room != nil {
				prev_room.remove_member(client)
				prev_room.broadcast(fmt.Sprintf("%s has left the room.", client.client_name))
			}
			room := s.get_room(payload)
			room.add_member(client)
			client.change_room(room)
			room.broadcast(fmt.Sprintf("%s has joined room.", client.client_name))
		case MESSAGE:
			room := client.get_current_room()
			if room == nil {
				client.write(ERROR, "You're not in any room, join a room to send message.")
				continue
			}
			room.broadcast(fmt.Sprintf("%s: %s", client.client_name, payload))
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
