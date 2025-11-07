package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"gourl/internal/database"
	"gourl/internal/models"

	"github.com/gin-gonic/gin"
)

// GetMyURLs returns all URLs created by the authenticated user
func GetMyURLs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get URLs for this user
	rows, err := database.DB.Query(
		"SELECT id, code, original_url, created_at FROM urls WHERE user_id = ? ORDER BY created_at DESC",
		id,
	)
	if err != nil {
		log.Printf("Error querying user URLs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	urls := []models.URL{}
	for rows.Next() {
		var url models.URL
		var createdAtStr string
		err := rows.Scan(&url.ID, &url.Code, &url.OriginalURL, &createdAtStr)
		if err != nil {
			log.Printf("Error scanning URL: %v", err)
			continue
		}
		
		// Parse created_at
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			url.CreatedAt = t
		} else if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			url.CreatedAt = t
		} else {
			url.CreatedAt = time.Now()
		}
		
		url.UserID = &id
		urls = append(urls, url)
	}

	c.JSON(http.StatusOK, gin.H{
		"urls": urls,
		"count": len(urls),
	})
}

// DeleteURL deletes a URL (only if user owns it)
func DeleteURL(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Verify user owns this URL
	var urlUserID sql.NullInt64
	err := database.DB.QueryRow(
		"SELECT user_id FROM urls WHERE code = ?",
		code,
	).Scan(&urlUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		log.Printf("Error querying URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check ownership
	if !urlUserID.Valid || urlUserID.Int64 != int64(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this URL"})
		return
	}

	// Delete URL (cascade will delete clicks)
	_, err = database.DB.Exec("DELETE FROM urls WHERE code = ?", code)
	if err != nil {
		log.Printf("Error deleting URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

// GetURLDetails returns detailed information about a URL (if user owns it)
func GetURLDetails(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Get URL details
	var url models.URL
	var urlUserID sql.NullInt64
	var createdAtStr string
	err := database.DB.QueryRow(
		"SELECT id, code, original_url, user_id, created_at FROM urls WHERE code = ?",
		code,
	).Scan(&url.ID, &url.Code, &url.OriginalURL, &urlUserID, &createdAtStr)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		log.Printf("Error querying URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check ownership
	if !urlUserID.Valid || urlUserID.Int64 != int64(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this URL"})
		return
	}

	// Parse created_at
	if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
		url.CreatedAt = t
	} else if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		url.CreatedAt = t
	} else {
		url.CreatedAt = time.Now()
	}

	url.UserID = &id

	// Get click count
	var clickCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM clicks WHERE url_id = ?", url.ID).Scan(&clickCount)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
		"click_count": clickCount,
	})
}

