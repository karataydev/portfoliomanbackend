package database

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/karataydev/portfoliomanbackend/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DBConnection struct {
	*sqlx.DB
}

func Connect() *DBConnection {
	db, err := sqlx.Connect("postgres", config.AppConfig.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Info("Connected to database successfully")
	return &DBConnection{db}
}

func (db *DBConnection) RunMigrations() error {
	// Create a new Postgres driver for migrate
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %v", err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	return nil
}
