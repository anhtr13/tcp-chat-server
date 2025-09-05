package server

type EVENT string

const (
	RENAME    EVENT = "/name" // Client events
	JOIN_ROOM EVENT = "/join"
	GET_ROOMS EVENT = "/rooms"
	ERROR     EVENT = "/err" // Server event
	MESSAGE   EVENT = "/msg" // Share event
)
