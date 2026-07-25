package tokenmgr

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	// Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey []byte
	issuer    string
}

func NewTokenManager(secretKey string, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken creates JWT token
func (tm *TokenManager) GenerateToken(userID int64, ttl time.Duration) (string, error) {
	const op = "pkg.tokenmgr.GenerateToken"

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
	return token.SignedString(tm.secretKey)
}

// GenerateRawToken generates a random cryptographic string of 32 bytes (64 hex symbols)
// Used during login/registration/rotation to send to the client
func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate raw token: %w", err)
	}

	return hex.EncodeToString(token), nil
}

// HashToken calculates the SHA-256 hash of a string
func (tm *TokenManager) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ParseToken validates JWT token
// func (tm *TokenManager) ParseToken(accessToken string) (int64, error) {
// 	// ... парсинг и валидация ...
// }
