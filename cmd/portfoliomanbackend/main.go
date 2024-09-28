package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/assetquotefeeder"
	"github.com/karataydev/portfoliomanbackend/internal/config"
	"github.com/karataydev/portfoliomanbackend/internal/database"
	"github.com/karataydev/portfoliomanbackend/internal/param"
	"github.com/karataydev/portfoliomanbackend/internal/portfolio"
	"github.com/karataydev/portfoliomanbackend/pkg/scheduler"
	"github.com/svarlamov/goyhfin"
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

	// Initialize param repository, service
	paramRepo := param.NewRepository(db)
	paramService := param.NewService(paramRepo)

	// Initialize portfolio repository, service, and handler
	portfolioRepo := portfolio.NewRepository(db)
	portfolioService := portfolio.NewService(portfolioRepo)
	portfolioHandler := portfolio.NewHandler(portfolioService)

	// AssetQuoteChanData Channel
	assetQuoteChan := make(chan asset.AssetQuoteChanData)

	// Initialize asset repository, service, and handler
	assetRepo := asset.NewRepository(db)
	assetService := asset.NewService(assetRepo, assetQuoteChan)
	assetHandler := asset.NewHandler(assetService)
	go assetService.AssetQuoteChanDataConsumer()

	// initialize scraper
	assetQuoteFeederService := assetquotefeeder.NewService(assetService, paramService, assetQuoteChan)

	// add initail data if not insreted before
	assetQuoteFeederService.InsertInitialData()
	sched := &scheduler.Scheduler{}
	sched.Add("daily quote data insert", 2, 22, 0, func() {
		inErr := assetQuoteFeederService.ScrapeAllAssets(goyhfin.OneDay, goyhfin.OneHour)
		if inErr != nil {
			log.Printf("error running ScrapeAllAssets %+v", inErr)
		}
	})
	sched.Start()

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
	api.Get("/portfolio/:portfolioId", portfolioHandler.GetPortfolio)
	api.Get("/portfolio/:portfolioId/allocations", portfolioHandler.GetPortfolioWithAllocations)

	// asset routes
	api.Get("/asset", assetHandler.GetAsset)
	api.Get("/asset/:assetId", assetHandler.GetAssets)

	// Start server
	log.Printf("Starting server on port %s", config.AppConfig.ServerPort)
	if err := app.Listen(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
