package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type DictionaryRepository interface {
	FindActivities(ctx context.Context) ([]*domain.Activity, error)
}

type DictionaryService struct {
	repo DictionaryRepository
	log  *slog.Logger
}

func NewDictionaryService(repo DictionaryRepository, logger *slog.Logger) *DictionaryService {
	return &DictionaryService{repo: repo, log: logger}
}

// Activities return a list of activities from the dictionary.
func (ds *DictionaryService) Activities(ctx context.Context) ([]*domain.Activity, error) {
	const op = "service.dictionaryService.Activities"
	log := ds.log.With(slog.String("op", op))

	activities, err := ds.repo.FindActivities(ctx)
	if err != nil {
		log.Debug("failed to list of activities", "error", logger.Err(err))
		return nil, err
	}

	return activities, nil
}
