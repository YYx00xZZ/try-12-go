package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	Driver string
	DSN    string
}

// LoadConfig reads database configuration from environment variables and
// returns a driver name plus DSN. Defaults to postgres, but allows overriding
// via DB_DRIVER.
func LoadConfig() (*Config, error) {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		if driver != "postgres" {
			return nil, fmt.Errorf("DATABASE_URL must be set when DB_DRIVER=%s", driver)
		}

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		if host == "" || port == "" || user == "" || name == "" {
			return nil, fmt.Errorf("missing required database environment variables")
		}

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

		dsn = dsnURL.String()
	}

	return &Config{Driver: driver, DSN: dsn}, nil
}

func NewDB() (*sql.DB, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	slog.Info("connecting to database", slog.String("driver", cfg.Driver), slog.String("dsn", maskDSN(cfg.DSN)))

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func maskDSN(dsn string) string {
	if strings.Contains(dsn, "://") {
		if parsed, err := url.Parse(dsn); err == nil {
			if parsed.User != nil {
				username := parsed.User.Username()
				if _, ok := parsed.User.Password(); ok {
					parsed.User = url.UserPassword(username, "****")
				} else {
					parsed.User = url.User(username)
				}
			}
			return parsed.String()
		}
	}

	return dsn
}
