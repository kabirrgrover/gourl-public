package main

import (
	"log"

	"gourl/pkg/config"
	"gourl/pkg/database"
	"gourl/pkg/handlers"
	"gourl/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)

	// Set up Gin router
	r := gin.Default()

	// Global middleware (applied to all routes)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.ErrorHandler())
	// Store config in context for handlers
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Serve static files (frontend)
	r.Static("/static", "./web/static")
	
	// Serve index.html for root path
	r.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})
	r.GET("/index.html", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})

	// Health check endpoint (no rate limiting)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (no rate limiting, but have their own protection)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// API routes with rate limiting
	api := r.Group("/api")
	api.Use(middleware.RateLimit(rateLimiter))
	{
		// Public endpoints
		api.POST("/shorten", handlers.CreateShortURL) // Optional auth
		api.POST("/shorten/bulk", handlers.BulkCreateShortURL) // Bulk shortening
		api.GET("/stats/:code", handlers.GetStats)
		api.GET("/stats/:code/enhanced", handlers.GetEnhancedStats)
		api.GET("/qr/:code", handlers.GenerateQRCode) // QR code generation
		
		// Protected endpoints (require authentication)
		protected := api.Group("")
		protected.Use(handlers.AuthMiddleware())
		{
			protected.GET("/my-urls", handlers.GetMyURLs)
			protected.GET("/urls/:code", handlers.GetURLDetails)
			protected.DELETE("/urls/:code", handlers.DeleteURL)
		}
	}

	// Redirect route (must be last to catch all codes, but not static files)
	r.GET("/:code", handlers.RedirectURL)

	// Start server
	log.Printf("Server starting on port %s (environment: %s)", cfg.Port, cfg.Environment)
	log.Printf("Rate limit: %d requests/second, burst: %d", cfg.RateLimitRPS, cfg.RateLimitBurst)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


