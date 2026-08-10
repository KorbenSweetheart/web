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

// Connections return a list of profile IDs of the users accepted by the user.
func (cs *ConnectionService) AcceptedConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.AcceptedConnections"
	// log := cs.log.With(slog.String("op", op))

	return nil, nil
}
