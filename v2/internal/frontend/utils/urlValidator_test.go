package utils_test

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/internal/frontend/utils"
)

// Test cases for ValidateAndOpenURL
func TestValidateAndOpenURL(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		shouldErr bool
		errMsg    string
	}{
		// Valid URLs
		{
			name:      "valid https URL",
			url:       "https://www.example.com",
			shouldErr: false,
		},
		{
			name:      "valid http URL",
			url:       "http://example.com",
			shouldErr: false,
		},
		{
			name:      "URL with query parameters",
			url:       "https://example.com/search?q=cats&dogs",
			shouldErr: false,
		},
		{
			name:      "URL with path parameters",
			url:       "https://example.com/path;param=value",
			shouldErr: false,
		},
		{
			name:      "URL with special characters in query",
			url:       "https://example.com/search?q=hello world&filter=price>100",
			shouldErr: false,
		},
		{
			name:      "URL with port",
			url:       "https://example.com:8080/path",
			shouldErr: false,
		},
		{
			name:      "URL with fragment",
			url:       "https://example.com/page#section",
			shouldErr: false,
		},

		// Invalid schemes
		{
			name:      "javascript scheme",
			url:       "javascript:alert('xss')",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "data scheme",
			url:       "data:text/html,<script>alert(1)</script>",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "file scheme",
			url:       "file:///etc/passwd",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "ftp scheme",
			url:       "ftp://files.example.com/file.txt",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},

		// Malformed URLs
		{
			name:      "not a URL",
			url:       "not-a-url",
			shouldErr: true,
			errMsg:    "scheme not allowed", // will have empty scheme
		},
		{
			name:      "missing scheme",
			url:       "example.com",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "malformed URL",
			url:       "https://",
			shouldErr: true,
			errMsg:    "missing host",
		},
		{
			name:      "empty host",
			url:       "http:///path",
			shouldErr: true,
			errMsg:    "missing host",
		},

		// Security issues
		{
			name:      "null byte in URL",
			url:       "https://example.com\x00/hidden",
			shouldErr: true,
			errMsg:    "null bytes not allowed",
		},
		{
			name:      "control characters",
			url:       "https://example.com\n/path",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "carriage return",
			url:       "https://example.com\r/path",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "URL with tab character (allowed)",
			url:       "https://example.com/path?q=hello\tworld",
			shouldErr: true,
			errMsg:    "control character",
		},

		// Edge cases
		{
			name:      "international domain",
			url:       "https://例え.テスト/path",
			shouldErr: false,
		},

		// URLs that might look suspicious but are valid
		{
			name:      "URL with dollar sign in query",
			url:       "https://example.com/search?price=$100",
			shouldErr: false,
		},
		{
			name:      "URL with parentheses",
			url:       "https://example.com/file(1).html",
			shouldErr: false,
		},
		{
			name:      "URL with pipe character",
			url:       "https://example.com/user/123|admin",
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We'll test only the validation part to avoid actually opening URLs
			err := utils.ValidateURL(tc.url)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error for URL %q, but got none", tc.url)
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("expected error containing %q, got %q", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for URL %q, but got: %v", tc.url, err)
				}
			}
		})
	}
}
