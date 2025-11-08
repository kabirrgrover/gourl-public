package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"gourl/pkg/database"
	"gourl/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetEnhancedStats returns detailed analytics with time-based stats
func GetEnhancedStats(c *gin.Context) {
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

	// Get clicks by day (last 30 days)
	clicksByDay := make(map[string]int)
	var query string
	if database.IsPostgres() {
		query = `
			SELECT DATE(clicked_at)::text as day, COUNT(*) as count
			FROM clicks
			WHERE url_id = $1 AND clicked_at >= NOW() - INTERVAL '30 days'
			GROUP BY DATE(clicked_at)
			ORDER BY day DESC
		`
	} else {
		query = `
			SELECT DATE(clicked_at) as day, COUNT(*) as count
			FROM clicks
			WHERE url_id = ? AND clicked_at >= datetime('now', '-30 days')
			GROUP BY DATE(clicked_at)
			ORDER BY day DESC
		`
	}
	rows, err := database.DB.Query(query, urlID)
	if err != nil {
		log.Printf("Error querying clicks by day: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var day string
			var count int
			if err := rows.Scan(&day, &count); err == nil {
				clicksByDay[day] = count
			}
		}
	}

	// Get top referrers
	topReferrers := []models.ReferrerStat{}
	if database.IsPostgres() {
		query = `
			SELECT referrer, COUNT(*) as count
			FROM clicks
			WHERE url_id = $1 AND referrer IS NOT NULL AND referrer != ''
			GROUP BY referrer
			ORDER BY count DESC
			LIMIT 10
		`
	} else {
		query = `
			SELECT referrer, COUNT(*) as count
			FROM clicks
			WHERE url_id = ? AND referrer IS NOT NULL AND referrer != ''
			GROUP BY referrer
			ORDER BY count DESC
			LIMIT 10
		`
	}
	rows, err = database.DB.Query(query, urlID)
	if err != nil {
		log.Printf("Error querying referrers: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ref models.ReferrerStat
			if err := rows.Scan(&ref.Referrer, &ref.Count); err == nil {
				topReferrers = append(topReferrers, ref)
			}
		}
	}

	// Get user agents breakdown
	userAgents := make(map[string]int)
	if database.IsPostgres() {
		query = `
			SELECT user_agent, COUNT(*) as count
			FROM clicks
			WHERE url_id = $1 AND user_agent IS NOT NULL AND user_agent != ''
			GROUP BY user_agent
			ORDER BY count DESC
			LIMIT 20
		`
	} else {
		query = `
			SELECT user_agent, COUNT(*) as count
			FROM clicks
			WHERE url_id = ? AND user_agent IS NOT NULL AND user_agent != ''
			GROUP BY user_agent
			ORDER BY count DESC
			LIMIT 20
		`
	}
	rows, err = database.DB.Query(query, urlID)
	if err != nil {
		log.Printf("Error querying user agents: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ua string
			var count int
			if err := rows.Scan(&ua, &count); err == nil {
				// Simplify user agent (extract browser name)
				ua = simplifyUserAgent(ua)
				userAgents[ua] += count
			}
		}
	}

	// Get countries breakdown
	countries := make(map[string]int)
	if database.IsPostgres() {
		query = `
			SELECT country, COUNT(*) as count
			FROM clicks
			WHERE url_id = $1 AND country IS NOT NULL AND country != ''
			GROUP BY country
			ORDER BY count DESC
			LIMIT 20
		`
	} else {
		query = `
			SELECT country, COUNT(*) as count
			FROM clicks
			WHERE url_id = ? AND country IS NOT NULL AND country != ''
			GROUP BY country
			ORDER BY count DESC
			LIMIT 20
		`
	}
	rows, err = database.DB.Query(query, urlID)
	if err != nil {
		log.Printf("Error querying countries: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var country string
			var count int
			if err := rows.Scan(&country, &count); err == nil {
				countries[country] = count
			}
		}
	}

	response := models.EnhancedStatsResponse{
		Code:          code,
		OriginalURL:   originalURL,
		CreatedAt:     createdAt,
		TotalClicks:   totalClicks,
		UniqueIPs:     uniqueIPs,
		ClicksByDay:   clicksByDay,
		TopReferrers:  topReferrers,
		UserAgents:    userAgents,
		Countries:     countries,
	}

	// Always include countries field, even if empty
	if countries == nil {
		response.Countries = make(map[string]int)
	}

	c.JSON(http.StatusOK, response)
}

// simplifyUserAgent extracts browser name from user agent string
func simplifyUserAgent(ua string) string {
	ua = strings.ToLower(ua)
	
	// Detect common browsers
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		return "Chrome"
	}
	if strings.Contains(ua, "firefox") {
		return "Firefox"
	}
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "Safari"
	}
	if strings.Contains(ua, "edg") {
		return "Edge"
	}
	if strings.Contains(ua, "opera") {
		return "Opera"
	}
	if strings.Contains(ua, "curl") {
		return "curl"
	}
	if strings.Contains(ua, "wget") {
		return "wget"
	}
	if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") {
		return "Bot/Crawler"
	}
	
	return "Other"
}

