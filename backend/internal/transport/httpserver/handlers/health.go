package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/transport/httpserver/dto"
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

// @Healthz godoc
// @Summary      Health check
// @Description  Check if the API server is alive
// @Tags         health
// @Produce      json
// @Success      200 {object} dto.OKResponse
// @Router       /healthz [get]
func (h *HealthHandler) Healthz(c *echo.Context) error {
	return c.JSON(http.StatusOK, dto.OKResponse{
		Message: "API server is up",
	})
}

// @Readyz godoc
// @Summary      Readiness check
// @Description  Check if the API server and database connection are ready
// @Tags         health
// @Produce      json
// @Success      200 {object} dto.OKResponse
// @Failure      503 {object} dto.ErrorResponse
// @Router       /readyz [get]
func (h *HealthHandler) Readyz(c *echo.Context) error {
	ctx := c.Request().Context()

	if err := h.dbPinger.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Message: "Database connection error",
		})
	}

	return c.JSON(http.StatusOK, dto.OKResponse{
		Message: "API server is ready and database is up",
	})
}
