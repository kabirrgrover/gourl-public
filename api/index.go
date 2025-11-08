package handler

import (
	"net/http"

	"gourl/pkg/config"
	"gourl/pkg/database"
	"gourl/pkg/handlers"
	"gourl/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/vercel/go-bridge/go/bridge"
)

var router *gin.Engine

func init() {
	// Initialize database when serverless function starts
	if err := database.InitDB(); err != nil {
		panic(err)
	}

	// Set up router once
	setupRouter()
}

func setupRouter() {
	cfg := config.LoadConfig()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)

	router = gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	router.Use(middleware.ErrorHandler())

	router.Static("/static", "./web/static")
	
	// Serve index.html for root path
	router.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})
	router.GET("/index.html", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.RateLimit(rateLimiter))
	{
		api.POST("/shorten", handlers.CreateShortURL)
		api.POST("/shorten/bulk", handlers.BulkCreateShortURL)
		api.GET("/stats/:code", handlers.GetStats)
		api.GET("/stats/:code/enhanced", handlers.GetEnhancedStats)
		api.GET("/qr/:code", handlers.GenerateQRCode)
	}

	// Protected routes (require auth)
	protected := api.Group("")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/my-urls", handlers.GetMyURLs)
		protected.GET("/urls/:code", handlers.GetURLDetails)
		protected.DELETE("/urls/:code", handlers.DeleteURL)
	}

	router.GET("/:code", handlers.RedirectURL)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Serve via Vercel bridge
	bridge.Start(router)
}
