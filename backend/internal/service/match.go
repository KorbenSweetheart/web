package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
)

type CandidateRepository interface {
	FindCandidates(ctx context.Context, id int64) ([]domain.Profile, error)
}

type RecommendationService struct {
	repo CandidateRepository
	log  *slog.Logger
}

func NewRecommendationService(r CandidateRepository, logger *slog.Logger) *RecommendationService {
	return &RecommendationService{repo: r, log: logger}
}

func (rs *RecommendationService) Matches(ctx context.Context, userID string, limit int) ([]*domain.Profile, error) {
	// op

	// get userprofile and check that profile is complete (all 5 touchpoints are set, and user has activities to check.)
	// send request to dp to find candidates.

	// if len(profile.Activities) == 0 {
	// 	return nil, fmt.Errorf("failed to get profile, id: %d, op: %s, error: %w", id, op, err)
	// }

	return nil, nil
}
