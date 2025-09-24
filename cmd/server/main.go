package main

import (
	"log"
	"net/http"
	"os"

	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/YYx00xZZ/try-12-go/internal/handler"
	postgresrepo "github.com/YYx00xZZ/try-12-go/internal/repository/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Config from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to the configured database
	pg, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pg.Close()

	// Echo app
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())

	// Routes
	e.GET("/health", handler.HealthCheck)
	userRepo := postgresrepo.NewUserRepository(pg)
	userHandler := handler.NewUserHandler(userRepo)
	e.GET("/users", userHandler.GetUsers)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}
