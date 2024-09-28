package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/karataydev/portfoliomanbackend/internal/app"
	"github.com/karataydev/portfoliomanbackend/internal/config"
	"github.com/karataydev/portfoliomanbackend/internal/database"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	db := database.Connect()
	defer db.Close()

	// migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	application := app.New(db)

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
