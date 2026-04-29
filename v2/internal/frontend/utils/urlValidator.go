package utils

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ValidateAndSanitizeURL(rawURL string) (string, error) {
	// Check for null bytes (can cause truncation issues in some systems)
	if strings.Contains(rawURL, "\x00") {
		return "", errors.New("null bytes not allowed in URL")
	}

	// Parse URL first - this handles most malformed URLs
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %v", err)
	}

	scheme := strings.ToLower(parsedURL.Scheme)

	if scheme == "javascript" || scheme == "data" || scheme == "file" || scheme == "ftp" || scheme == "" {
		return "", errors.New("scheme not allowed")
	}

	// Ensure there's actually a host for http/https URLs
	if (scheme == "http" || scheme == "https") && parsedURL.Host == "" {
		return "", fmt.Errorf("missing host for %s URL", scheme)
	}

	sanitizedURL := parsedURL.String()

	// Check for control characters that might cause issues
	// (but allow legitimate URL characters like &, ;, etc.)
	for i, r := range sanitizedURL {
		// Block control characters except tab, but allow other printable chars
		if r < 32 && r != 9 { // 9 is tab, which might be legitimate
			return "", fmt.Errorf("control character at position %d not allowed", i)
		}
	}

	// Shell metacharacter check
	shellDangerous := `[;\|` + "`" + `$\\<>*{}\[\]()~! \t\n\r]`
	if matched, _ := regexp.MatchString(shellDangerous, sanitizedURL); matched {
		return "", errors.New("shell metacharacters not allowed")
	}

	// Unicode danger check
	unicodeDangerous := "[\u0000-\u001F\u007F\u00A0\u1680\u2000-\u200F\u2028-\u202F\u205F\u2060\u3000\uFEFF]"
	if matched, _ := regexp.MatchString(unicodeDangerous, sanitizedURL); matched {
		return "", errors.New("unicode dangerous characters not allowed")
	}

	return sanitizedURL, nil
}
