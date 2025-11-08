package handlers

import (
	"fmt"
	"net/http"

	"gourl/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code for a short URL
func GenerateQRCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Get URL to verify it exists
	var originalURL string
	err := database.DB.QueryRow(
		"SELECT original_url FROM urls WHERE code = ?",
		code,
	).Scan(&originalURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Build short URL using configurable base URL
	baseURL := getBaseURL(c)
	shortURL := baseURL + "/" + code

	// Get size parameter (default 256)
	size := 256
	if sizeParam := c.Query("size"); sizeParam != "" {
		var parsedSize int
		if _, err := fmt.Sscanf(sizeParam, "%d", &parsedSize); err == nil {
			size = parsedSize
			if size > 1024 {
				size = 1024 // Max size
			}
			if size < 64 {
				size = 64 // Min size
			}
		}
	}

	// Generate QR code
	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Set size
	qr.DisableBorder = false

	// Convert to PNG
	png, err := qr.PNG(size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code image"})
		return
	}

	// Set headers and return image
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, "image/png", png)
}

