package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
)

type ConnectionRepository interface {
	FindConnectionRecord(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error)
	CreateConnectionRecord(ctx context.Context, conn *domain.Connection) error
	UpdateConnectionRecord(ctx context.Context, conn *domain.Connection) error
	DeleteConnectionRecord(ctx context.Context, conn *domain.Connection) error
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

	if fromUserID == toUserID {
		return nil, fmt.Errorf("can't connect to self, op: %s, fromId: %d, toID: %d", op, fromUserID, toUserID)
	}

	var err error

	_, err = cs.repo.FindConnectionRecord(ctx, fromUserID, toUserID)
	if err == nil {
		return nil, domain.ErrConnectionAlreadyExists
	}

	if !errors.Is(err, domain.ErrConnectionNotFound) {
		return nil, err
	}

	conn := &domain.Connection{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     domain.Pending,
	}

	if err = cs.repo.CreateConnectionRecord(ctx, conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// RespondUserConnectionRequest changes pending connection request status to accepted or dismissed.
func (cs *ConnectionService) RespondUserConnectionRequest(ctx context.Context, userID, targetUserID int64, status domain.ConnectionStatus) (*domain.Connection, error) {
	const op = "service.connectionService.RespondUserConnectionRequest"
	// log := cs.log.With(slog.String("op", op))

	var conn *domain.Connection
	var err error

	conn, err = cs.repo.FindConnectionRecord(ctx, targetUserID, userID)
	if err != nil {
		return nil, err
	}

	conn.Status = status

	if err := cs.repo.UpdateConnectionRecord(ctx, conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// RemoveConnectionToUser deletes connection with status accepted between users.
func (cs *ConnectionService) RemoveConnectionToUser(ctx context.Context, fromUserID, toUserID int64) error {
	const op = "service.connectionService.RemoveConnectionToUser"
	// log := cs.log.With(slog.String("op", op))

	conn, err := cs.repo.FindConnectionRecord(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}

	if err := cs.repo.DeleteConnectionRecord(ctx, conn); err != nil {
		return err
	}

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
