package utils

import (
	"match-me-api/internal/pkg/tokenmgr"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

// JWTMiddleware creates Echo JWT middleware with config and returns it.
func JWTMiddlewareWithConfig(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "header:Authorization:Bearer ,cookie:access_token",

		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(tokenmgr.CustomClaims)
		},

		SuccessHandler: func(c *echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token) // ContextKey:  "jwt_token", // changes "user" to "jwt_token"
			if !ok || token == nil {
				return nil
			}

			claims, ok := token.Claims.(*tokenmgr.CustomClaims)
			if !ok {
				return nil
			}

			c.Set("user_id", claims.UserID)

			return nil
		},
	})
}
