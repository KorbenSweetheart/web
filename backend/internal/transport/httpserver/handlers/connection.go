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

// @CreateConnection godoc
// @Summary      Create connection request
// @Description  Create a pending connection request to another user
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ConnectionRequest true "Target user ID"
// @Success      201 {object} dto.ConnectionResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /connections [post]
func (h *ConnectionHandler) CreateConnection(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	var req dto.ConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	connection, err := h.ConnectionService.ConnectToUser(ctx, myID, req.ToUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConnectionAlreadyExists):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Connection already exists",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.ConnectionResponse{
		FromUserID: connection.FromUserID,
		ToUserID:   connection.ToUserID,
		Status:     connection.Status.String(),
		Timestamp:  connection.UpdatedAt,
	})
}

// @UpdateConnectionStatus godoc
// @Summary      Respond to connection request
// @Description  Accept or decline an incoming connection request
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Target user ID"
// @Param        request body dto.UpdateConnectionRequest true "Connection action status ('accepted' or 'declined')"
// @Success      200 {object} dto.ConnectionResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /connections/{id} [patch]
func (h *ConnectionHandler) UpdateConnectionStatus(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user id: NAN",
		})
	}

	var req dto.UpdateConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	status, err := domain.ParseConnectionStatus(req.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid connection status",
		})
	}

	connection, err := h.ConnectionService.RespondUserConnectionRequest(ctx, myID, targetUserID, status)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConnectionNotFound):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Connection not found",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.ConnectionResponse{
		FromUserID: connection.FromUserID,
		ToUserID:   connection.ToUserID,
		Status:     connection.Status.String(),
		Timestamp:  connection.UpdatedAt,
	})
}

// @DeleteConnection godoc
// @Summary      Delete connection
// @Description  Remove an existing connection with another user
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Target user ID"
// @Success      200 {object} dto.OKResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /connections/{id} [delete]
func (h *ConnectionHandler) DeleteConnection(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user id: NAN",
		})
	}

	if err := h.ConnectionService.RemoveConnectionToUser(ctx, myID, targetUserID); err != nil {
		if errors.Is(err, domain.ErrConnectionNotFound) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Connection not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, dto.OKResponse{
		Message: "Connection deleted successfully",
	})
}

// @Connections godoc
// @Summary      Get accepted connections
// @Description  Get a list of user IDs for accepted connections
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} dto.ConnectionIDResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /connections [get]
func (h *ConnectionHandler) Connections(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	connectedUserIDs, err := h.ConnectionService.AcceptedConnections(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	response := make([]dto.ConnectionIDResponse, 0, len(connectedUserIDs))
	for _, id := range connectedUserIDs {
		response = append(response, dto.ConnectionIDResponse{ID: id})
	}

	return c.JSON(http.StatusOK, response)
}

// @ConnectionRequests godoc
// @Summary      Get pending connection requests
// @Description  Get a list of user IDs with incoming pending connection requests
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} dto.ConnectionIDResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /connections/requests [get]
func (h *ConnectionHandler) ConnectionRequests(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	PendingUserIDs, err := h.ConnectionService.PendingConnections(ctx, myID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	response := make([]dto.ConnectionIDResponse, 0, len(PendingUserIDs))
	for _, id := range PendingUserIDs {
		response = append(response, dto.ConnectionIDResponse{ID: id})
	}

	return c.JSON(http.StatusOK, response)
}
