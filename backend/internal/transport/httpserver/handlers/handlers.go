package handlers

import (
	"match-me-api/internal/transport/websocket"
)

// Handlers combine prepared HTTP handlers to pass them into the router
type Handlers struct {
	Health     *HealthHandler
	Dictionary *DictionaryHandler
	Auth       *AuthHandler
	User       *UserHandler
	Match      *MatchHandler
	Conn       *ConnectionHandler
	Chat       *ChatHandler
	WS         *websocket.Handler
}

