package service

import (
	"context"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type CandidateRepository interface {
	FindCandidates(ctx context.Context, id int64) ([]*domain.Profile, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
}

type RecommendationService struct {
	repo CandidateRepository
	log  *slog.Logger
}

func NewRecommendationService(r CandidateRepository, logger *slog.Logger) *RecommendationService {
	return &RecommendationService{repo: r, log: logger}
}

func (rs *RecommendationService) Recommendations(ctx context.Context, userID int64) ([]*domain.Profile, error) {
	const op = "service.recommendationService.Recommendations"
	log := rs.log.With(slog.String("op", op))

	// 1. get userprofile and check that profile is complete (all 5 touchpoints are set, and user has activities to check.)
	p, err := rs.repo.ProfileByID(ctx, userID)
	if err != nil {
		log.Debug("failed to get profile by id", "id", userID, "error", logger.Err(err))
		return nil, err
	}

	isComplete := isProfileComplete(p)
	if !isComplete {
		return nil, domain.ErrIncompleteProfile
	}

	activityIDs := make([]int64, 0, len(p.Activities))
	for i, _ := range p.Activities {
		activityIDs = append(activityIDs, p.Activities[i].ActivityID)
	}

	// 2. send request to db to find candidates.
	candidates, err := FindCandidates(ctx, p.UserID, p.MaxRadius, p.Lon, p.Lat, activityIDs)([]*domain.Profile, error)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates profiles, id: %d, op: %s, error: %w", userID, op, err)
	}

	// 3. Filter and sort the candidates based on some criteria or weights

	return candidates, nil
}

func isProfileComplete(p *domain.Profile) bool {
	// min 5 touch points:
	// 1. distance - calculated inside db (postgis index) - not nill
	// 2. activities (match by 3 touch points: type, experience level, interest level)
	// 3. interaction mode - has default value - not nill
	// 4. age - potential addition
	// 5. gender - potential addition

	if len(p.Activities) == 0 {
		return false
	}

	return true
}
