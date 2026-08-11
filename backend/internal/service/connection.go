package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
)

type ConnectionRepository interface {
	CreateConnectionRecord(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error)
	UpdateConnectionRecord(ctx context.Context, userID, targetUserID int64, status domain.ConnectionStatus) (*domain.Connection, error)
	DeleteConnectionRecord(ctx context.Context, userID, targetUserID int64) error
	AcceptedConnectionRecords(ctx context.Context, userID int64) ([]int64, error)
	PendingConnectionRecords(ctx context.Context, userID int64) ([]int64, error)
	// AllConnectionRecords(ctx context.Context, userID int64) ([]int64, error) // needed for recommendations
}

type ConnectionService struct {
	repo ConnectionRepository
	log  *slog.Logger
}

func NewConnectionService(r ConnectionRepository, logger *slog.Logger) *ConnectionService {
	return &ConnectionService{repo: r, log: logger}
}

// ConnectToUser creates connection with pending status for the provided userID.
func (cs *ConnectionService) ConnectToUser(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error) {
	const op = "service.connectionService.ConnectToUser"
	// log := cs.log.With(slog.String("op", op))

	connection, err := cs.repo.CreateConnectionRecord(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// RespondUserConnectionRequest changes pending connection request status to accepted or dismissed.
func (cs *ConnectionService) RespondUserConnectionRequest(ctx context.Context, userID, targetUserID int64, status domain.ConnectionStatus) (*domain.Connection, error) {
	const op = "service.connectionService.RespondUserConnectionRequest"
	// log := cs.log.With(slog.String("op", op))

	connection, err := cs.repo.UpdateConnectionRecord(ctx, userID, targetUserID, status)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// RemoveConnectionToUser deletes connection with status accepted between users.
func (cs *ConnectionService) RemoveConnectionToUser(ctx context.Context, fromUserID, toUserID int64) error {
	const op = "service.connectionService.RemoveConnectionToUser"
	// log := cs.log.With(slog.String("op", op))

	return nil
}

// AcceptedConnections return a list of profile IDs of the users accepted by the user.
func (cs *ConnectionService) AcceptedConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.AcceptedConnections"
	// log := cs.log.With(slog.String("op", op))

	connections, err := cs.repo.AcceptedConnectionRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	if connections == nil {
		return []int64{}, nil
	}

	return connections, nil
}

// PendingConnectionRequests return a list of profile IDs of the users accepted by the user.
func (cs *ConnectionService) PendingConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.PendingConnectionRequests"
	// log := cs.log.With(slog.String("op", op))

	connections, err := cs.repo.PendingConnectionRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	if connections == nil {
		return []int64{}, nil
	}

	return connections, nil
}
