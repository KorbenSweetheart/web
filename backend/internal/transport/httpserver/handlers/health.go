package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func CheckHealth(c *echo.Context) error {
	// if err := db.Ping(c.Request().Context()); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]any{
	// 		"status":   "error",
	// 		"database": "unhealthy",
	// 	})
	// }

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "healthy",
	})
}
