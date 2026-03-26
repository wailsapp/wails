package originvalidator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type OriginValidator struct {
	allowedOrigins []string
}

// NewOriginValidator creates a new validator from a comma-separated string of allowed origins
func NewOriginValidator(startUrl *url.URL, allowedOriginsString string) *OriginValidator {
	allowedOrigins := startUrl.Scheme + "://" + startUrl.Host
	if allowedOriginsString != "" {
		allowedOrigins += "," + allowedOriginsString
	}
	validator := &OriginValidator{}
	validator.parseAllowedOrigins(allowedOrigins)
	return validator
}

// parseAllowedOrigins parses the comma-separated origins string
func (v *OriginValidator) parseAllowedOrigins(originsString string) {
	if originsString == "" {
		v.allowedOrigins = []string{}
		return
	}

	origins := strings.Split(originsString, ",")
	var trimmedOrigins []string

	for _, origin := range origins {
		trimmed := strings.TrimSuffix(strings.TrimSpace(origin), "/")
		if trimmed != "" {
			trimmedOrigins = append(trimmedOrigins, trimmed)
		}
	}

	v.allowedOrigins = trimmedOrigins
}

// IsOriginAllowed checks if the given origin is allowed
func (v *OriginValidator) IsOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range v.allowedOrigins {
		if v.matchesOriginPattern(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

// matchesOriginPattern checks if origin matches the pattern (supports wildcards)
func (v *OriginValidator) matchesOriginPattern(pattern, origin string) bool {
	// Exact match
	if pattern == origin {
		return true
	}

	// Wildcard pattern matching
	if strings.Contains(pattern, "*") {
		regexPattern := v.wildcardPatternToRegex(pattern)
		matched, err := regexp.MatchString(regexPattern, origin)
		if err != nil {
			return false
		}
		return matched
	}

	return false
}

// wildcardPatternToRegex converts wildcard pattern to regex.
//
// Wildcards are only expanded when they represent a full domain component,
// i.e. immediately after "://", ".", or ":". A wildcard appended to a
// partial label (e.g. "com*") is stripped so that it cannot match across
// domain boundaries (GHSA-47hv-j4px-h3c9).
func (v *OriginValidator) wildcardPatternToRegex(wildcardPattern string) string {
	// Escape special regex characters except *
	specialChars := []string{"\\", ".", "+", "?", "^", "$", "{", "}", "(", ")", "|", "[", "]"}

	escaped := wildcardPattern
	for _, specialChar := range specialChars {
		escaped = strings.ReplaceAll(escaped, specialChar, "\\"+specialChar)
	}

	// Replace * only when preceded by a separator so it matches a full
	// component. Order matters: handle :// before : to avoid partial overlap.
	// Use {0,} instead of * as the regex quantifier to avoid collision with
	// the literal * cleanup in the next step.
	escaped = strings.ReplaceAll(escaped, "//*", "//[^.:/]{0,}")
	escaped = strings.ReplaceAll(escaped, "\\.*", "\\.[^.:/]{0,}")
	escaped = strings.ReplaceAll(escaped, ":*", ":[^.:/]{0,}")

	// Strip any remaining * that is not preceded by a separator.
	// This prevents suffix-based bypasses like "com*" matching "community".
	escaped = strings.ReplaceAll(escaped, "*", "")

	// Anchor the pattern to match the entire string
	return "^" + escaped + "$"
}

// GetOriginFromURL extracts origin from URL string
func (v *OriginValidator) GetOriginFromURL(urlString string) (string, error) {
	if urlString == "" {
		return "", fmt.Errorf("empty URL")
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("URL missing scheme or host")
	}

	// Build origin (scheme + host)
	origin := parsedURL.Scheme + "://" + parsedURL.Host

	return origin, nil
}
