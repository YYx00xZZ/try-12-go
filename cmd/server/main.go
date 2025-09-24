package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/YYx00xZZ/try-12-go/internal/handler"
)

func main() {
	// Read config from environment variables (12-factor principle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // sensible default
	}

	// Create new Echo instance
	e := echo.New()

	// Middleware
	e.HideBanner = true
	e.HidePort = true
	e.Use(handler.RequestLogger)

	// Routes
	e.GET("/health", handler.HealthCheck)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}
