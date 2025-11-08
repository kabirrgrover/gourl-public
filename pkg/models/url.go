package models

import "time"

// URL represents a shortened URL in the database
type URL struct {
	ID         int       `json:"id" db:"id"`
	Code       string    `json:"code" db:"code"`
	OriginalURL string   `json:"original_url" db:"original_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UserID     *int      `json:"user_id,omitempty" db:"user_id"` // Optional: for authenticated users
}

// Click represents a click/access event on a shortened URL
type Click struct {
	ID         int       `json:"id" db:"id"`
	URLID      int       `json:"url_id" db:"url_id"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	Referrer   string    `json:"referrer" db:"referrer"`
	ClickedAt  time.Time `json:"clicked_at" db:"clicked_at"`
}

// User represents an API user
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"` // Never expose in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateURLRequest represents the request body for creating a short URL
type CreateURLRequest struct {
	URL        string     `json:"url" binding:"required"`
	CustomCode string     `json:"custom_code,omitempty"` // Optional custom alias
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`  // Optional expiration date
}

// BulkCreateURLRequest represents bulk URL creation
type BulkCreateURLRequest struct {
	URLs []CreateURLRequest `json:"urls" binding:"required,min=1,max=100"`
}

// BulkCreateURLResponse represents bulk creation response
type BulkCreateURLResponse struct {
	URLs  []CreateURLResponse `json:"urls"`
	Count int                 `json:"count"`
}

// CreateURLResponse represents the response after creating a short URL
type CreateURLResponse struct {
	ShortURL   string    `json:"short_url"`
	OriginalURL string   `json:"original_url"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
}

// StatsResponse represents analytics data for a short URL
type StatsResponse struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	TotalClicks int       `json:"total_clicks"`
	UniqueIPs   int       `json:"unique_ips"`
}

// EnhancedStatsResponse includes time-based analytics
type EnhancedStatsResponse struct {
	Code            string            `json:"code"`
	OriginalURL     string            `json:"original_url"`
	CreatedAt       time.Time         `json:"created_at"`
	TotalClicks     int               `json:"total_clicks"`
	UniqueIPs       int               `json:"unique_ips"`
	ClicksByDay     map[string]int    `json:"clicks_by_day"`     // Date -> count
	TopReferrers     []ReferrerStat    `json:"top_referrers"`      // Top 10 referrers
	UserAgents       map[string]int    `json:"user_agents"`        // Browser/device breakdown
	Countries        map[string]int    `json:"countries"`          // Country -> count
}

// ReferrerStat represents referrer statistics
type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
