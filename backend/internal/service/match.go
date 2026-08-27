package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"

	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

const (
	weightDist   = 0.25 // 25% spatial proximity
	weightAct    = 0.65 // 65% activity and skill alignment
	weightMode   = 0.10 // 10% interaction mode compatibility
	alphaAct     = 0.30 // preserves mutual high interest when skills diverge
	topActWeight = 0.60 // weight for primary anchor activity vs breadth

	minActivityScore = 0.15 // minimum activity similarity required to qualify
	minTotalScore    = 0.30 // minimum total match score required for recommendation
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

	isComplete := isProfileComplete(profile)
	if !isComplete {
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

// isProfileComplete checks that a profile has all required touchpoints filled in.
func isProfileComplete(p *domain.Profile) bool {
	// min 5 touch points:
	// 1. distance - calculated inside db (postgis index) - not nill
	// 2. activities (match by 3 touch points: type, experience level, interest level)
	// 3. interaction mode - has default value - not nill
	// 4. age - potential addition
	// 5. gender - potential addition

	if p == nil {
		return false
	}
	if len(p.Activities) == 0 {
		return false
	}
	if p.Lat == 0 && p.Lon == 0 {
		return false
	}
	if p.InteractionMode == 0 {
		return false
	}

	return true
}

// rankCandidates scores, sorts descending, and returns top candidates with their match scores.
func rankCandidates(req *domain.Profile, candidates []*domain.Profile, limit int) []domain.ScoredProfile {
	if len(candidates) == 0 {
		return make([]domain.ScoredProfile, 0)
	}

	// Pre-index requester activities once to prevent allocations inside candidate loop
	reqMap := make(map[int64]domain.ProfileActivity, len(req.Activities))
	for i := range req.Activities {
		reqMap[req.Activities[i].ActivityID] = req.Activities[i]
	}

	scored := make([]domain.ScoredProfile, 0, len(candidates))
	for _, cand := range candidates {
		actScore := calculateMultiActivityScore(reqMap, cand.Activities)
		if actScore < minActivityScore {
			continue
		}

		distKm := haversineDistanceKm(req.Lat, req.Lon, cand.Lat, cand.Lon)
		distScore := calculateDistanceScore(distKm, req.MaxRadius)
		modeScore := calculateInteractionModeScore(req.InteractionMode, cand.InteractionMode)

		totalScore := (weightDist * distScore) + (weightAct * actScore) + (weightMode * modeScore)
		if totalScore < minTotalScore {
			continue
		}

		scored = append(scored, domain.ScoredProfile{
			Profile: cand,
			Score:   totalScore,
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

// calculateDistanceScore computes spatial proximity score using Gaussian decay.
func calculateDistanceScore(distKm, maxRadiusKm float64) float64 {
	if distKm <= 0 {
		return 1.0
	}
	if maxRadiusKm <= 0 || distKm >= maxRadiusKm {
		return 0.0
	}

	sigma := maxRadiusKm / 2.0
	return math.Exp(-0.5 * math.Pow(distKm/sigma, 2))
}

// haversineDistanceKm computes great-circle distance between two GPS coordinates in kilometers.
func haversineDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	sinDLat := math.Sin(dLat / 2.0)
	sinDLon := math.Sin(dLon / 2.0)

	a := sinDLat*sinDLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinDLon*sinDLon
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))

	return earthRadiusKm * c
}

// calculateMultiActivityScore computes multi-activity alignment using a two-tier weighted blend.
func calculateMultiActivityScore(reqMap map[int64]domain.ProfileActivity, candActs []domain.ProfileActivity) float64 {
	if len(reqMap) == 0 || len(candActs) == 0 {
		return 0.0
	}

	maxActScore := 0.0
	sumActScore := 0.0
	sharedCount := 0

	for i := range candActs {
		candAct := candActs[i]
		if reqAct, ok := reqMap[candAct.ActivityID]; ok {
			score := calculateActivityScore(reqAct, candAct)
			if score > maxActScore {
				maxActScore = score
			}
			sumActScore += score
			sharedCount++
		}
	}

	if sharedCount == 0 {
		return 0.0
	}

	unionCount := len(reqMap) + len(candActs) - sharedCount
	if unionCount <= 0 {
		unionCount = 1
	}

	actScore := topActWeight*maxActScore + (1.0-topActWeight)*(sumActScore/math.Sqrt(float64(unionCount)))
	if actScore > 1.0 {
		actScore = 1.0
	}

	return actScore
}

// calculateActivityScore computes normalized similarity for a single shared activity.
func calculateActivityScore(req, cand domain.ProfileActivity) float64 {
	// Interest alignment: range [0.04, 1.0]
	interestScore := float64(int(req.InterestLevel)*int(cand.InterestLevel)) / 25.0

	// Experience penalty: range [0.0, 1.0]
	expDiff := math.Abs(float64(int(req.Experience) - int(cand.Experience)))
	expScore := 1.0 - (expDiff / 4.0)

	return interestScore * (alphaAct + (1.0-alphaAct)*expScore)
}

// calculateInteractionModeScore evaluates compatibility between interaction preferences.
func calculateInteractionModeScore(reqMode, candMode domain.InteractionMode) float64 {
	if reqMode == domain.OpenToAnything || candMode == domain.OpenToAnything || reqMode == candMode {
		return 1.0
	}
	return 0.0
}
