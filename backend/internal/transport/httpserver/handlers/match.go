package handlers

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"match-me-api/internal/domain"
	"match-me-api/internal/transport/httpserver/dto"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type MatchEngine interface { // alternative name RecommendationEngine
	MatchedProfiles(ctx context.Context, userID int64) ([]domain.ScoredProfile, error)
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

// @Recommendations godoc
// @Summary      Get recommendations
// @Description  Get a list of recommended matching user profiles with compatibility scores
// @Tags         recommendations
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.RecommendationsResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /recommendations [get]
func (h *MatchHandler) Recommendations(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	profiles, err := h.matchService.MatchedProfiles(ctx, myID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, domain.ErrIncompleteProfile):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Incomplete user profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to get recommendations",
			})
		}
	}

	recommendations := make([]dto.RecommendationItemResponse, 0, len(profiles))
	for i := range profiles {
		recommendations = append(recommendations, dto.RecommendationItemResponse{
			UserID: profiles[i].Profile.UserID,
			Score:  int(math.Round(profiles[i].Score * 100)),
		})
	}

	return c.JSON(http.StatusOK, dto.RecommendationsResponse{
		Recommendations: recommendations,
	})
}

// @DismissRecommendation godoc
// @Summary      Dismiss recommendation
// @Description  Dismiss a user recommendation so they will not appear in recommendations again
// @Tags         recommendations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Target user ID to dismiss"
// @Success      200 {object} dto.DismissRecommendationResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /recommendations/{id}/dismiss [post]
func (h *MatchHandler) DismissRecommendation(c *echo.Context) error {
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

	if myID == targetUserID {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
	}

	if err := h.matchService.DismissRecommendation(ctx, myID, targetUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to dismiss recommendation",
		})
	}

	return c.JSON(http.StatusOK, dto.DismissRecommendationResponse{
		Message: "Recommendation dismissed",
	})
}
