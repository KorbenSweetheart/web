package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ConnectionManager interface {
	ConnectToUser(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error)
	RespondUserConnectionRequest(ctx context.Context, userID, targetUserID int64, status domain.ConnectionStatus) (*domain.Connection, error)
	RemoveConnectionToUser(ctx context.Context, userID, targetUserID int64) error
	AcceptedConnections(ctx context.Context, userID int64) ([]int64, error)
	PendingConnections(ctx context.Context, userID int64) ([]int64, error)
}

type ConnectionHandler struct {
	ConnectionService ConnectionManager
	validator         *validator.Validate
	log               *slog.Logger
}

func NewConnectionHandler(cm ConnectionManager, v *validator.Validate, logger *slog.Logger) *ConnectionHandler {
	return &ConnectionHandler{ConnectionService: cm, validator: v, log: logger}
}

// CreateConnection returns a created connection JSON object between users.
// POST: /connections, body: {"to_user_id": 123}
func (h *ConnectionHandler) CreateConnection(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	var req dto.ConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	connection, err := h.ConnectionService.ConnectToUser(ctx, myID, req.ToUserID)
	if err != nil {
		if errors.Is(err, domain.ErrConnectionAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": domain.ErrConnectionAlreadyExists.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to create connection"})
	}

	return c.JSON(http.StatusCreated, dto.ConnectionResponse{
		FromUserID: connection.FromUserID,
		ToUserID:   connection.ToUserID,
		Status:     connection.Status.String(),
		Timestamp:  connection.UpdatedAt,
	})
}

// UpdateConnectionStatus updates the connection record based on the user's chosen action: accept or decline the request.
// PATCH: /connections/:id Body: {"status": "accepted" or "declined"}
func (h *ConnectionHandler) UpdateConnectionStatus(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	var req dto.UpdateConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	status, err := domain.ParseConnectionStatus(req.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	connection, err := h.ConnectionService.RespondUserConnectionRequest(ctx, myID, targetUserID, status)
	if err != nil {
		if errors.Is(err, domain.ErrConnectionNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": domain.ErrConnectionNotFound.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to update connection status"})
	}

	return c.JSON(http.StatusOK, dto.ConnectionResponse{
		FromUserID: connection.FromUserID,
		ToUserID:   connection.ToUserID,
		Status:     connection.Status.String(),
		Timestamp:  connection.UpdatedAt,
	})
}

// DeleteConnection deletes the connection between two users.
func (h *ConnectionHandler) DeleteConnection(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	if err := h.ConnectionService.RemoveConnectionToUser(ctx, myID, targetUserID); err != nil {
		if errors.Is(err, domain.ErrConnectionNotFound) {
			c.JSON(http.StatusBadRequest, map[string]any{"error": domain.ErrConnectionNotFound.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to delete connection"})
	}

	return c.JSON(http.StatusNoContent, map[string]any{"status": "Connection deleted successfully"}) // TODO: maybe need to return empty response.
}

// Connections return a list of profile IDs of accepted connections for the user.
func (h *ConnectionHandler) Connections(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	connectedUserIDs, err := h.ConnectionService.AcceptedConnections(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get connections"})
	}

	response := make([]dto.ConnectionIDResponse, 0, len(connectedUserIDs))
	for _, id := range connectedUserIDs {
		response = append(response, dto.ConnectionIDResponse{ID: id})
	}

	return c.JSON(http.StatusOK, response)
}

// ConnectionRequests return a list of profile IDs of incoming pending connections for the user.
func (h *ConnectionHandler) ConnectionRequests(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	PendingUserIDs, err := h.ConnectionService.PendingConnections(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get pending connections"})
	}

	response := make([]dto.ConnectionIDResponse, 0, len(PendingUserIDs))
	for _, id := range PendingUserIDs {
		response = append(response, dto.ConnectionIDResponse{ID: id})
	}

	return c.JSON(http.StatusOK, response)
}
