package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/YYx00xZZ/try-12-go/internal/observability"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func main() {
	ctx := context.Background()
	shutdown, logger, err := observability.Setup(ctx, "try-12-go-migrate")
	if err != nil {
		slog.Error("failed to set up observability", slog.Any("err", err))
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			logger.Error("failed to cleanly shutdown observability", slog.Any("err", err))
		}
	}()

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

	ctx, span := otel.Tracer("migration").Start(ctx, "migrate.up")
	defer span.End()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.Error("migration failed", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("migrations applied successfully")
	span.SetStatus(codes.Ok, "applied")
}
