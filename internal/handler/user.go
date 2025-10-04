package handler

import (
	"log/slog"
	"net/http"

	"github.com/YYx00xZZ/try-12-go/internal/repository"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := otel.Tracer("handler.user").Start(c.Request().Context(), "UserHandler.GetUsers")
	defer span.End()

	users, err := h.repo.List(ctx)
	if err != nil {
		slog.Error("failed to list users", slog.Any("err", err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	span.SetAttributes(attribute.Int("user.count", len(users)))
	slog.Info("listed users", slog.Int("count", len(users)))

	return c.JSON(http.StatusOK, users)
}
