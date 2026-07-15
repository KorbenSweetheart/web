package service

import (
	"context"
	"match-me-api/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	// some logic here
	return us.repo.UserByID(ctx, id)
}
