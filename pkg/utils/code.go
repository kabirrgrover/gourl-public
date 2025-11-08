package utils

import (
	"crypto/rand"
	"strings"
)

const (
	// CodeLength is the default length for short URL codes
	CodeLength = 8
	// Base62 characters for URL-safe codes
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenerateCode generates a random base62 code of specified length
func GenerateCode(length int) (string, error) {
	if length <= 0 {
		length = CodeLength
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	code := make([]byte, length)
	for i := range bytes {
		code[i] = base62Chars[bytes[i]%62]
	}

	return string(code), nil
}

// ValidateURL performs basic URL validation
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// ValidateCustomCode validates a custom code/alias
func ValidateCustomCode(code string) (bool, string) {
	if len(code) < 3 {
		return false, "Custom code must be at least 3 characters"
	}
	if len(code) > 20 {
		return false, "Custom code must be at most 20 characters"
	}
	
	// Reserved codes that cannot be used
	reserved := []string{"api", "static", "health", "index.html", "favicon.ico", "admin", "dashboard", "login", "register", "logout"}
	for _, r := range reserved {
		if strings.EqualFold(code, r) {
			return false, "This code is reserved and cannot be used"
		}
	}
	
	// Only allow alphanumeric and hyphens/underscores
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return false, "Custom code can only contain letters, numbers, hyphens, and underscores"
		}
	}
	
	return true, ""
}


