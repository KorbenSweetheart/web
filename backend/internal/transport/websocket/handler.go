package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, validate origin against allowed domains
		return true
	},
}

// Handler handles WebSocket connection upgrades for authenticated users.
type Handler struct {
	hub *Hub
	log *slog.Logger
}

// NewHandler creates a new WebSocket Handler instance.
func NewHandler(hub *Hub, logger *slog.Logger) *Handler {
	return &Handler{
		hub: hub,
		log: logger,
	}
}

// Upgrade upgrades the incoming HTTP request to a WebSocket session.
func (h *Handler) Upgrade(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Failed to upgrade to WebSocket"})
	}

	client := NewClient(userID, conn, h.hub, h.log)

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()

	return nil
}
