package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type CandidateRepository interface {
	FindCandidates(ctx context.Context, profile *domain.Profile) ([]*domain.Profile, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
}

type MatchService struct {
	repo CandidateRepository
	log  *slog.Logger
}

func NewMatchService(r CandidateRepository, logger *slog.Logger) *MatchService {
	return &MatchService{repo: r, log: logger}
}

// MatchedProfiles returns a list of 10 ranked profiles matched to the user based on Match Score.
func (ms *MatchService) MatchedProfiles(ctx context.Context, userID int64) ([]*domain.Profile, error) {
	const op = "service.matchService.Recommendations"
	log := ms.log.With(slog.String("op", op))

	profile, err := ms.repo.ProfileByID(ctx, userID)
	if err != nil {
		log.Debug("failed to get profile by id", "id", userID, "error", logger.Err(err))
		return nil, err
	}

	isComplete := isProfileComplete(profile)
	if !isComplete {
		return nil, domain.ErrIncompleteProfile
	}

	candidates, err := ms.repo.FindCandidates(ctx, profile)
	if err != nil {
		log.Debug("failed to find candidates", "id", userID, "error", logger.Err(err))
		return nil, err
	}

	// TODO: implement Match Score and algorithm
	// add filter and sort the get top 10 best candidates based on Match Score (criteria or weights)
	// don't forget to exclude users who was previously rejected. maybe add it to db level.

	return candidates, nil
}

// isProfileComplete checks that a profile has all required touchpoints filled in.
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
