package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rps      int // requests per second
	burst    int // burst size
	cleanup  *time.Ticker
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rps,
		burst:    burst,
		cleanup:  time.NewTicker(5 * time.Minute), // Clean up old visitors every 5 minutes
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors removes old visitors periodically
func (rl *RateLimiter) cleanupVisitors() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:   rl.burst - 1,
			lastSeen: now,
		}
		return true
	}

	// Refill tokens based on time passed
	elapsed := now.Sub(v.lastSeen).Seconds()
	tokensToAdd := int(elapsed * float64(rl.rps))
	if tokensToAdd > 0 {
		v.tokens = min(v.tokens+tokensToAdd, rl.burst)
		v.lastSeen = now
	}

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// RateLimit middleware limits requests per IP
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (for proxies)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}

// ErrorHandler middleware handles panics and errors gracefully
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Log the error
			gin.DefaultErrorWriter.Write([]byte(err.Error() + "\n"))
			
			// Return appropriate error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

