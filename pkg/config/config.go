package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	DBPath          string
	RateLimitRPS    int // Requests per second per IP
	RateLimitBurst  int // Burst size
	CORSAllowedOrigins []string
	Environment     string // "development" or "production"
	BaseURL         string // Base URL for short links (e.g., https://yoursite.com)
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DBPath:          getEnv("DB_PATH", "gourl.db"),
		RateLimitRPS:    getEnvAsInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst:  getEnvAsInt("RATE_LIMIT_BURST", 20),
		Environment:     getEnv("ENV", "development"),
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		BaseURL:         getEnv("BASE_URL", ""), // Empty means auto-detect from request
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice gets an environment variable as a slice (comma-separated) or returns default
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

