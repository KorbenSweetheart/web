package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCostFactor = 12
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_ -]{3,32}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type AuthProvider interface {
	CreateAccount(ctx context.Context, user *domain.Account) error
	SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error
	AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type TokenProvider interface {
	GenerateToken(userID int64, ttl time.Duration) (string, error)
	GenerateRefreshToken() (string, error)
	HashToken(token string) string
}

type AuthService struct {
	storage         AuthProvider
	tokenMgr        TokenProvider
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	log             *slog.Logger
}

func NewAuthService(ap AuthProvider, tp TokenProvider, atTTL, rtTTL time.Duration, logger *slog.Logger) *AuthService {
	return &AuthService{storage: ap, tokenMgr: tp, AccessTokenTTL: atTTL, RefreshTokenTTL: rtTTL, log: logger}
}

// Register creates a user account and profile, adds a record to the db table, prepopulates the ID, and returns it to the caller.
func (as *AuthService) Register(ctx context.Context, name, email, password string) (*domain.Account, error) {
	const op = "service.authService.Register"
	log := as.log.With(slog.String("op", op))

	if err := validateRegistrationInput(name, email, password); err != nil {
		log.Debug("invalid registration input")
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCostFactor) // Rounds (Cost Factor): 12
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash from password: %w", err)
	}

	u := &domain.Account{
		Email:        email,
		PasswordHash: string(hash),
		Profile: domain.Profile{
			Name: name, // TODO: maybe add it during registration
		},
	}

	log.Info("creating account")

	if err := as.storage.CreateAccount(ctx, u); err != nil {
		if errors.Is(err, domain.ErrEmailIsTaken) {
			log.Debug("email already taken")
			return nil, err
		}
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	log.Info("account created successfully")

	return u, nil
}

// Login verifies the login attempt; on success, generate a JWT token and return it to the caller.
func (as *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	const op = "service.authService.Login"
	log := as.log.With(slog.String("op", op))

	log.Info("starting login")

	email = strings.ToLower(strings.TrimSpace(email))

	account, err := as.storage.AccountByEmail(ctx, email)
	if err != nil {
		log.Debug("failed to get account by email", "email", email, "error", logger.Err(err))
		return "", "", fmt.Errorf("failed to get account by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Debug("failed to login: a password and hash do not match", "email", email, "error", logger.Err(err))
			return "", "", domain.ErrInvalidCreds
		} else {
			log.Debug("failed to compare hash and password", "email", email, "error", logger.Err(err))
			return "", "", fmt.Errorf("failed to compare hash and password: %w", err)
		}
	}

	log.Debug("generating JWT token")

	accessToken, err := as.tokenMgr.GenerateToken(account.ID, as.AccessTokenTTL)
	if err != nil {
		log.Debug("failed to generate jwt token", "user:", account.ID, "error", logger.Err(err))
		return "", "", fmt.Errorf("failed to generate jwt token: %w", err)
	}

	log.Debug("JWT token generated successfully")

	log.Debug("generating refresh token")

	rawRefreshToken, err := as.tokenMgr.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	RefreshToken := &domain.RefreshToken{
		AccountID: account.ID,
		TokenHash: as.tokenMgr.HashToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(as.RefreshTokenTTL),
	}

	if err := as.storage.SaveRefreshToken(ctx, RefreshToken); err != nil {
		log.Debug("failed to save refresh token to db", "user:", account.ID, "error", logger.Err(err))
		return "", "", fmt.Errorf("failed to save refresh token to db: %w", err)
	}

	log.Debug("refresh token generated and saved to db successfully")

	log.Info("login completed successfully")

	return accessToken, rawRefreshToken, nil
}

// func (us *UserService) Logout(ctx context.Context, JWT string) error {
// 	const op = "service.authService.Logout"
// 	log := us.log.With(slog.String("op", op))

// 	log.Info("starting logout")

// 	// Delete refresh token from DB

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
func validateRegistrationInput(username, email, pw string) error {
	var errs []error

	// TODO: maybe use for username creation later
	if !usernameRegex.MatchString(username) {
		errs = append(errs, domain.ErrInvalidUsernameFormat)
	}

	email = strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(email) {
		errs = append(errs, domain.ErrInvalidEmailFormat)
	}

	if len(pw) < 8 {
		errs = append(errs, domain.ErrInvalidPassFormat)
	}

	return errors.Join(errs...)
}
