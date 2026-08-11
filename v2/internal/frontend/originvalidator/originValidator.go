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

// wildcardPatternToRegex converts a wildcard pattern to an anchored regex.
//
// A '*' is treated as a wildcard only when it spans a COMPLETE origin
// component — bounded on the left by "://", "." or ":" AND on the right by
// ".", ":", "/" or the end of the pattern. Such a wildcard expands to a single
// non-empty component matcher ([^.:/@]+) that cannot cross a scheme/host/port
// or userinfo boundary. Any other '*' (a partial-label wildcard such as "myapp.com*",
// "*myapp.com" or "myapp.*com") is treated as a literal character, so it
// cannot widen the match — a misused trailing/partial wildcard simply fails
// closed instead of allowing suffix or cross-boundary bypasses
// (GHSA-47hv-j4px-h3c9). \A and \z anchor the whole string with no
// trailing-newline leniency (unlike ^...$).
func (v *OriginValidator) wildcardPatternToRegex(wildcardPattern string) string {
	runes := []rune(wildcardPattern)
	var b strings.Builder
	b.WriteString(`\A`)
	for i, c := range runes {
		if c == '*' && isComponentWildcard(runes, i) {
			// One non-empty component. '@' is excluded alongside the . : /
			// separators so a wildcard component can never absorb a userinfo
			// delimiter (e.g. "myapp.com:*" must not match
			// "myapp.com:x@evilcom", whose real host is "evilcom").
			b.WriteString(`[^.:/@]+`)
			continue
		}
		b.WriteString(regexp.QuoteMeta(string(c)))
	}
	b.WriteString(`\z`)
	return b.String()
}

// isComponentWildcard reports whether the '*' at index i spans a complete
// origin component: a separator boundary ("://", "." or ":") immediately to its
// left, and a separator boundary (".", ":", "/" or end of pattern) immediately
// to its right. Only such wildcards are safe to expand; any other '*' would let
// the wildcard bleed across a host/port boundary and is left literal.
func isComponentWildcard(runes []rune, i int) bool {
	// Left boundary.
	leftOK := false
	switch {
	case i == 0:
		leftOK = false
	case runes[i-1] == '.' || runes[i-1] == ':':
		leftOK = true
	case runes[i-1] == '/' && i >= 3 && string(runes[i-3:i]) == "://":
		leftOK = true
	}
	if !leftOK {
		return false
	}
	// Right boundary.
	if i == len(runes)-1 {
		return true
	}
	switch runes[i+1] {
	case '.', ':', '/':
		return true
	}
	return false
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
