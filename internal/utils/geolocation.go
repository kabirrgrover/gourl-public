package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetCountryFromIP gets country from IP address using free IP geolocation API
func GetCountryFromIP(ipAddress string) string {
	// Skip localhost/private IPs
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || ipAddress == "localhost" {
		return "Local (Testing)"
	}
	
	// Check if it's a private IP address
	ip := net.ParseIP(ipAddress)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return "Local (Private Network)"
		}
	}
	
	// Check common private IP patterns
	if strings.HasPrefix(ipAddress, "192.168.") ||
		strings.HasPrefix(ipAddress, "10.") ||
		strings.HasPrefix(ipAddress, "172.16.") ||
		strings.HasPrefix(ipAddress, "172.17.") ||
		strings.HasPrefix(ipAddress, "172.18.") ||
		strings.HasPrefix(ipAddress, "172.19.") ||
		strings.HasPrefix(ipAddress, "172.20.") ||
		strings.HasPrefix(ipAddress, "172.21.") ||
		strings.HasPrefix(ipAddress, "172.22.") ||
		strings.HasPrefix(ipAddress, "172.23.") ||
		strings.HasPrefix(ipAddress, "172.24.") ||
		strings.HasPrefix(ipAddress, "172.25.") ||
		strings.HasPrefix(ipAddress, "172.26.") ||
		strings.HasPrefix(ipAddress, "172.27.") ||
		strings.HasPrefix(ipAddress, "172.28.") ||
		strings.HasPrefix(ipAddress, "172.29.") ||
		strings.HasPrefix(ipAddress, "172.30.") ||
		strings.HasPrefix(ipAddress, "172.31.") {
		return "Local (Private Network)"
	}
	
	// Use ip-api.com (free, no API key needed, 45 req/min limit)
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=country", ipAddress)
	
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "Unknown"
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown"
	}
	
	var result struct {
		Country string `json:"country"`
		Status  string `json:"status"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return "Unknown"
	}
	
	if result.Status == "success" && result.Country != "" {
		return result.Country
	}
	
	return "Unknown"
}

