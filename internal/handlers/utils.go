package handlers

import (
	"strings"

	"gourl/internal/config"

	"github.com/gin-gonic/gin"
)

// getBaseURL returns the base URL for short links
// Uses BASE_URL from config if set, otherwise detects from request
func getBaseURL(c *gin.Context) string {
	// Get config from context
	cfgInterface, exists := c.Get("config")
	if !exists {
		// Fallback if config not in context
		return getBaseURLFromRequest(c)
	}
	
	cfg, ok := cfgInterface.(*config.Config)
	if !ok || cfg.BaseURL == "" {
		// Use request-based detection
		return getBaseURLFromRequest(c)
	}
	
	// Use configured BASE_URL
	return strings.TrimSuffix(cfg.BaseURL, "/")
}

// getBaseURLFromRequest detects base URL from the HTTP request
func getBaseURLFromRequest(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// Also check X-Forwarded-Proto header (for proxies/load balancers)
	if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
		scheme = "https"
	}
	
	host := c.Request.Host
	// Remove port if it's the default port
	if strings.HasSuffix(host, ":80") && scheme == "http" {
		host = strings.TrimSuffix(host, ":80")
	}
	if strings.HasSuffix(host, ":443") && scheme == "https" {
		host = strings.TrimSuffix(host, ":443")
	}
	
	return scheme + "://" + host
}

