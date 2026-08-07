package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	dbPinger  Pinger
	validator *validator.Validate
	log       *slog.Logger
}

func NewHealthHandler(p Pinger) *HealthHandler {
	return &HealthHandler{dbPinger: p}
}

func (h *HealthHandler) Healthz(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *echo.Context) error {
	ctx := c.Request().Context()

	if err := h.dbPinger.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]any{
			"status":   "unavailable",
			"database": "down",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":   "ready",
		"database": "up",
	})
}
