package application

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CORSConfig defines the configuration for Cross-Origin Resource Sharing (CORS)
type CORSConfig struct {
	// Enabled determines if CORS headers should be added to responses
	Enabled bool

	// AllowedOrigins is a list of origins that are allowed to access the runtime API.
	// Each origin can be an exact match or use wildcards:
	// - "https://example.com" (exact match)
	// - "https://*.example.com" (wildcard subdomain)
	// - "*" (allow all - not recommended for production)
	// If empty and CORS is enabled, no origins will be allowed.
	AllowedOrigins []string

	// AllowedMethods is a list of HTTP methods allowed for CORS requests
	// Default: ["GET", "POST", "OPTIONS"]
	AllowedMethods []string

	// AllowedHeaders is a list of headers that the client is allowed to send
	// Default includes standard Wails headers
	AllowedHeaders []string

	// MaxAge indicates how long the results of a preflight request can be cached
	// Default: 5 minutes
	MaxAge time.Duration
}

// DefaultCORSConfig returns the default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		Enabled: false,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type",
			"X-Wails-Window-ID",
			"X-Wails-Window-Name",
			"X-Wails-Client-ID",
		},
		MaxAge: 5 * time.Minute,
	}
}

// IsOriginAllowed checks if the given origin is allowed based on the CORS configuration
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if !c.Enabled || origin == "" {
		return false
	}

	for _, allowedOrigin := range c.AllowedOrigins {
		// Exact match
		if allowedOrigin == origin {
			return true
		}

		// Wildcard match for all origins (not recommended)
		if allowedOrigin == "*" {
			return true
		}

		// Wildcard subdomain match (e.g., "https://*.example.com")
		if strings.Contains(allowedOrigin, "*") {
			pattern := strings.ReplaceAll(allowedOrigin, "*", "")
			// Check if origin ends with the pattern (for subdomain wildcards)
			if strings.HasPrefix(allowedOrigin, "*") && strings.HasSuffix(origin, pattern) {
				return true
			}
			// Check if origin starts with the pattern (for path wildcards)
			if strings.HasSuffix(allowedOrigin, "*") && strings.HasPrefix(origin, pattern) {
				return true
			}
			// Check for middle wildcard (e.g., "https://*.example.com")
			if strings.HasPrefix(allowedOrigin, "https://*.") || strings.HasPrefix(allowedOrigin, "http://*.") {
				// Extract the domain pattern after the wildcard
				parts := strings.SplitN(allowedOrigin, "*.", 2)
				if len(parts) == 2 {
					scheme := parts[0] // "https://" or "http://"
					domain := parts[1]  // "example.com"

					// Check if origin matches the pattern
					if strings.HasPrefix(origin, scheme) {
						originWithoutScheme := strings.TrimPrefix(origin, scheme)
						// Check if it's the exact domain or a subdomain
						if originWithoutScheme == domain || strings.HasSuffix(originWithoutScheme, "."+domain) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// setCORSHeaders sets the appropriate CORS headers on the response
func (c *CORSConfig) setCORSHeaders(rw http.ResponseWriter, r *http.Request) {
	if !c.Enabled {
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	// Check if the origin is allowed
	if c.IsOriginAllowed(origin) {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials", "true")

		// Set additional headers for preflight requests
		if r.Method == http.MethodOptions {
			if len(c.AllowedMethods) > 0 {
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ", "))
			}
			if len(c.AllowedHeaders) > 0 {
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ", "))
			}
			if c.MaxAge > 0 {
				rw.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", int(c.MaxAge.Seconds())))
			}
		}
	}
}

// handlePreflight handles CORS preflight requests
func (c *CORSConfig) handlePreflight(rw http.ResponseWriter, r *http.Request) bool {
	if !c.Enabled || r.Method != http.MethodOptions {
		return false
	}

	origin := r.Header.Get("Origin")
	if origin == "" || !c.IsOriginAllowed(origin) {
		return false
	}

	c.setCORSHeaders(rw, r)
	rw.WriteHeader(http.StatusOK)
	return true
}