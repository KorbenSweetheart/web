package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm"
)

// CreateAccount adds a single user account record to database.
func (s *Storage) FindCandidates(ctx context.Context, id int64) ([]*domain.Profile, error) {
	const op = "storage.postgres.FindCandidates"
	// log := s.log.With(slog.String("op", op))

	// TODO: do this in service layer
	var profile domain.Profile

	if err := s.db.WithContext(ctx).
		Preload("Activities.Activity").
		Where("user_id = ?", id).
		First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get profile, id: %d, op: %s, error: %w", id, op, err)
		}
	}

	// if len(profile.Activities) == 0 {
	// 	return nil, fmt.Errorf("failed to get profile, id: %d, op: %s, error: %w", id, op, err)
	// }

	// activityIDs := make([]int64, 0, len(profile.Activities))
	// for i, _ := range profile.Activities {
	// 	activityIDs = append(activityIDs, profile.Activities[i].ActivityID)
	// }

	var candidates []*domain.Profile
	MaxRadiusMeters := profile.MaxRadius * 1000.0

	// WKT-строка точки текущего юзера: POINT(lon lat)
	PointWKT := fmt.Sprintf("SRID=4326;POINT(%f %f)", profile.Lon, profile.Lat)

	if err := s.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Preload("Activities.Activity").
		Where("user_id <> ?", profile.UserID).
		// Проверка взаимного радиуса PostGIS
		Where(
			"ST_DWithin(location, ST_GeogFromText(?), ?) AND ST_Distance(location, ST_GeogFromText(?)) <= (max_radius * 1000)",
			PointWKT,
			MaxRadiusMeters,
			PointWKT,
		).
		// Пересечение хотя бы по одному виду спорта
		Where("user_id IN (SELECT profile_user_id FROM profile_activities WHERE activity_id IN (?))", activityIDs).
		Find(&candidates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get profile, id: %d, op: %s, error: %w", id, op, err)
		}
	}

	return candidates, nil
}
