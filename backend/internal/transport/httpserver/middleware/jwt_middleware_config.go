package middleware

import (
	"fmt"
	"match-me-api/internal/pkg/tokenmgr"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

// JWTMiddleware creates Echo JWT middleware with config and returns it.
func JWTMiddlewareWithConfig(secretKey, cookieName string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: fmt.Sprintf("header:Authorization:Bearer ,cookie:%s", cookieName),
		// ContextKey:  "jwt_token", // changes "user" to "jwt_token", don't forget to change it belov

		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(tokenmgr.CustomClaims)
		},

		SuccessHandler: func(c *echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token) // if we want to change the name, uncomment ContextKey:  "jwt_token" above,
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
