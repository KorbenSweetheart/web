package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
)

// CreateConnection .
func (s *Storage) CreateConnection(ctx context.Context, fromUserID, toUserID int64) error {
	const op = "storage.postgres.CreateConnection"

	return nil
}

// UpdateConnection .
func (s *Storage) UpdateConnection(ctx context.Context, fromUserID, toUserID int64, status domain.ConnectionStatus) error {
	const op = "storage.postgres.UpdateConnection"

	return nil
}

// DeleteConnection .
func (s *Storage) DeleteConnection(ctx context.Context, fromUserID, toUserID int64) error {
	const op = "storage.postgres.DeleteConnection"

	return nil
}

// AcceptedConnections returns a list of IDs of all users' connected profiles with status 'accepted'.
func (s *Storage) AcceptedConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AcceptedConnections"

	userIDs := make([]int64, 0)

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Select("CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END", userID).
		Where("(from_user_id = ? OR to_user_id = ?) AND status = ?", userID, userID, domain.Accepted).
		Pluck("case", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get connected userIDs, op: %s, userid: %d, error: %w", op, userID, err)
	}

	return userIDs, nil
}

// PendingConnections returns a list of IDs of incoming connection requests for the provided userID.
func (s *Storage) PendingConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.PendingConnections"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Where("to_user_id = ? AND status = ?", userID, domain.Pending).
		Pluck("from_user_id", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending connection requests, op: %s, userid: %d, error: %w", op, userID, err)
	}

	return userIDs, nil
}

// AllConnections returns a list of IDs of all types of connections for the provided userID.
func (s *Storage) AllConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AllConnections"

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
