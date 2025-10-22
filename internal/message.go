package server

type message struct {
	Event   EVENT  `json:"event"`
	Payload string `json:"payload"`
}
