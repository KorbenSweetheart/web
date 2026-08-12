package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm"
)

// FindConnectionRecord returns a connection record if it exists between users in any direction.
func (s *Storage) FindConnectionRecord(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error) {
	const op = "storage.postgres.FindConnectionRecord"

	var conn domain.Connection

	// checking does the connection in ANY direction exists
	err := s.db.WithContext(ctx).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			fromUserID, toUserID, toUserID, fromUserID).
		First(&conn).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrConnectionNotFound
		}
		return nil, fmt.Errorf("%s: failed to find connection record fromId: %d toID: %d: %w",
			op, fromUserID, toUserID, err)
	}

	return &conn, nil
}

// CreateConnectionRecord creates connection record with status pending between 2 users.
func (s *Storage) CreateConnectionRecord(ctx context.Context, conn *domain.Connection) error {
	const op = "storage.postgres.CreateConnectionRecord"

	if err := s.db.WithContext(ctx).Create(conn).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrConnectionAlreadyExists
		}
		return fmt.Errorf("%s: failed to create connection record fromId: %d toID: %d: %w",
			op, conn.FromUserID, conn.ToUserID, err)
	}

	return nil
}

// UpdateConnectionRecord changes connection record status between 2 users by changing "pending" status to "accepted" or "declined".
func (s *Storage) UpdateConnectionRecord(ctx context.Context, conn *domain.Connection) error {
	const op = "storage.postgres.UpdateConnectionRecord"

	err := s.db.WithContext(ctx).Save(conn).Error
	if err != nil {
		return fmt.Errorf("%s: failed to update connection record fromId: %d  toID: %d: %w",
			op, conn.FromUserID, conn.ToUserID, err)

	}

	return nil
}

// DeleteConnectionRecord deletes connection record between 2 users.
func (s *Storage) DeleteConnectionRecord(ctx context.Context, conn *domain.Connection) error {
	const op = "storage.postgres.DeleteConnectionRecord"

	err := s.db.WithContext(ctx).Delete(conn).Error
	if err != nil {
		return fmt.Errorf("%s: failed to delete connection record fromId: %d toID: %d: %w",
			op, conn.FromUserID, conn.ToUserID, err)
	}

	return nil
}

// AcceptedConnectionRecords returns a list of IDs of all users' connected profiles with status 'accepted'.
func (s *Storage) AcceptedConnectionRecords(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AcceptedConnectionRecords"

	userIDs := make([]int64, 0)

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Select("CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END", userID).
		Where("(from_user_id = ? OR to_user_id = ?) AND status = ?", userID, userID, domain.Accepted).
		Pluck("case", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get accepted connections for userid: %d: %w", op, userID, err)
	}

	return userIDs, nil
}

// PendingConnectionRecords returns a list of IDs of incoming connection requests for the provided userID.
func (s *Storage) PendingConnectionRecords(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.PendingConnectionRecords"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Where("to_user_id = ? AND status = ?", userID, domain.Pending).
		Pluck("from_user_id", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get pending connections for userid: %d: %w", op, userID, err)
	}

	return userIDs, nil
}

// AllConnectionRecords returns a list of IDs of all types of connections for the provided userID.
func (s *Storage) AllConnectionRecords(ctx context.Context, userID int64) ([]int64, error) {
	const op = "storage.postgres.AllConnectionRecords"

	var userIDs []int64

	err := s.db.WithContext(ctx).Model(&domain.Connection{}).
		Select("CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END", userID).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Pluck("case", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get all connections for userid: %d: %w", op, userID, err)
	}

	return userIDs, nil
}
