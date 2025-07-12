package server

type EVENT string

const (
	RENAME         EVENT = "/name" // Client events
	JOIN_ROOM            = "/join"
	GET_USER_ROOMS       = "/room/me"
	GET_ALL_ROOMS        = "/room/all"
	ERROR                = "/err" // Server event
	MESSAGE              = "/msg" // Share event
)
