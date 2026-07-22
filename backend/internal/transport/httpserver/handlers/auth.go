package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.Profile, error)
}

type AuthHandler struct {
	authService AuthService
	validator   *validator.Validate
	log         *slog.Logger
}

func NewAuthHandler(as AuthService, v *validator.Validate, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: as, validator: v, log: logger}
}

// Register
func (h *AuthHandler) Register(c *echo.Context) error {
	ctx := c.Request().Context()

	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid json",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "validation failed: " + err.Error(),
		})
	}

	email := req.Email
	password := req.Password

	u, err := h.authService.Register(ctx, email, password)
	if err != nil {
		// TODO: maybe add switch for different types of error, to return different statuses.
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"id":      u.ID,
		"email":   u.Email,
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) Login(c *echo.Context) error {
	// ctx := c.Request().Context()

	return c.JSON(http.StatusOK, map[string]any{
		"status":   "ok",
		"database": "healthy",
	})
}
