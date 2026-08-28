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
	DefaultRadius    = 10.0
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_ -]{3,32}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, user *domain.Account) error
	AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	DeleteRefreshTokenByHash(ctx context.Context, hash string) error
	DeleteRefreshTokenByAccountID(ctx context.Context, id int64) error
}

type TokenProcessor interface {
	GenerateToken(userID int64, ttl time.Duration) (string, error)
	GenerateRefreshToken() (string, error)
	HashToken(token string) string
}

type AuthService struct {
	repo            AccountRepository
	tokenRepo       TokenRepository
	tokenMgr        TokenProcessor
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	log             *slog.Logger
}

func NewAuthService(r AccountRepository, tr TokenRepository, tp TokenProcessor, atTTL, rtTTL time.Duration, logger *slog.Logger) *AuthService {
	return &AuthService{repo: r, tokenRepo: tr, tokenMgr: tp, AccessTokenTTL: atTTL, RefreshTokenTTL: rtTTL, log: logger}
}

// Register creates a user account and profile, adds a record to the db table, prepopulates the ID, and returns it to the caller.
func (as *AuthService) Register(ctx context.Context, name, email, password string) (*domain.Account, error) {
	const op = "service.authService.Register"
	log := as.log.With(slog.String("op", op))

	if err := validateRegistrationInput(name, email, password); err != nil {
		log.Debug("invalid registration input")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCostFactor) // Rounds (Cost Factor): 12
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate password hash: %w", op, err)
	}

	u := &domain.Account{
		Email:        email,
		PasswordHash: string(hash),
		Profile: domain.Profile{
			Name:       name,
			Email:      email,
			PictureURL: "",
		},
	}

	log.Info("creating account")

	if err := as.repo.CreateAccount(ctx, u); err != nil {
		if errors.Is(err, domain.ErrEmailIsTaken) {
			log.Debug("email already taken")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return nil, fmt.Errorf("%s: failed to create account: %w", op, err)
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

	account, err := as.repo.AccountByEmail(ctx, email)
	if err != nil {
		log.Debug("failed to get account by email", "email", email, "error", logger.Err(err))
		return "", "", fmt.Errorf("%s: failed to get account by email: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Debug("failed to login: a password and hash do not match", "email", email, "error", logger.Err(err))
			return "", "", domain.ErrInvalidCreds
		} else {
			log.Debug("failed to compare hash and password", "email", email, "error", logger.Err(err))
			return "", "", fmt.Errorf("%s: failed to compare hash and password: %w", op, err)
		}
	}

	log.Debug("generating JWT token")

	accessToken, err := as.tokenMgr.GenerateToken(account.ID, as.AccessTokenTTL)
	if err != nil {
		log.Debug("failed to generate jwt token", "user:", account.ID, "error", logger.Err(err))
		return "", "", fmt.Errorf("%s: failed to generate jwt token: %w", op, err)
	}

	log.Debug("JWT token generated successfully")

	log.Debug("generating refresh token")

	rawRefreshToken, err := as.tokenMgr.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	RefreshToken := &domain.RefreshToken{
		AccountID: account.ID,
		TokenHash: as.tokenMgr.HashToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(as.RefreshTokenTTL),
	}

	if err := as.tokenRepo.SaveRefreshToken(ctx, RefreshToken); err != nil {
		log.Debug("failed to save refresh token to db", "user:", account.ID, "error", logger.Err(err))
		return "", "", fmt.Errorf("%s: failed to save refresh token to db: %w", op, err)
	}

	log.Debug("refresh token generated and saved to db successfully")

	log.Info("login completed successfully")

	return accessToken, rawRefreshToken, nil
}

// Refresh validates the refresh token, performs token rotation, and returns new access and refresh tokens.
func (as *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (string, string, error) {
	const op = "service.authService.Refresh"
	log := as.log.With(slog.String("op", op))

	log.Info("starting token refresh")

	if rawRefreshToken == "" {
		return "", "", domain.ErrInvalidOrExpiredToken
	}

	tokenHash := as.tokenMgr.HashToken(rawRefreshToken)

	rt, err := as.tokenRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidOrExpiredToken) {
			return "", "", err
		}
		return "", "", fmt.Errorf("%s: failed to get refresh token: %w", op, err)
	}

	if rt.ExpiresAt.Before(time.Now()) {
		_ = as.tokenRepo.DeleteRefreshTokenByHash(ctx, tokenHash)
		return "", "", domain.ErrInvalidOrExpiredToken
	}

	accessToken, err := as.tokenMgr.GenerateToken(rt.AccountID, as.AccessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate jwt token: %w", op, err)
	}

	newRawRefreshToken, err := as.tokenMgr.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	newRT := &domain.RefreshToken{
		AccountID: rt.AccountID,
		TokenHash: as.tokenMgr.HashToken(newRawRefreshToken),
		ExpiresAt: time.Now().Add(as.RefreshTokenTTL),
	}

	if err := as.tokenRepo.DeleteRefreshTokenByHash(ctx, tokenHash); err != nil {
		return "", "", fmt.Errorf("%s: failed to delete old refresh token: %w", op, err)
	}

	if err := as.tokenRepo.SaveRefreshToken(ctx, newRT); err != nil {
		return "", "", fmt.Errorf("%s: failed to save new refresh token: %w", op, err)
	}

	log.Info("token refresh completed successfully")

	return accessToken, newRawRefreshToken, nil
}

func (as *AuthService) Logout(ctx context.Context, userID int64) error {
	const op = "service.authService.Logout"
	log := as.log.With(slog.String("op", op))

	log.Info("starting logout")

	// Delete refresh token from DB

	err := as.tokenRepo.DeleteRefreshTokenByAccountID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete refresh token: %w", op, err)
	}

	log.Info("logout completed successfully")

	return nil
}

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
