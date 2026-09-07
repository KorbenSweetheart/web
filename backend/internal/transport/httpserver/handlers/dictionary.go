package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type DictionaryProvider interface {
	Activities(ctx context.Context) ([]*domain.Activity, error)
}

type DictionaryHandler struct {
	DictionaryService DictionaryProvider
	validator         *validator.Validate
	log               *slog.Logger
}

func NewDictionaryHandler(service DictionaryProvider, v *validator.Validate, logger *slog.Logger) *DictionaryHandler {
	return &DictionaryHandler{
		DictionaryService: service,
		validator:         v,
		log:               logger,
	}
}

// @Activities godoc
// @Summary      Get activities dictionary
// @Description  Get a list of all supported sports and activities
// @Tags         dictionary
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} dto.ActivityResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /activities [get]
func (h *DictionaryHandler) Activities(c *echo.Context) error {
	ctx := c.Request().Context()

	activities, err := h.DictionaryService.Activities(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get activities",
		})
	}

	activitiesResponse := make([]dto.ActivityResponse, 0, len(activities))

	if len(activities) > 0 {
		for i := range activities {
			activity := dto.ActivityResponse{
				ID:    activities[i].ID,
				Title: activities[i].Title,
			}
			activitiesResponse = append(activitiesResponse, activity)
		}
	}

	return c.JSON(http.StatusOK, activitiesResponse)
}
