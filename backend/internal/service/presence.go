// Package service implements core business logic and orchestration for domain entities.
package service

import (
	"context"
	"fmt"
	"log/slog"
)

// PresenceStorage defines the underlying storage contract for user presence.
// Implemented by Memory Storage (now) and Redis Storage (in multi-instance scale).
type PresenceStorage interface {
	AddConnection(ctx context.Context, userID int64, connID string) (isFirst bool, err error)
	RemoveConnection(ctx context.Context, userID int64, connID string) (isLast bool, err error)
	IsUserOnline(ctx context.Context, userID int64) (bool, error)
	BatchOnlineStatus(ctx context.Context, userIDs []int64) (map[int64]bool, error)
}

// PresenceService manages user presence state and business rules.
type PresenceService struct {
	storage PresenceStorage
	log     *slog.Logger
}

// NewPresenceService creates a new PresenceService instance.
func NewPresenceService(storage PresenceStorage, logger *slog.Logger) *PresenceService {
	return &PresenceService{
		storage: storage,
		log:     logger,
	}
}

// UserConnected registers a new connection session for a user.
// Returns isFirst = true if the user transitioned from offline to online.
func (s *PresenceService) UserConnected(ctx context.Context, userID int64, connID string) (bool, error) {
	const op = "service.presenceService.UserConnected"

	isFirst, err := s.storage.AddConnection(ctx, userID, connID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isFirst, nil
}

// UserDisconnected deregisters a connection session for a user.
// Returns isLast = true if the user transitioned from online to offline.
func (s *PresenceService) UserDisconnected(ctx context.Context, userID int64, connID string) (bool, error) {
	const op = "service.presenceService.UserDisconnected"

	isLast, err := s.storage.RemoveConnection(ctx, userID, connID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isLast, nil
}

// IsOnline checks whether a user has at least one active connection.
func (s *PresenceService) IsOnline(ctx context.Context, userID int64) (bool, error) {
	const op = "service.presenceService.IsOnline"

	online, err := s.storage.IsUserOnline(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return online, nil
}

// BatchIsOnline retrieves the presence status for multiple users.
func (s *PresenceService) BatchIsOnline(ctx context.Context, userIDs []int64) (map[int64]bool, error) {
	const op = "service.presenceService.BatchIsOnline"

	if len(userIDs) == 0 {
		return make(map[int64]bool), nil
	}

	statuses, err := s.storage.BatchOnlineStatus(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return statuses, nil
}
