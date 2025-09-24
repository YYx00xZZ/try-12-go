package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/YYx00xZZ/try-12-go/internal/handler"
)

func main() {
	// Config from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to Postgres
	pg, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer pg.Close()

	// Echo app
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(handler.RequestLogger)

	// Routes
	e.GET("/health", handler.HealthCheck)
	userHandler := handler.NewUserHandler(pg)
	e.GET("/users", userHandler.GetUsers)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}
