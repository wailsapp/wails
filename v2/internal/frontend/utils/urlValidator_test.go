package utils_test

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/internal/frontend/utils"
)

// Test cases for ValidateAndOpenURL
func TestValidateURL(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		shouldErr bool
		errMsg    string
		expected  string
	}{
		// Valid URLs
		{
			name:      "valid https URL",
			url:       "https://www.example.com",
			shouldErr: false,
			expected:  "https://www.example.com",
		},
		{
			name:      "valid http URL",
			url:       "http://example.com",
			shouldErr: false,
			expected:  "http://example.com",
		},
		{
			name:      "URL with query parameters",
			url:       "https://example.com/search?q=cats&dogs",
			shouldErr: false,
			expected:  "https://example.com/search?q=cats&dogs",
		},
		{
			name:      "URL with port",
			url:       "https://example.com:8080/path",
			shouldErr: false,
			expected:  "https://example.com:8080/path",
		},
		{
			name:      "URL with fragment",
			url:       "https://example.com/page#section",
			shouldErr: false,
			expected:  "https://example.com/page#section",
		},
		{
			name:      "urlencode params",
			url:       "http://google.com/ ----browser-subprocess-path=C:\\\\Users\\\\Public\\\\test.bat",
			shouldErr: false,
			expected:  "http://google.com/%20----browser-subprocess-path=C:%5C%5CUsers%5C%5CPublic%5C%5Ctest.bat",
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
			name:      "URL with tab character",
			url:       "https://example.com/path?q=hello\tworld",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "URL with path parameters",
			url:       "https://example.com/path;param=value",
			shouldErr: true,
			errMsg:    "shell metacharacters not allowed",
		},
		{
			name:      "URL with special characters in query",
			url:       "https://example.com/search?q=hello world&filter=price>100",
			shouldErr: true,
			errMsg:    "shell metacharacters not allowed",
		},
		{
			name:      "URL with special characters in query and params",
			url:       "https://example.com/search?q=hello ----browser-subprocess-path=C:\\\\Users\\\\Public\\\\test.bat",
			shouldErr: true,
			errMsg:    "shell metacharacters not allowed",
		},
		{
			name:      "URL with dollar sign in query",
			url:       "https://example.com/search?price=$100",
			shouldErr: true,
			errMsg:    "shell metacharacters not allowed",
		},
		{
			name:      "URL with parentheses",
			url:       "https://example.com/file(1).html",
			shouldErr: true,
			errMsg:    "shell metacharacters not allowed",
		},
		{
			name:      "URL with unicode",
			url:       "https://example.com/search?q=hello\u2001foo",
			shouldErr: true,
			errMsg:    "unicode dangerous characters not allowed",
		},

		// Edge cases
		{
			name:      "international domain",
			url:       "https://例え.テスト/path",
			shouldErr: false,
			expected:  "https://%E4%BE%8B%E3%81%88.%E3%83%86%E3%82%B9%E3%83%88/path",
		},
		{
			name:      "URL with pipe character",
			url:       "https://example.com/user/123|admin",
			shouldErr: false,
			expected:  "https://example.com/user/123%7Cadmin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We'll test only the validation part to avoid actually opening URLs
			sanitized, err := utils.ValidateAndSanitizeURL(tc.url)

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
				if sanitized != tc.expected {
					t.Errorf("unexpected sanitized URL for %q: expected %q, got %q", tc.url, tc.expected, sanitized)
				}
			}
		})
	}
}
