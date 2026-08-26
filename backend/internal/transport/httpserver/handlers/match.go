package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type MatchEngine interface { // alternative name RecommendationEngine
	MatchedProfiles(ctx context.Context, userID int64) ([]*domain.Profile, error)
	DismissRecommendation(ctx context.Context, userID, targetUserID int64) error
	// FriendRequest(ctx context.Context, fromID, toID int64) error
}

type MatchHandler struct {
	matchService MatchEngine
	validator    *validator.Validate
	log          *slog.Logger
}

func NewMatchHandler(ms MatchEngine, v *validator.Validate, logger *slog.Logger) *MatchHandler {
	return &MatchHandler{matchService: ms, validator: v, log: logger}
}

// Recommendations return a list of profiles that match the user profile search criteria.
func (h *MatchHandler) Recommendations(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profiles, err := h.matchService.MatchedProfiles(ctx, myID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    myID,
				"error": "User not found",
			})
		case errors.Is(err, domain.ErrIncompleteProfile):
			return c.JSON(http.StatusBadRequest, map[string]any{
				"id":    myID,
				"error": "Incomplete user profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get recommendations"})
		}
	}

	recommendations := make([]int64, 0, len(profiles))
	if len(profiles) > 0 {
		for i := range profiles {
			recommendations = append(recommendations, profiles[i].UserID)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"recommendations": recommendations,
	})
}

// DismissRecommendation dismisses a recommendation for the authenticated user.
func (h *MatchHandler) DismissRecommendation(c *echo.Context) error {
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

	if myID == targetUserID {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request"})
	}

	if err := h.matchService.DismissRecommendation(ctx, myID, targetUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to dismiss recommendation"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "dismissed",
	})
}

