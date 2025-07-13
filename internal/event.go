package server

type EVENT string

const (
	RENAME    EVENT = "/name" // Client events
	JOIN_ROOM       = "/join"
	GET_ROOMS       = "/rooms"
	ERROR           = "/err" // Server event
	MESSAGE         = "/msg" // Share event
)
