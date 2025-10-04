package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/YYx00xZZ/try-12-go/docs"
	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/YYx00xZZ/try-12-go/internal/handler"
	"github.com/YYx00xZZ/try-12-go/internal/observability"
	"github.com/YYx00xZZ/try-12-go/internal/repository"
	mongorepo "github.com/YYx00xZZ/try-12-go/internal/repository/mongo"
	postgresrepo "github.com/YYx00xZZ/try-12-go/internal/repository/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Try 12 Go API
// @version 1.0
// @description API for managing users
// @BasePath /
func main() {
	ctx := context.Background()
	shutdown, logger, err := observability.Setup(ctx, "try-12-go")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	backend := os.Getenv("DB_BACKEND")
	if backend == "" {
		backend = "postgres"
	}
	backend = strings.ToLower(backend)

	var (
		userRepo    repository.UserRepository
		repoCleanup func()
	)

	switch backend {
	case "postgres":
		pg, err := db.NewDB()
		if err != nil {
			slog.Error("failed to connect to postgres", slog.Any("err", err))
			os.Exit(1)
		}
		repoCleanup = func() {
			pg.Close()
		}
		userRepo = postgresrepo.NewUserRepository(pg)
	case "mongo":
		mongoClient, mongoCfg, err := db.NewMongoClient(ctx)
		if err != nil {
			slog.Error("failed to connect to mongodb", slog.Any("err", err))
			os.Exit(1)
		}
		repoCleanup = func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := mongoClient.Disconnect(shutdownCtx); err != nil {
				slog.Error("failed to disconnect mongodb", slog.Any("err", err))
			}
		}
		collection := mongoClient.Database(mongoCfg.Database).Collection(mongoCfg.Collection)
		userRepo = mongorepo.NewUserRepository(collection)
	default:
		logger.Error("unsupported DB_BACKEND", slog.String("backend", backend))
		os.Exit(1)
	}

	if repoCleanup != nil {
		defer repoCleanup()
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogLatency:  true,
		LogMethod:   true,
		LogRemoteIP: true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
				slog.String("remote_ip", v.RemoteIP),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
			}
			slog.LogAttrs(c.Request().Context(), slog.LevelInfo, "http_request", attrs...)
			return nil
		},
	}))
	e.Use(observability.TraceMiddleware("try-12-go"))

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/health", handler.HealthCheck)
	userHandler := handler.NewUserHandler(userRepo)
	e.GET("/users", userHandler.GetUsers)

	slog.Info("starting server", slog.String("port", port), slog.String("db_backend", backend))
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slog.Error("server shutdown", slog.Any("err", err))
		os.Exit(1)
	}
}
