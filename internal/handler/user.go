package handler

import (
	"log/slog"
	"net/http"

	"github.com/YYx00xZZ/try-12-go/internal/repository"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUsers returns the first page of users.
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {array} repository.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.repo.List(ctx)
	if err != nil {
		slog.Error("failed to list users", slog.Any("err", err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	slog.Info("listed users", slog.Int("count", len(users)))

	return c.JSON(http.StatusOK, users)
}
