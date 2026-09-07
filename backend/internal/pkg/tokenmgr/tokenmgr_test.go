package tokenmgr

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashToken(t *testing.T) {
	// 1. Setup
	tm := NewTokenManager("secret", "unittest")

	// 2. Table Test Cases
	tests := []struct {
		name     string
		token    string
		wantHash string
	}{
		{
			name:     "valid",
			token:    "test",
			wantHash: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name:     "empty string",
			token:    "",
			wantHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	// 3. Execute Table Tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hashedToken := tm.HashToken(tc.token)

			if hashedToken != tc.wantHash {
				t.Errorf("got hash %s, want %s", hashedToken, tc.wantHash)
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tm := NewTokenManager("secret", "unittest")

	const tokenLen = 64

	refToken1, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(refToken1) != tokenLen {
		t.Errorf("token len %d, want len %d", len(refToken1), tokenLen)
	}

	decoded, err := hex.DecodeString(refToken1)
	if err != nil {
		t.Errorf("token is not a valid hex string: %v", err)
	}
	if len(decoded) != 32 {
		t.Errorf("decoded token length = %d bytes, want 32", len(decoded))
	}

	refToken2, err := tm.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if refToken1 == refToken2 {
		t.Fatalf("expected different hashes, but got identical: %s", refToken2)
	}
}

func TestGenerateToken(t *testing.T) {
	tm := NewTokenManager("secret", "unittest")

	// 2. Table Test Cases
	tests := []struct {
		name         string
		id           int64
		ttl          time.Duration
		wantErr      error
		wantParseErr error
	}{
		{
			name:    "valid standard token",
			id:      7,
			ttl:     time.Minute * 10,
			wantErr: nil,
		},
		{
			name:    "valid large user id",
			id:      math.MaxInt64,
			ttl:     time.Hour,
			wantErr: nil,
		},
		{
			name:    "error - zero user id",
			id:      int64(0),
			ttl:     time.Minute,
			wantErr: ErrInvalidUserID,
		},
		{
			name:    "error - negative user id",
			id:      -10,
			ttl:     time.Minute,
			wantErr: ErrInvalidUserID,
		},
		{
			name:         "expiration in the past",
			id:           1,
			ttl:          -1 * time.Minute,
			wantErr:      nil,
			wantParseErr: jwt.ErrTokenExpired,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := tm.GenerateToken(tc.id, tc.ttl)

			// 1. Generation error assertion
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error generating token: %v", err)
			}
			if token == "" {
				t.Fatal("expected non-empty token")
			}

			// 2. Parse and verify claims
			parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return tm.secretKey, nil
			})

			// 3. Expected parse error assertion (e.g. expired token)
			if tc.wantParseErr != nil {
				if !errors.Is(err, tc.wantParseErr) {
					t.Fatalf("expected parse error %v, got %v", tc.wantParseErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error parsing token: %v", err)
			}

			if !parsedToken.Valid {
				t.Fatalf("invalid token")
			}

			claims, ok := parsedToken.Claims.(*CustomClaims)
			if !ok {
				t.Fatalf("unexpected claims type: %T", parsedToken.Claims)
			}

			if claims.UserID != tc.id {
				t.Errorf("claims.UserID = %d, want %d", claims.UserID, tc.id)
			}
			if claims.Subject != fmt.Sprintf("%d", tc.id) {
				t.Errorf("claims.Subject = %q, want %q", claims.Subject, fmt.Sprintf("%d", tc.id))
			}
			if claims.Issuer != tm.issuer {
				t.Errorf("claims.Issuer = %q, want %q", claims.Issuer, tm.issuer)
			}

			now := time.Now()
			expectedExp := now.Add(tc.ttl)
			if claims.ExpiresAt == nil || claims.ExpiresAt.Sub(expectedExp).Abs() > 2*time.Second {
				t.Errorf("claims.ExpiresAt = %v, want approx %v", claims.ExpiresAt, expectedExp)
			}
			if claims.IssuedAt == nil || claims.IssuedAt.Sub(now).Abs() > 2*time.Second {
				t.Errorf("claims.IssuedAt = %v, want approx %v", claims.IssuedAt, now)
			}
			if claims.NotBefore == nil || claims.NotBefore.Sub(now).Abs() > 2*time.Second {
				t.Errorf("claims.NotBefore = %v, want approx %v", claims.NotBefore, now)
			}
		})
	}

	t.Run("fails verification with wrong secret key", func(t *testing.T) {
		token, err := tm.GenerateToken(42, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}

		_, err = jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("wrong-secret-key"), nil
		})

		if !errors.Is(err, jwt.ErrSignatureInvalid) {
			t.Fatalf("expected ErrSignatureInvalid, got %v", err)
		}
	})
}
