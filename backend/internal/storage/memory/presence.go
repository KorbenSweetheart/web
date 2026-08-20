// Package memory provides in-memory implementations for ephemeral infrastructure storage.
package memory

import (
	"context"
	"sync"
)

// Storage manages active user connection mappings in memory.
type Storage struct {
	mu          sync.RWMutex
	connections map[int64]map[string]struct{}
}

// NewStorage creates a new in-memory presence Storage instance.
func NewStorage() *Storage {
	return &Storage{
		connections: make(map[int64]map[string]struct{}),
	}
}

// AddConnection registers a connection ID for a user.
// Returns isFirst = true if this is the user's first active connection.
func (s *Storage) AddConnection(_ context.Context, userID int64, connID string) (isFirst bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userConns, exists := s.connections[userID]
	if !exists {
		userConns = make(map[string]struct{})
		s.connections[userID] = userConns
	}

	userConns[connID] = struct{}{}
	isFirst = len(userConns) == 1

	return isFirst, nil
}

// RemoveConnection deregisters a connection ID for a user.
// Returns isLast = true if the user has no remaining active connections.
func (s *Storage) RemoveConnection(_ context.Context, userID int64, connID string) (isLast bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userConns, exists := s.connections[userID]
	if !exists {
		return false, nil
	}

	delete(userConns, connID)

	if len(userConns) == 0 {
		delete(s.connections, userID)
		return true, nil
	}

	return false, nil
}

// IsUserOnline checks whether the given user has at least one active connection.
func (s *Storage) IsUserOnline(_ context.Context, userID int64) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conns, exists := s.connections[userID]
	return exists && len(conns) > 0, nil
}

// BatchOnlineStatus retrieves the presence status for multiple users in a single operation.
func (s *Storage) BatchOnlineStatus(_ context.Context, userIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool, len(userIDs))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range userIDs {
		conns, exists := s.connections[id]
		result[id] = exists && len(conns) > 0
	}

	return result, nil
}
