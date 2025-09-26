package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/YYx00xZZ/try-12-go/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// UserRepository persists and retrieves users using a SQL database.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository wires a Postgres-backed user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// List returns a page of users, ordered by their primary key.
func (r *UserRepository) List(ctx context.Context) ([]repository.User, error) {
	ctx, span := otel.Tracer("repository.postgres").Start(ctx, "UserRepository.List")
	defer span.End()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM users ORDER BY id LIMIT 10")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.Error("query users failed", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	users := make([]repository.User, 0)
	for rows.Next() {
		var u repository.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			slog.Error("scan user row failed", slog.Any("err", err))
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.Error("rows iteration failed", slog.Any("err", err))
		return nil, err
	}

	span.SetAttributes(attribute.Int("repository.user.count", len(users)))
	slog.Info("fetched users", slog.Int("count", len(users)))
	return users, nil
}
