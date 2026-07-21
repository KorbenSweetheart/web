package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.Profile, error)
}

type AuthHandler struct {
	authService AuthService
	log         *slog.Logger
}

func NewAuthHandler(as AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: as, log: logger}
}

// /users/{id}
func (h *AuthHandler) Register(c *echo.Context) error {
	ctx := c.Request().Context()

	user, err := h.authService.Register(ctx, email, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// /users/{id} (id)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(c *echo.Context) error {
	// ctx := c.Request().Context()

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "healthy",
	})
}
