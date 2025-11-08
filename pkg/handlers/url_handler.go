package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"gourl/pkg/database"
	"gourl/pkg/models"
	"gourl/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateShortURL handles POST /shorten requests
func CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate URL
	if !utils.ValidateURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL must start with http:// or https://"})
		return
	}

	// Handle custom code if provided
	var code string
	if req.CustomCode != "" {
		// Validate custom code
		valid, errMsg := utils.ValidateCustomCode(req.CustomCode)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		
		// Check if custom code already exists
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE code = ?)", req.CustomCode).Scan(&exists)
		if err != nil {
			log.Printf("Error checking custom code existence: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "This custom code is already taken"})
			return
		}
		code = req.CustomCode
	} else {
		// Generate unique random code
		var err error
		code, err = utils.GenerateCode(utils.CodeLength)
		if err != nil {
			log.Printf("Error generating code: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
			return
		}

		// Check if code already exists (very unlikely but handle it)
		for i := 0; i < 5; i++ {
			var exists bool
			err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE code = ?)", code).Scan(&exists)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Error checking code existence: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
			if !exists {
				break
			}
			// Regenerate if exists
			code, err = utils.GenerateCode(utils.CodeLength)
			if err != nil {
				log.Printf("Error regenerating code: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
				return
			}
		}
	}

	// Get user ID if authenticated (optional)
	var userID interface{}
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(int); ok {
			userID = id
		}
	} else {
		userID = nil
	}

	// Insert into database
	now := time.Now()
	// Format time for SQLite compatibility
	createdAt := now.Format("2006-01-02 15:04:05")
	
	var expiresAt interface{}
	if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt.Format("2006-01-02 15:04:05")
	} else {
		expiresAt = nil
	}
	
	result, err := database.DB.Exec(
		"INSERT INTO urls (code, original_url, user_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?)",
		code, req.URL, userID, createdAt, expiresAt,
	)
	if err != nil {
		log.Printf("Error inserting URL: %v", err)
		log.Printf("Code: %s, URL: %s, UserID: %v, Time: %v", code, req.URL, userID, now)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL", "details": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	// Build full short URL using configurable base URL
	baseURL := getBaseURL(c)
	shortURL := baseURL + "/" + code

	response := models.CreateURLResponse{
		ShortURL:    shortURL,
		OriginalURL: req.URL,
		Code:        code,
		CreatedAt:   now,
	}

	log.Printf("Created short URL: %s -> %s (ID: %d)", code, req.URL, id)
	c.JSON(http.StatusCreated, response)
}

// RedirectURL handles GET /{code} requests and redirects to original URL
func RedirectURL(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Exclude reserved paths
	reservedPaths := []string{"api", "static", "health", "index.html", "favicon.ico"}
	for _, reserved := range reservedPaths {
		if code == reserved {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
	}

	var urlID int
	var originalURL string
	var expiresAt sql.NullString
	err := database.DB.QueryRow(
		"SELECT id, original_url, expires_at FROM urls WHERE code = ?",
		code,
	).Scan(&urlID, &originalURL, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		log.Printf("Error querying URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if URL has expired
	if expiresAt.Valid && expiresAt.String != "" {
		if expTime, err := time.Parse("2006-01-02 15:04:05", expiresAt.String); err == nil {
			if time.Now().After(expTime) {
				c.JSON(http.StatusGone, gin.H{"error": "This short URL has expired"})
				return
			}
		}
	}

	// Log the click asynchronously (don't block redirect)
	go logClick(urlID, c)

	log.Printf("Redirecting %s -> %s", code, originalURL)
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// logClick logs a click event to the database
func logClick(urlID int, c *gin.Context) {
	// Extract request information
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	
	// Get country from IP (in goroutine to not block redirect)
	go func() {
		country := utils.GetCountryFromIP(ipAddress)
		
		// Insert click into database
		_, err := database.DB.Exec(
			"INSERT INTO clicks (url_id, ip_address, user_agent, referrer, country) VALUES (?, ?, ?, ?, ?)",
			urlID, ipAddress, userAgent, referrer, country,
		)
		if err != nil {
			log.Printf("Error logging click: %v", err)
		}
	}()
}

// GetStats handles GET /api/stats/{code} requests and returns analytics
func GetStats(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Get URL information
	var urlID int
	var originalURL string
	var createdAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, original_url, created_at FROM urls WHERE code = ?",
		code,
	).Scan(&urlID, &originalURL, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		log.Printf("Error querying URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get total clicks count
	var totalClicks int
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM clicks WHERE url_id = ?",
		urlID,
	).Scan(&totalClicks)
	if err != nil {
		log.Printf("Error counting clicks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get unique IP addresses count
	var uniqueIPs int
	err = database.DB.QueryRow(
		"SELECT COUNT(DISTINCT ip_address) FROM clicks WHERE url_id = ?",
		urlID,
	).Scan(&uniqueIPs)
	if err != nil {
		log.Printf("Error counting unique IPs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	response := models.StatsResponse{
		Code:        code,
		OriginalURL: originalURL,
		CreatedAt:   createdAt,
		TotalClicks: totalClicks,
		UniqueIPs:   uniqueIPs,
	}

	c.JSON(http.StatusOK, response)
}

