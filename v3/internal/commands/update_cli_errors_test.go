package commands

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/debug"
)

type updateDiscoveryTransport func(*http.Request) (*http.Response, error)

func (f updateDiscoveryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestUpdateCLIDiscoveryFailures(t *testing.T) {
	networkErr := errors.New("network unavailable")
	cases := []struct {
		name, body, want string
		status           int
		err              error
	}{
		{"network", "", "network unavailable", 0, networkErr},
		{"rate_limit", `{ "message": "rate limited" }`, "HTTP 403", http.StatusForbidden, nil},
		{"invalid_json", "{", "unexpected end of JSON input", http.StatusOK, nil},
		{"no_stable", `[{"name":"v3.0.0-beta.1"}]`, "", http.StatusOK, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			previousClient, previousModulePath := http.DefaultClient, debug.LocalModulePath
			t.Cleanup(func() {
				http.DefaultClient, debug.LocalModulePath = previousClient, previousModulePath
			})
			debug.LocalModulePath = ""
			calls := 0
			http.DefaultClient = &http.Client{Transport: updateDiscoveryTransport(func(*http.Request) (*http.Response, error) {
				calls++
				if tc.err != nil {
					return nil, tc.err
				}
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			err := UpdateCLI(&UpdateCLIOptions{})
			if tc.want == "" {
				if err != nil {
					t.Errorf("missing stable release returned error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error=%v, want %q", err, tc.want)
			}
			if tc.err != nil && !errors.Is(err, tc.err) {
				t.Errorf("network error identity lost: %v", err)
			}
			if calls != 1 {
				t.Errorf("made %d requests, want 1", calls)
			}
		})
	}
}
