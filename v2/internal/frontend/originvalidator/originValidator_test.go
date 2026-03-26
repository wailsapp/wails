package originvalidator

import (
	"net/url"
	"testing"
)

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func TestWildcardPatternDoesNotCrossDomainBoundaries(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins string
		origin         string
		expected       bool
	}{
		// Legitimate subdomain wildcard — should still work
		{
			name:           "subdomain wildcard matches subdomain",
			allowedOrigins: "https://*.myapp.com",
			origin:         "https://api.myapp.com",
			expected:       true,
		},
		{
			name:           "subdomain wildcard matches different subdomain",
			allowedOrigins: "https://*.myapp.com",
			origin:         "https://www.myapp.com",
			expected:       true,
		},
		{
			name:           "subdomain wildcard rejects different domain",
			allowedOrigins: "https://*.myapp.com",
			origin:         "https://evil.com",
			expected:       false,
		},
		// Trailing wildcard — the vulnerability vector
		{
			name:           "trailing wildcard rejects different TLD (bypass attempt)",
			allowedOrigins: "https://myapp.com*",
			origin:         "https://myapp.community",
			expected:       false,
		},
		{
			name:           "trailing wildcard rejects attacker subdomain (bypass attempt)",
			allowedOrigins: "https://myapp.com*",
			origin:         "https://myapp.com.attacker.com",
			expected:       false,
		},
		{
			name:           "trailing wildcard rejects arbitrary suffix",
			allowedOrigins: "https://myapp.com*",
			origin:         "https://myapp.comXXXXX",
			expected:       false,
		},
		// Exact match still works
		{
			name:           "exact match",
			allowedOrigins: "https://myapp.com",
			origin:         "https://myapp.com",
			expected:       true,
		},
		{
			name:           "exact match rejects different origin",
			allowedOrigins: "https://myapp.com",
			origin:         "https://evil.com",
			expected:       false,
		},
		// Wildcard does not cross into path or port
		{
			name:           "wildcard does not match across port separator",
			allowedOrigins: "https://localhost*",
			origin:         "https://localhost:8080",
			expected:       false,
		},
		{
			name:           "wildcard does not match across path separator",
			allowedOrigins: "https://myapp*",
			origin:         "https://myapp/evil",
			expected:       false,
		},
		// Empty origin
		{
			name:           "empty origin is rejected",
			allowedOrigins: "https://*.myapp.com",
			origin:         "",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startURL := mustParseURL("https://wails.localhost")
			v := NewOriginValidator(startURL, tt.allowedOrigins)
			got := v.IsOriginAllowed(tt.origin)
			if got != tt.expected {
				t.Errorf("IsOriginAllowed(%q) with pattern %q = %v, want %v",
					tt.origin, tt.allowedOrigins, got, tt.expected)
			}
		})
	}
}
