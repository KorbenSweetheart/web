package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
)

type RecommendationProvider interface {
	Recomendations(ctx context.Context, id int64) ([]domain.Profile, error)
}

type MatchService struct {
	storage RecommendationProvider
	log     *slog.Logger
}

func NewSwipeHandler(rp RecommendationProvider, logger *slog.Logger) *MatchService {
	return &MatchService{storage: rp, log: logger}
}

func (ms *MatchService) Recommendations(ctx context.Context, userID string, limit int) ([]int64, error) {
	return nil, nil
	// ms.repo.GetRecommendations(ctx, userID, limit)

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
