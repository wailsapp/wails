package application

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ValidateAndSanitizeURL(rawURL string) (string, error) {
	if strings.Contains(rawURL, "\x00") {
		return "", errors.New("null bytes not allowed in URL")
	}

	for i, r := range rawURL {
		if r < 32 && r != 9 {
			return "", fmt.Errorf("control character at position %d not allowed", i)
		}
	}

	shellDangerous := `[;|` + "`" + `$\\<>*{}\[\]()~! \t\n\r]`
	if matched, _ := regexp.MatchString(shellDangerous, rawURL); matched {
		return "", errors.New("shell metacharacters not allowed")
	}

	unicodeDangerous := "[\u0000-\u001F\u007F\u00A0\u1680\u2000-\u200F\u2028-\u202F\u205F\u3000\uFEFF\u200B-\u200D\u2060\u2061\u2062\u2063\u2064\u206A-\u206F\uFFF0-\uFFFF]"
	if matched, _ := regexp.MatchString(unicodeDangerous, rawURL); matched {
		return "", errors.New("dangerous unicode characters not allowed")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %v", err)
	}

	scheme := strings.ToLower(parsedURL.Scheme)

	if scheme == "javascript" || scheme == "data" || scheme == "file" || scheme == "ftp" || scheme == "" {
		return "", errors.New("scheme not allowed")
	}

	if (scheme == "http" || scheme == "https") && parsedURL.Host == "" {
		return "", fmt.Errorf("missing host for %s URL", scheme)
	}

	sanitizedURL := parsedURL.String()
	return sanitizedURL, nil
}
