package service

import (
	"context"
	"log/slog"
)

type ConnectionRepository interface {
	AcceptedConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error)
	PendingConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error)
	AllConnectionUserIDs(ctx context.Context, userID int64) ([]int64, error)
}

type ConnectionService struct {
	repo ConnectionRepository
	log  *slog.Logger
}

func NewConnectionService(r ConnectionRepository, logger *slog.Logger) *ConnectionService {
	return &ConnectionService{repo: r, log: logger}
}

// AcceptedConnections return a list of profile IDs of the users accepted by the user.
// GET /connections
func (cs *ConnectionService) AcceptedConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.AcceptedConnections"
	// log := cs.log.With(slog.String("op", op))

	connections, err := cs.repo.AcceptedConnectionUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if connections == nil {
		return []int64{}, nil
	}

	return connections, nil
}

// PendingConnectionRequests return a list of profile IDs of the users accepted by the user.
// GET /connections/requests
func (cs *ConnectionService) PendingConnectionRequests(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.PendingConnectionRequests"
	// log := cs.log.With(slog.String("op", op))

	return nil, nil
}

// AddConnection creates connection with pending status for the provided userID.
// POST /connections
func (cs *ConnectionService) AddConnection(ctx context.Context, fromUserID, toUserID int64) ([]int64, error) {
	const op = "service.connectionService.AddConnection"
	// log := cs.log.With(slog.String("op", op))

	return nil, nil
}

// ChangeConnectionStatus changes pending connection request status to accepted or dismissed.
// PUT /connections/:id
func (cs *ConnectionService) ChangeConnectionStatus(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.ChangeConnectionStatus"
	// log := cs.log.With(slog.String("op", op))

	return nil, nil
}

// DeleteConnection deletes connection with status accepted between users.
// DELETE /connections/:id
func (cs *ConnectionService) DeleteConnection(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.DeleteConnection"
	// log := cs.log.With(slog.String("op", op))

	return nil, nil
}
