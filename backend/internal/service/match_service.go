package service

import (
	"context"
	"match-me-api/internal/domain"
)

// approximate realization
type MatchService struct {
	repo domain.UserRepository
}

func NewSwipeHandler(repo domain.UserRepository) *MatchService {
	return &MatchService{repo: repo}
}

func (ms *MatchService) Recommendations(ctx context.Context, userID string, limit int) ([]*domain.Profile, error) {
	return ms.repo.GetRecommendations(ctx, userID, limit)

	// type request struct {
	// 	FromID string `json:"from_id"`
	// 	ToID   string `json:"to_id"`
	// }

	// var req request
	// if err := c.Bind(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	// }

	// match, err := h.app.HandleMatch(req.FromID, req.ToID)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// }

	// return c.JSON(http.StatusOK, match)
}
