package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func CreateAccount(c *echo.Context) error {
	// c.

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "healthy",
	})
}
