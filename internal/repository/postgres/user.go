package postgres

import (
	"context"
	"database/sql"

	"github.com/YYx00xZZ/try-12-go/internal/repository"
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
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM users ORDER BY id LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]repository.User, 0)
	for rows.Next() {
		var u repository.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
