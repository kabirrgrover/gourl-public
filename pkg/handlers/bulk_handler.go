package handlers

import (
	"log"
	"net/http"
	"time"

	"gourl/pkg/database"
	"gourl/pkg/models"
	"gourl/pkg/utils"

	"github.com/gin-gonic/gin"
)

// BulkCreateShortURL handles POST /shorten/bulk requests
func BulkCreateShortURL(c *gin.Context) {
	var req models.BulkCreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.URLs) == 0 || len(req.URLs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must provide between 1 and 100 URLs"})
		return
	}

	// Get user ID if authenticated
	var userID interface{}
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(int); ok {
			userID = id
		}
	} else {
		userID = nil
	}

	responses := []models.CreateURLResponse{}
	now := time.Now()
	createdAt := now.Format("2006-01-02 15:04:05")
	baseURL := getBaseURL(c)

	for _, urlReq := range req.URLs {
		// Validate URL
		if !utils.ValidateURL(urlReq.URL) {
			responses = append(responses, models.CreateURLResponse{
				OriginalURL: urlReq.URL,
				Code:        "",
			})
			continue
		}

		// Handle custom code or generate
		var code string
		if urlReq.CustomCode != "" {
			valid, errMsg := utils.ValidateCustomCode(urlReq.CustomCode)
			if !valid {
				log.Printf("Invalid custom code for %s: %s", urlReq.URL, errMsg)
				responses = append(responses, models.CreateURLResponse{
					OriginalURL: urlReq.URL,
					Code:        "",
				})
				continue
			}
			code = urlReq.CustomCode
		} else {
			var err error
			code, err = utils.GenerateCode(utils.CodeLength)
			if err != nil {
				log.Printf("Error generating code: %v", err)
				responses = append(responses, models.CreateURLResponse{
					OriginalURL: urlReq.URL,
					Code:        "",
				})
				continue
			}
		}

		// Check if code exists
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE code = ?)", code).Scan(&exists)
		if err != nil || exists {
			responses = append(responses, models.CreateURLResponse{
				OriginalURL: urlReq.URL,
				Code:        "",
			})
			continue
		}

		// Insert into database
		var expiresAt interface{}
		if urlReq.ExpiresAt != nil {
			expiresAt = urlReq.ExpiresAt.Format("2006-01-02 15:04:05")
		} else {
			expiresAt = nil
		}

		_, err = database.DB.Exec(
			"INSERT INTO urls (code, original_url, user_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?)",
			code, urlReq.URL, userID, createdAt, expiresAt,
		)
		if err != nil {
			log.Printf("Error inserting URL: %v", err)
			responses = append(responses, models.CreateURLResponse{
				OriginalURL: urlReq.URL,
				Code:        "",
			})
			continue
		}

		shortURL := baseURL + "/" + code
		responses = append(responses, models.CreateURLResponse{
			ShortURL:    shortURL,
			OriginalURL: urlReq.URL,
			Code:        code,
			CreatedAt:   now,
		})
	}

	c.JSON(http.StatusCreated, models.BulkCreateURLResponse{
		URLs:  responses,
		Count: len(responses),
	})
}

