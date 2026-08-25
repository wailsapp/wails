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
	tests := []originCase{
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

	runOriginCases(t, tests)
}

type originCase struct {
	name           string
	allowedOrigins string
	origin         string
	expected       bool
}

func runOriginCases(t *testing.T, tests []originCase) {
	t.Helper()
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

// TestWildcardComponentBoundaries covers the regression cases requested in the
// review of GHSA-47hv-j4px-h3c9: a '*' must be a complete origin component
// (separator on BOTH sides) to be a wildcard; partial-label wildcards must not
// expand; and wildcards must never cross host/port/path or userinfo boundaries.
func TestWildcardComponentBoundaries(t *testing.T) {
	runOriginCases(t, []originCase{
		// Subdomain wildcard: matches exactly one subdomain label.
		{"subdomain wildcard matches one label", "https://*.myapp.com", "https://api.myapp.com", true},
		{"subdomain wildcard rejects nested suffix bypass", "https://*.myapp.com", "https://api.myapp.com.evil.com", false},
		{"subdomain wildcard rejects userinfo bypass", "https://*.myapp.com", "https://api.myapp.com@evil.com", false},
		{"subdomain wildcard rejects apex", "https://*.myapp.com", "https://myapp.com", false},
		{"subdomain wildcard rejects multi-level", "https://*.myapp.com", "https://a.b.myapp.com", false},
		{"subdomain wildcard rejects lookalike domain", "https://*.myapp.com", "https://api.notmyapp.com", false},

		// Partial-label wildcards must NOT expand (treated literally => fail closed).
		{"prefix partial wildcard rejects lookalike", "https://*myapp.com", "https://evilmyapp.com", false},
		{"prefix partial wildcard rejects apex", "https://*myapp.com", "https://myapp.com", false},
		{"infix partial wildcard rejects fill", "https://myapp.*com", "https://myapp.evilcom", false},
		{"infix partial wildcard rejects single label", "https://myapp.*com", "https://myapp.xcom", false},
		{"trailing partial wildcard rejects suffix", "https://myapp.com*", "https://myapp.community", false},
		{"trailing partial wildcard rejects apex (fails closed)", "https://myapp.com*", "https://myapp.com", false},

		// Port wildcard: complete component after ':'.
		{"port wildcard matches a port", "https://myapp.com:*", "https://myapp.com:8080", true},
		{"port wildcard does not cross into host", "https://myapp.com:*", "https://myapp.com:8080.evil.com", false},
		{"port wildcard does not cross into path", "https://myapp.com:*", "https://myapp.com:8080/evil", false},
		{"port wildcard requires a port", "https://myapp.com:*", "https://myapp.com", false},
		{"port wildcard rejects userinfo bypass", "https://myapp.com:*", "https://myapp.com:x@evilcom", false},

		// Multiple complete-component wildcards.
		{"sub+tld wildcard matches", "https://*.myapp.*", "https://api.myapp.com", true},
		{"sub+tld wildcard rejects nested suffix", "https://*.myapp.*", "https://api.myapp.com.evil.com", false},
	})
}
