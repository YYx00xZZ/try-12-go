package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck reports the service availability.
// @Summary Health status
// @Tags health
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
