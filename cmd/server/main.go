package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/YYx00xZZ/try-12-go/internal/handler"
	"github.com/YYx00xZZ/try-12-go/internal/observability"
	postgresrepo "github.com/YYx00xZZ/try-12-go/internal/repository/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	pg, err := db.NewDB()
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("err", err))
		os.Exit(1)
	}
	defer pg.Close()

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
				slog.Duration("latency", v.Latency),
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

	e.GET("/health", handler.HealthCheck)
	userRepo := postgresrepo.NewUserRepository(pg)
	userHandler := handler.NewUserHandler(userRepo)
	e.GET("/users", userHandler.GetUsers)

	slog.Info("starting server", slog.String("port", port))
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slog.Error("server shutdown", slog.Any("err", err))
		os.Exit(1)
	}
}
