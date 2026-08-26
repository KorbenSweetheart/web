package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm/clause"
)

// FindCandidates returns a list of profiles that fit the required criteria.
func (s *Storage) FindCandidates(ctx context.Context, profile *domain.Profile) ([]*domain.Profile, error) {
	const op = "storage.postgres.FindCandidates"
	// log := s.log.With(slog.String("op", op))

	candidates := make([]*domain.Profile, 0)
	maxRadiusMeters := profile.MaxRadius * 1000.0

	activityIDs := make([]int64, 0, len(profile.Activities))
	for i := range profile.Activities {
		activityIDs = append(activityIDs, profile.Activities[i].ActivityID)
	}

	if len(activityIDs) == 0 {
		return candidates, nil
	}

	// WKT - Well-Known Text representation of the geometry/geography: POINT(lon lat)
	pointWKT := fmt.Sprintf("SRID=4326;POINT(%f %f)", profile.Lon, profile.Lat)

	err := s.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Preload("Activities.Activity").
		Where("user_id <> ?", profile.UserID).
		// The users are within each other's radius in PostGIS.
		Where(
			"ST_DWithin(location, ST_GeogFromText(?), ?) AND ST_DWithin(location, ST_GeogFromText(?), max_radius * 1000)",
			pointWKT, maxRadiusMeters, pointWKT,
		).
		Where("EXISTS (SELECT 1 FROM profile_activities WHERE profile_user_id = profiles.user_id AND activity_id IN (?))", activityIDs).
		// exclude users who already connected in any direction, with any status.
		Where("NOT EXISTS (SELECT 1 FROM connections WHERE (from_user_id = ? AND to_user_id = profiles.user_id) OR (from_user_id = profiles.user_id AND to_user_id = ?))",
			profile.UserID, profile.UserID,
		).
		// exclude users where either user has dismissed the other.
		Where("NOT EXISTS (SELECT 1 FROM recommendations WHERE ((from_user_id = ? AND to_user_id = profiles.user_id) OR (from_user_id = profiles.user_id AND to_user_id = ?)) AND status = ?)",
			profile.UserID, profile.UserID, domain.Dismissed,
		).
		Find(&candidates).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get candidates profiles for userid: %d: %w", op, profile.UserID, err)
	}

	return candidates, nil
}

func (s *Storage) IsCandidate(ctx context.Context, userID, targetUserID int64) (bool, error) {
	const op = "storage.postgres.IsCandidate"

	var count int64
	// Evaluates ST_DWithin between userID and targetUserID, plus common activities
	err := s.db.WithContext(ctx).
		Table("profiles AS p1, profiles AS p2").
		Where("p1.user_id = ? AND p2.user_id = ?", userID, targetUserID).
		Where("ST_DWithin(p1.location, p2.location, p1.max_radius * 1000)").
		Where("ST_DWithin(p2.location, p1.location, p2.max_radius * 1000)").
		Where("EXISTS (SELECT 1 FROM profile_activities pa1 JOIN profile_activities pa2 ON pa1.activity_id = pa2.activity_id WHERE pa1.profile_user_id = p1.user_id AND pa2.profile_user_id = p2.user_id)").
		Where("NOT EXISTS (SELECT 1 FROM recommendations WHERE ((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND status = ?)",
			userID, targetUserID, targetUserID, userID, domain.Dismissed,
		).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return count > 0, nil
}

// DismissRecommendationRecord creates or updates a recommendation record with dismissed status.
func (s *Storage) DismissRecommendationRecord(ctx context.Context, rec *domain.Recommendation) error {
	const op = "storage.postgres.DismissRecommendationRecord"

	err := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "from_user_id"}, {Name: "to_user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
		}).
		Create(rec).Error

	if err != nil {
		return fmt.Errorf("%s: failed to dismiss recommendation fromId: %d toID: %d: %w",
			op, rec.FromUserID, rec.ToUserID, err)
	}

	return nil
}

