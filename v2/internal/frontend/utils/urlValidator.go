package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Alternative simpler validation function if you don't need the struct approach
func ValidateURL(rawURL string) error {
	// Check for null bytes (can cause truncation issues in some systems)
	if strings.Contains(rawURL, "\x00") {
		return errors.New("null bytes not allowed in URL")
	}

	// Parse URL first - this handles most malformed URLs
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	scheme := parsedURL.Scheme

	if scheme == "javascript" || scheme == "data" || scheme == "file" || scheme == "ftp" || scheme == "" {
		return errors.New("scheme not allowed")
	}

	// Ensure there's actually a host for http/https URLs
	if (scheme == "http" || scheme == "https") && parsedURL.Host == "" {
		return fmt.Errorf("missing host for %s URL", scheme)
	}

	// Optional: Check for control characters that might cause issues
	// (but allow legitimate URL characters like &, ;, etc.)
	for i, r := range rawURL {
		// Block control characters except tab, but allow other printable chars
		if r < 32 && r != 9 { // 9 is tab, which might be legitimate
			return fmt.Errorf("control character at position %d not allowed", i)
		}
	}

	return nil
}
