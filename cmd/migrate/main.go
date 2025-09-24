package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YYx00xZZ/try-12-go/internal/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("database config invalid: %v", err)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DSN,
	)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}

	// Run migrations (up)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}
