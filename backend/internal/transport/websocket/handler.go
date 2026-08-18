package websocket

import (
	"net/http"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin determines if the incoming HTTP request is allowed to upgrade
	CheckOrigin: func(r *http.Request) bool {
		// In production, validate origin against allowed domains (e.g., http://localhost:5173 for Vite)
		return true
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) HandleConnection(c *echo.Context) error {
	// 1. Extract authenticated user context
	// If you use standard JWT middleware, user_id is in c.Get("user_id")
	// If using query parameter auth: token := c.QueryParam("token")
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized or missing user context")
	}

	// 2. Upgrade the HTTP connection to a WebSocket stream
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to upgrade to websocket")
	}

	// 3. Instantiate the Client session
	client := &Client{
		ID:     uuid.New(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256), // Buffered channel prevents slow clients from blocking
		Hub:    h.hub,
	}

	// 4. Register client with central Hub
	h.hub.Register <- client

	// 5. Start the read and write pumps in dedicated goroutines
	go client.writePump()
	go client.readPump()

	return nil
}
