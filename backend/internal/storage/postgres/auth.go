package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm"
)

func (s *Storage) SaveRefreshToken(ctx context.Context, rtRecord *domain.RefreshToken) error {
	const op = "storage.postgres.SaveRefreshToken"

	if err := s.db.WithContext(ctx).Create(rtRecord).Error; err != nil {
		return fmt.Errorf("%s: failed to create refresh token: %w", op, err)
	}

	return nil
}

func (s *Storage) GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	const op = "storage.postgres.GetRefreshTokenByHash"

	var rt domain.RefreshToken
	err := s.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&rt).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidOrExpiredToken
		}
		return nil, fmt.Errorf("%s: failed to get refresh token: %w", op, err)
	}

	return &rt, nil
}

func (s *Storage) DeleteRefreshTokenByHash(ctx context.Context, hash string) error {
	const op = "storage.postgres.DeleteRefreshTokenByHash"

	err := s.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		Delete(&domain.RefreshToken{}).Error

	if err != nil {
		return fmt.Errorf("%s: failed to delete refresh token by hash: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteRefreshTokenByAccountID(ctx context.Context, id int64) error {
	const op = "storage.postgres.DeleteRefreshTokenByAccountID"

	err := s.db.WithContext(ctx).
		Where("account_id = ?", id).
		Delete(&domain.RefreshToken{}).Error

	if err != nil {
		return fmt.Errorf("%s: failed to delete refresh token: %w", op, err)
	}

	return nil
}

