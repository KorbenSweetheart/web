package websocket

import (
	"fmt"
	"log/slog"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"time"

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

// @Upgrade godoc
// @Summary      WebSocket upgrade
// @Description  Upgrade HTTP connection to WebSocket protocol for real-time chat and presence
// @Tags         realtime
// @Security     BearerAuth
// @Success      101 "Switching Protocols to WebSocket"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      503 {object} dto.ErrorResponse
// @Router       /ws [get]
func (h *Handler) Upgrade(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Failed to upgrade to WebSocket",
		})
	}

	connID := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	client := NewClient(connID, userID, conn, h.hub, h.log)

	select {
	case h.hub.Register <- client:
	case <-h.hub.done:
		_ = conn.Close()
		return c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Message: "Server is shutting down",
		})
	}

	go client.WritePump()
	go client.ReadPump()

	return nil
}
