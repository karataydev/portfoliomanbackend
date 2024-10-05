package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/assetquotefeeder"
	"github.com/karataydev/portfoliomanbackend/internal/auth"
	"github.com/karataydev/portfoliomanbackend/internal/config"
	"github.com/karataydev/portfoliomanbackend/internal/database"
	"github.com/karataydev/portfoliomanbackend/internal/investmentgrowth"
	"github.com/karataydev/portfoliomanbackend/internal/param"
	"github.com/karataydev/portfoliomanbackend/internal/portfolio"
	"github.com/karataydev/portfoliomanbackend/internal/transaction"
	"github.com/karataydev/portfoliomanbackend/internal/user"
	"github.com/karataydev/portfoliomanbackend/pkg/scheduler"
	"github.com/svarlamov/goyhfin"
)

type App struct {
	db       *database.DBConnection
	fiberApp *fiber.App
	// dependecie
	paramService *param.Service

	portfolioService *portfolio.Service
	portfolioHandler *portfolio.Handler

	assetService *asset.Service
	assetHandler *asset.Handler

	assetQuoteFeederService *assetquotefeeder.Service
	transactionService      *transaction.Service
	transactionHandler      *transaction.Handler

	userService *user.Service
	userHandler *user.Handler

	tokenService *auth.TokenService

	investmentGrowthService *investmentgrowth.Service
	investmentGrowthHandler *investmentgrowth.Handler

	scheduler *scheduler.Scheduler
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

	assetQuoteChan := make(chan asset.AssetQuoteChanData)
	assetRepo := asset.NewRepository(a.db)
	a.assetService = asset.NewService(assetRepo, assetQuoteChan)
	go a.assetService.AssetQuoteChanDataConsumer()

	a.assetQuoteFeederService = assetquotefeeder.NewService(a.assetService, a.paramService, assetQuoteChan)

	transactionRepo := transaction.NewRepository(a.db)
	a.transactionService = transaction.NewService(transactionRepo, a.assetService)

	portfolioRepo := portfolio.NewRepository(a.db)
	a.portfolioService = portfolio.NewService(portfolioRepo, a.transactionService, a.assetService)

	// Initialize auth services
	rsaKeys, err := auth.NewRSAKeysFromByte([]byte(config.AppConfig.PrivateKey), []byte(config.AppConfig.PublicKey))
	if err != nil {
		log.Fatalf("Failed to initialize RSA keys: %v", err)
	}

	googleValidator := auth.NewGoogleValidator(config.AppConfig.GoogleClientId)

	a.tokenService = auth.NewTokenService(rsaKeys, config.AppConfig.TokenDuration, googleValidator)

	// investment growth service
	a.investmentGrowthService = investmentgrowth.NewService(a.portfolioService, a.assetService)
	a.investmentGrowthHandler = investmentgrowth.NewHandler(a.investmentGrowthService)

	// Initialize user service
	userRepo := user.NewRepository(a.db)
	a.userService = user.NewService(userRepo, a.tokenService)
}

func (a *App) initHandlers() {
	a.portfolioHandler = portfolio.NewHandler(a.portfolioService)
	a.assetHandler = asset.NewHandler(a.assetService)
	a.transactionHandler = transaction.NewHandler(a.transactionService)
	a.userHandler = user.NewHandler(a.userService)
}

func (a *App) setupRoutes() {
	api := a.fiberApp.Group("/api")

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.Post("/signup", a.userHandler.SignUp)
	authGroup.Post("/signin", a.userHandler.SignIn)

	protected := api.Group("")
	protected.Use(auth.JwtAuthMiddleware(a.tokenService))

	protected.Post("/portfolio/add-transaction", a.portfolioHandler.AddTransactionToPortfolio)
	protected.Get("/portfolio/user-portfolios", a.portfolioHandler.GetUserPortfolios)
	protected.Get("/portfolio/followed-portfolios", a.portfolioHandler.GetFollowedPortfolios)

	protected.Get("/portfolio/:portfolioId", a.portfolioHandler.GetPortfolio)
	protected.Get("/portfolio/:portfolioId/allocations", a.portfolioHandler.GetPortfolioWithAllocations)

	protected.Post("/portfolio/:portfolioId/follow", a.portfolioHandler.FollowPortfolio)
	protected.Delete("/portfolio/:portfolioId/unfollow", a.portfolioHandler.UnfollowPortfolio)
	protected.Get("/portfolio/:portfolioId/follower-count", a.portfolioHandler.GetFollowerCount)
	protected.Get("/portfolio/:portfolioId/is-following", a.portfolioHandler.IsFollowing)

	protected.Get("/investment-growth/:symbol", a.investmentGrowthHandler.CalculateInvestmentGrowth)

	protected.Get("/asset", a.assetHandler.GetAsset)
	protected.Get("/asset/market-overview", a.assetHandler.GetMarketOverview)

	protected.Get("/asset/:assetId", a.assetHandler.GetAssets)


	protected.Get("/transaction", a.transactionHandler.Get)
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
	a.assetQuoteFeederService.InsertInitialData()
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
	app.Use(requestid.New())

	return app
}
