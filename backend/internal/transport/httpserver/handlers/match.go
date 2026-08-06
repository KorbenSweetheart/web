package handlers

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"

	"github.com/go-playground/validator/v10"
)

type MatchService interface {
	Account(ctx context.Context, id int64) (*domain.Account, error)
	Profile(ctx context.Context, id int64) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
}

type MatchHandler struct {
	matchService MatchService
	validator    *validator.Validate
	log          *slog.Logger
}

func NewMatchHandler(ms MatchService, v *validator.Validate, logger *slog.Logger) *MatchHandler {
	return &MatchHandler{matchService: ms, validator: v, log: logger}
}

// func (mh *MatchHandler) Recom
