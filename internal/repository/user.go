package repository

import "context"

// User represents the core user data exposed to consumers of the repository.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserRepository captures the behaviour required to work with users.
type UserRepository interface {
	List(ctx context.Context) ([]User, error)
}
