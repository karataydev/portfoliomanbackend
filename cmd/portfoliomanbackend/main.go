package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/karataydev/portfoliomanbackend/internal/config"
	"github.com/karataydev/portfoliomanbackend/internal/database"
	"github.com/karataydev/portfoliomanbackend/internal/portfolio"
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

	// Initialize repository, service, and handler
	repo := portfolio.NewRepository(db)
	service := portfolio.NewService(repo)
	handler := portfolio.NewHandler(service)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Use CORS middleware globally
	app.Use(cors.New())

	// Use logger middleware
	app.Use(logger.New())

	// Create API group
	api := app.Group("/api")

	// Setup routes
	api.Get("/portfolio/:portfolioID", handler.GetPortfolio)
	api.Get("/portfolio/:portfolioID/allocations", handler.GetPortfolioWithAllocations)

	// Start server
	log.Printf("Starting server on port %s", config.AppConfig.ServerPort)
	if err := app.Listen(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
