package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type RecommendationEngine interface { // alternative name RecommendationEngine
	Recommendations(ctx context.Context, userID int64) ([]domain.Profile, error)
	Like(ctx context.Context, fromID, toID int64) (bool, error)
}

type MatchHandler struct {
	recService RecommendationEngine
	validator  *validator.Validate
	log        *slog.Logger
}

func NewMatchHandler(rs RecommendationEngine, v *validator.Validate, logger *slog.Logger) *MatchHandler {
	return &MatchHandler{recService: rs, validator: v, log: logger}
}

// Recommendations return a list of profiles that match the user profile search criteria.
func (h *MatchHandler) Recommendations(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.recService.Profile(ctx, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    myID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, UserSummaryResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}
