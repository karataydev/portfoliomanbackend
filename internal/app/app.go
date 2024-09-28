package app

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

type App struct {
	db                      *database.DBConnection
	fiberApp                *fiber.App
	paramService            *param.Service
	portfolioService        *portfolio.Service
	portfolioHandler        *portfolio.Handler
	assetService            *asset.Service
	assetHandler            *asset.Handler
	assetQuoteFeederService *assetquotefeeder.Service
	scheduler               *scheduler.Scheduler
}

func New(db *database.DBConnection) *App {
	app := &App{
		db:        db,
		fiberApp:  createFiberApp(),
		scheduler: scheduler.New(),
	}

	app.initServices()
	app.initHandlers()
	app.setupRoutes()
	app.setupScheduler()

	return app
}

func (a *App) initServices() {
	paramRepo := param.NewRepository(a.db)
	a.paramService = param.NewService(paramRepo)

	portfolioRepo := portfolio.NewRepository(a.db)
	a.portfolioService = portfolio.NewService(portfolioRepo)

	assetQuoteChan := make(chan asset.AssetQuoteChanData)
	assetRepo := asset.NewRepository(a.db)
	a.assetService = asset.NewService(assetRepo, assetQuoteChan)
	go a.assetService.AssetQuoteChanDataConsumer()

	a.assetQuoteFeederService = assetquotefeeder.NewService(a.assetService, a.paramService, assetQuoteChan)
}

func (a *App) initHandlers() {
	a.portfolioHandler = portfolio.NewHandler(a.portfolioService)
	a.assetHandler = asset.NewHandler(a.assetService)
}

func (a *App) setupRoutes() {
	api := a.fiberApp.Group("/api")

	api.Get("/portfolio/:portfolioId", a.portfolioHandler.GetPortfolio)
	api.Get("/portfolio/:portfolioId/allocations", a.portfolioHandler.GetPortfolioWithAllocations)
	api.Get("/asset", a.assetHandler.GetAsset)
	api.Get("/asset/:assetId", a.assetHandler.GetAssets)
}

func (a *App) setupScheduler() {
	a.scheduler.Add("daily quote data insert", 2, 22, 0, func() {
		err := a.assetQuoteFeederService.ScrapeAllAssets(goyhfin.OneDay, goyhfin.OneHour)
		if err != nil {
			// Consider using a proper logging package here
			println("error running ScrapeAllAssets:", err.Error())
		}
	})
}

func (a *App) Run() error {
	a.scheduler.Start()
	log.Printf("Starting server on port %s", config.AppConfig.ServerPort)
	return a.fiberApp.Listen(":" + config.AppConfig.ServerPort)
}

func createFiberApp() *fiber.App {
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

	app.Use(cors.New())
	app.Use(logger.New())

	return app
}
