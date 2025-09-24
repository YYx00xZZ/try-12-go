package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		fmt.Println("Connecting to Postgres using DATABASE_URL")
		db, err := sql.Open("postgres", url)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)

	fmt.Println("Connecting with DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
