package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	if host != "" && port != "" && user != "" && name != "" {
		dsnURL := &url.URL{
			Scheme: "postgres",
			Host:   fmt.Sprintf("%s:%s", host, port),
			Path:   name,
		}
		if password != "" {
			dsnURL.User = url.UserPassword(user, password)
		} else {
			dsnURL.User = url.User(user)
		}
		query := dsnURL.Query()
		query.Set("sslmode", sslmode)
		dsnURL.RawQuery = query.Encode()

		dsn := dsnURL.String()
		maskedDSN := dsn
		if password != "" {
			maskedDSN = strings.Replace(dsn, password, "****", 1)
		}
		fmt.Println("Connecting with DSN:", maskedDSN)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}

	if urlValue := os.Getenv("DATABASE_URL"); urlValue != "" {
		fmt.Println("Connecting to Postgres using DATABASE_URL")
		db, err := sql.Open("postgres", urlValue)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, fmt.Errorf("missing required database environment variables")
}
