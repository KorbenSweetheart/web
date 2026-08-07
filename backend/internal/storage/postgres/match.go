package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
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
			"ST_DWithin(location, ST_GeogFromText(?), ?) AND ST_Distance(location, ST_GeogFromText(?)) <= (max_radius * 1000)",
			pointWKT,
			maxRadiusMeters,
			pointWKT,
		).
		// The users have similar activities.
		Where("user_id IN (SELECT profile_user_id FROM profile_activities WHERE activity_id IN (?))", activityIDs).
		Find(&candidates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get candidates profiles, id: %d, op: %s, error: %w", profile.UserID, op, err)
	}

	return candidates, nil
}
