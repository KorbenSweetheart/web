package handlers

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Matcher interface { // alternative name RecommendationEngine
	// Recommendations(ctx context.Context, userID int64) ([]domain.Match, error)
	Like(ctx context.Context, fromID, toID int64) (bool, error)
}

type MatchHandler struct {
	matchService Matcher
	validator    *validator.Validate
	log          *slog.Logger
}

func NewMatchHandler(ms Matcher, v *validator.Validate, logger *slog.Logger) *MatchHandler {
	return &MatchHandler{matchService: ms, validator: v, log: logger}
}

// func (mh *MatchHandler) Recom
