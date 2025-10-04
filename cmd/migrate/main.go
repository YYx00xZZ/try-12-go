package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	backend := strings.ToLower(os.Getenv("DB_BACKEND"))
	if backend != "" && backend != "postgres" {
		logger.Info("skipping migrations for non-postgres backend", slog.String("backend", backend))
		return
	}

	cfg, err := db.LoadConfig()
	if err != nil {
		slog.Error("database config invalid", slog.Any("err", err))
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DSN,
	)
	if err != nil {
		slog.Error("migration init failed", slog.Any("err", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migration failed", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("migrations applied successfully")
}
