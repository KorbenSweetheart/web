package tokenmgr

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidUserID = errors.New("invalid user id")

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	// Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey []byte
	issuer    string
}

// NewTokenManager creates a new instance of TokenManager
func NewTokenManager(secretKey string, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken creates JWT token using HS256 algorithm
func (tm *TokenManager) GenerateToken(userID int64, ttl time.Duration) (string, error) {
	const op = "pkg.tokenmgr.GenerateToken"

	if userID <= 0 {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidUserID)
	}

	claims := CustomClaims{
		UserID: userID,
		// Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			// ID: // "jti" - JWT ID can be a UUID and can be used to revoke the access token immediately (Token Revocation List / Blacklist в Redis)
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign JWT token: %w", op, err)
	}
	return signedToken, nil
}

// GenerateRawToken generates a random cryptographic string of 32 bytes (64 hex symbols)
// Used during login/registration/rotation to send to the client
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	const op = "pkg.tokenmgr.GenerateRefreshToken"

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate raw token: %w", op, err)
	}

	return hex.EncodeToString(token), nil
}

// HashToken calculates the SHA-256 hash of a string
func (tm *TokenManager) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
