package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
)

// AcceptedConnectionUserIDs returns a list of IDs of all users' connected profiles with status 'accepted'.
func (s *Storage) AcceptedConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AcceptedConnectionUserIDs"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Select("CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END", userID).
		Where("(from_user_id = ? OR to_user_id = ?) AND status = ?", userID, userID, domain.StatusAccepted).
		Pluck("case", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get connected userIDs, op: %s, userid: %d, error: %w", op, userID, err)
	}

	return userIDs, nil
}

// PendingConnectionUserIDs returns a list of IDs of incoming connection requests for the provided userID.
func (s *Storage) PendingConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.PendingConnectionUserIDs"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Where("to_user_id = ? AND status = ?", userID, domain.StatusPending).
		Pluck("from_user_id", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending connection requests, op: %s, userid: %d, error: %w", op, userID, err)
	}

	return userIDs, nil
}

// AllConnectionUserIDs returns a list of IDs of all types of connections for the provided userID.
func (s *Storage) AllConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AllConnectionUserIDs"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Select("CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END", userID).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Pluck("case", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get all connections for the userID, op: %s, userid: %d, error: %w", op, userID, err)
	}

	return userIDs, nil
}
