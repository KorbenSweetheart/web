package service

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type CandidateRepository interface {
	FindCandidates(ctx context.Context, profile *domain.Profile) ([]*domain.Profile, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
	DismissRecommendationRecord(ctx context.Context, rec *domain.Recommendation) error
}

type MatchService struct {
	repo CandidateRepository
	log  *slog.Logger
}

func NewMatchService(r CandidateRepository, logger *slog.Logger) *MatchService {
	return &MatchService{repo: r, log: logger}
}

// DismissRecommendation marks a recommendation as dismissed for the user.
func (ms *MatchService) DismissRecommendation(ctx context.Context, userID, targetUserID int64) error {
	const op = "service.matchService.DismissRecommendation"
	log := ms.log.With(slog.String("op", op))

	rec := &domain.Recommendation{
		FromUserID: userID,
		ToUserID:   targetUserID,
		Status:     domain.Dismissed,
	}

	if err := ms.repo.DismissRecommendationRecord(ctx, rec); err != nil {
		log.Debug("failed to dismiss recommendation", "userID", userID, "targetUserID", targetUserID, "error", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// MatchedProfiles returns a list of 10 ranked profiles matched to the user based on Match Score.
func (ms *MatchService) MatchedProfiles(ctx context.Context, userID int64) ([]domain.ScoredProfile, error) {
	const op = "service.matchService.Recommendations"
	log := ms.log.With(slog.String("op", op))

	profile, err := ms.repo.ProfileByID(ctx, userID)
	if err != nil {
		log.Debug("failed to get profile by id", "id", userID, "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !domain.IsProfileComplete(profile) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrIncompleteProfile)
	}

	candidates, err := ms.repo.FindCandidates(ctx, profile)
	if err != nil {
		log.Debug("failed to find candidates", "id", userID, "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ranked := rankCandidates(profile, candidates, 10)
	return ranked, nil
}

// rankCandidates scores, sorts descending, and returns top candidates with their match scores.
func rankCandidates(req *domain.Profile, candidates []*domain.Profile, limit int) []domain.ScoredProfile {
	if len(candidates) == 0 {
		return make([]domain.ScoredProfile, 0)
	}

	scored := make([]domain.ScoredProfile, 0, len(candidates))
	for _, cand := range candidates {
		score, isMatch := domain.CalculateMatchScore(req, cand)
		if !isMatch {
			continue
		}

		scored = append(scored, domain.ScoredProfile{
			Profile: cand,
			Score:   score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if limit > 0 && len(scored) > limit {
		result := make([]domain.ScoredProfile, limit)
		copy(result, scored[:limit])
		return result
	}

	return scored
}
