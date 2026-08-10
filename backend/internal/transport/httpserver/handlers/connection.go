package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ConnectionManager interface {
	AcceptedConnections(ctx context.Context, userID int64) ([]int64, error)
	// FriendRequest(ctx context.Context, fromID, toID int64) error
}

type ConnectionHandler struct {
	ConnectionService ConnectionManager
	validator         *validator.Validate
	log               *slog.Logger
}

func NewConnectionHandler(cm ConnectionManager, v *validator.Validate, logger *slog.Logger) *ConnectionHandler {
	return &ConnectionHandler{ConnectionService: cm, validator: v, log: logger}
}

// Connections return a list of profile IDs of accepted connections for the user.
func (h *ConnectionHandler) Connections(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	connections, err := h.ConnectionService.AcceptedConnections(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get recommendations"})

	}

	return c.JSON(http.StatusOK, map[string]any{
		"connections": connections,
	})
}
