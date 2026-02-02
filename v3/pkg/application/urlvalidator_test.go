package application_test

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestValidateURL(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		shouldErr bool
		errMsg    string
		expected  string
	}{
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
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "javascript scheme",
			url:       "javascript:alert('XSS')",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "data scheme",
			url:       "data:text/html,<script>alert('XSS')</script>",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "file scheme",
			url:       "file:///etc/passwd",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "ftp scheme",
			url:       "ftp://ftp.example.com/file",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "missing scheme",
			url:       "example.com",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "empty string",
			url:       "",
			shouldErr: true,
			errMsg:    "scheme not allowed",
		},
		{
			name:      "null byte in URL",
			url:       "https://example.com\x00/malicious",
			shouldErr: true,
			errMsg:    "null bytes not allowed",
		},
		{
			name:      "control character",
			url:       "https://example.com\x01",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "shell injection with semicolon",
			url:       "https://example.com/;rm -rf /",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "shell injection with pipe",
			url:       "https://example.com/|cat /etc/passwd",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "shell injection with backtick",
			url:       "https://example.com/`whoami`",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "shell injection with dollar",
			url:       "https://example.com/$(whoami)",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "unicode null",
			url:       "https://example.com/\u0000",
			shouldErr: true,
			errMsg:    "null bytes not allowed",
		},
		{
			name:      "missing host for http",
			url:       "http:///path",
			shouldErr: true,
			errMsg:    "missing host",
		},
		{
			name:      "missing host for https",
			url:       "https:///path",
			shouldErr: true,
			errMsg:    "missing host",
		},
		{
			name:      "URL with newline",
			url:       "https://example.com/path\n/newline",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "URL with carriage return",
			url:       "https://example.com/path\r/return",
			shouldErr: true,
			errMsg:    "control character",
		},
		{
			name:      "URL with tab",
			url:       "https://example.com/path\t/tab",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with space in path",
			url:       "https://example.com/path with spaces",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with angle brackets",
			url:       "https://example.com/<script>",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with parentheses",
			url:       "https://example.com/(test)",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with curly braces",
			url:       "https://example.com/{test}",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with square brackets",
			url:       "https://example.com/[test]",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with asterisk",
			url:       "https://example.com/*",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with tilde",
			url:       "https://example.com/~user",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "URL with exclamation",
			url:       "https://example.com/!test",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
		{
			name:      "zero-width space",
			url:       "https://example.com/\u200B",
			shouldErr: true,
			errMsg:    "dangerous unicode",
		},
		{
			name:      "zero-width non-joiner",
			url:       "https://example.com/\u200C",
			shouldErr: true,
			errMsg:    "dangerous unicode",
		},
		{
			name:      "zero-width joiner",
			url:       "https://example.com/\u200D",
			shouldErr: true,
			errMsg:    "dangerous unicode",
		},
		{
			name:      "right-to-left override",
			url:       "https://example.com/\u202E",
			shouldErr: true,
			errMsg:    "dangerous unicode",
		},
		{
			name:      "invalid URL format",
			url:       "ht!tp://[invalid",
			shouldErr: true,
			errMsg:    "shell metacharacters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := application.ValidateAndSanitizeURL(tc.url)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tc.errMsg)
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected '%s', got '%s'", tc.expected, result)
				}
			}
		})
	}
}