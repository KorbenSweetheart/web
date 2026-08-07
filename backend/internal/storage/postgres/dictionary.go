package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
)

// FindActivities returns a list of activities from the dictionary.
func (s *Storage) FindActivities(ctx context.Context) ([]*domain.Activity, error) {
	const op = "storage.postgres.FindActivities"

	activities := make([]*domain.Activity, 0)

	err := s.db.WithContext(ctx).Find(&activities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get activities dictionary, op: %s, error: %w", op, err)
	}

	return activities, nil
}
