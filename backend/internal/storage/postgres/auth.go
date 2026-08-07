package postgres

import (
	"context"
	"fmt"
	"match-me-api/internal/domain"
)

func (s *Storage) SaveRefreshToken(ctx context.Context, rtRecord *domain.RefreshToken) error {
	const op = "storage.postgres.SaveRefreshToken"
	// log := s.log.With(slog.String("op", op))

	if err := s.db.WithContext(ctx).Create(rtRecord).Error; err != nil {
		return fmt.Errorf("failed to create refresh token, op: %s, error: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteRefreshTokenByAccountID(ctx context.Context, id int64) error {
	const op = "storage.postgres.DeleteRefreshToken"

	err := s.db.WithContext(ctx).
		Where("account_id = ?", id).
		Delete(&domain.RefreshToken{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete refresh token, op: %s, error: %w", op, err)
	}

	return nil
}
