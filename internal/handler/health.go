package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck handler for liveness/readiness probe
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Simple request logger middleware
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Log basic request info to stdout (12-factor logging)
		c.Logger().Infof("%s %s", c.Request().Method, c.Request().RequestURI)
		return next(c)
	}
}
