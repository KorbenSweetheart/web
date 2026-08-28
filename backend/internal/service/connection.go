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
	ProfilesByIDs(ctx context.Context, ids []int64) ([]*domain.Profile, error)
	IsDismissed(ctx context.Context, userID, targetUserID int64) (bool, error)
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
		return nil, fmt.Errorf("%s: can't connect to self, fromId: %d, toID: %d", op, fromUserID, toUserID)
	}

	existingConn, err := cs.repo.FindConnectionRecord(ctx, fromUserID, toUserID)
	if err != nil && !errors.Is(err, domain.ErrConnectionNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		switch existingConn.Status {
		case domain.Accepted:
			return nil, fmt.Errorf("%s: %w", op, domain.ErrConnectionAlreadyExists)
		case domain.Pending:
			// Case when pending connection already exist from toUserID.
			// Then we can convert it to "accepted".
			if existingConn.ToUserID == fromUserID {
				existingConn.Status = domain.Accepted

				if err := cs.repo.UpdateConnectionRecord(ctx, existingConn); err != nil {
					return nil, fmt.Errorf("%s: %w", op, err)
				}

				return existingConn, nil
			}

			return nil, fmt.Errorf("%s: %w", op, domain.ErrConnectionAlreadyExists)

		case domain.Declined:
			// Case when one of the users declined connection previously.
			return nil, fmt.Errorf("%s: %w", op, domain.ErrConnectionAlreadyExists)
		}
	}

	// Fetch both profiles in a single batch query
	profiles, err := cs.repo.ProfilesByIDs(ctx, []int64{fromUserID, toUserID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(profiles) < 2 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
	}

	var fromProfile, toProfile *domain.Profile
	if profiles[0].UserID == fromUserID {
		fromProfile, toProfile = profiles[0], profiles[1]
	} else {
		fromProfile, toProfile = profiles[1], profiles[0]
	}

	if !domain.IsCandidate(fromProfile, toProfile) {
		return nil, fmt.Errorf("%s: not a candidate: %w", op, domain.ErrConnectionNotAllowed)
	}

	isDismissed, err := cs.repo.IsDismissed(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if isDismissed {
		return nil, fmt.Errorf("%s: recommendation dismissed: %w", op, domain.ErrConnectionNotAllowed)
	}

	conn := &domain.Connection{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     domain.Pending,
	}

	if err = cs.repo.CreateConnectionRecord(ctx, conn); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn.Status = status

	if err := cs.repo.UpdateConnectionRecord(ctx, conn); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return conn, nil
}

// RemoveConnectionToUser deletes connection with status accepted between users.
func (cs *ConnectionService) RemoveConnectionToUser(ctx context.Context, fromUserID, toUserID int64) error {
	const op = "service.connectionService.RemoveConnectionToUser"
	// log := cs.log.With(slog.String("op", op))

	conn, err := cs.repo.FindConnectionRecord(ctx, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := cs.repo.DeleteConnectionRecord(ctx, conn); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// AcceptedConnections return a list of profile IDs of the users accepted by the user.
func (cs *ConnectionService) AcceptedConnections(ctx context.Context, userID int64) ([]int64, error) {
	const op = "service.connectionService.AcceptedConnections"
	// log := cs.log.With(slog.String("op", op))

	connections, err := cs.repo.AcceptedConnectionRecords(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if connections == nil {
		return []int64{}, nil
	}

	return connections, nil
}
