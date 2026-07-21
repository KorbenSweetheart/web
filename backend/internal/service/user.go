package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/storage"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCostFactor = 12

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type UserProvider interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	// UserByID(ctx context.Context, id int64) (*domain.User, error)
	// ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
	// UpdateProfile(ctx context.Context, profile *domain.Profile) error
}

type UserService struct {
	storage UserProvider
	log     *slog.Logger
}

func NewUserService(up UserProvider, logger *slog.Logger) *UserService {
	return &UserService{storage: up, log: logger}
}

// Register creates a user profile, adds a record to the db table, prepopulates the ID, and returns it to the caller.
func (us *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	const op = "service.Register"
	log := us.log.With(slog.String("op", op))

	if err := validateRegistrationInput(email, password); err != nil {
		log.Debug("invalid registration input")
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCostFactor) // Rounds (Cost Factor): 12
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash from password: %w", err)
	}

	u := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	log.Info("registering user")

	if err := us.storage.CreateUser(ctx, u); err != nil {
		if errors.Is(err, storage.ErrEmailIsTaken) {
			log.Debug("email already taken")
			return nil, err
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info("user registered successfully")

	return u, nil
}

func (us *UserService) Login(ctx context.Context, email, password string) (*domain.Profile, error) {
	const op = "service.Login"
	log := us.log.With(slog.String("op", op))

	log.Info("starting login")

	email = strings.ToLower(strings.TrimSpace(email))

	user, err := us.storage.UserByEmail(ctx, email)
	if err != nil {
		log.Debug("login failed", "email", email, "error", err)
		return nil, storage.ErrInvalidCreds
	}

	// TODO: Add salt mandatory part
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Debug("login failed", "email", email, "error", err)
		return nil, storage.ErrInvalidCreds
	}

	log.Info("creating JWT token") // Is it so?

	// TODO: JWT logic here

	log.Info("JWT token created successfully")

	log.Info("login completed successfully")

	return nil, nil
}

// func (us *UserService) Logout(ctx context.Context, JWT string) error {
// 	const op = "service.Logout"
// 	log := us.log.With(slog.String("op", op))

// 	log.Info("starting logout")

// 	// using UUID delete the session record from the table
// 	err := us.storage.DeleteSession(ctx, UUID)
// 	if err != nil {
// 		if errors.Is(err, domain.ErrSessionNotFound) {
// 			return fmt.Errorf("failed to delete session: %w", err)
// 		}
// 		return fmt.Errorf("unexpected error: %w", err)
// 	}

// 	log.Info("logout completed successfully")

// 	return nil
// }

// validateRegistrationInput is a helper function that validates the registration input data.
func validateRegistrationInput(email, pw string) error {
	var errs []error

	// TODO: maybe use for username creation later
	// if !usernameRegex.MatchString(username) {
	// 	errs = append(errs, storage.ErrInvalidUsernameFormat)
	// }

	email = strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(email) {
		errs = append(errs, storage.ErrInvalidEmailFormat)
	}

	if len(pw) < 8 {
		errs = append(errs, storage.ErrInvalidPassFormat)
	}

	return errors.Join(errs...)
}
