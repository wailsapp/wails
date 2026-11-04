package github

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGitHubRequestDeadlines(t *testing.T) {
	for _, lookup := range []struct {
		name, payload string
		run           func()
	}{
		{"tags", "[]", func() { _, _ = GetVersionTags() }},
		{"notes", `{"body":null}`, func() { _, _ = GetReleaseNotes("v3.0.0") }},
	} {
		for _, tc := range []struct {
			name              string
			configured, limit time.Duration
		}{
			{"default", 0, 15 * time.Second},
			{"negative", -time.Second, 15 * time.Second},
			{"longer", time.Minute, 15 * time.Second},
			{"shorter", 2 * time.Second, 2 * time.Second},
		} {
			t.Run(lookup.name+"/"+tc.name, func(t *testing.T) {
				previous := http.DefaultClient
				t.Cleanup(func() { http.DefaultClient = previous })
				called := false
				client := &http.Client{Timeout: tc.configured, Transport: tagTransport(func(r *http.Request) (*http.Response, error) {
					called = true
					deadline, ok := r.Context().Deadline()
					if !ok {
						t.Error("request has no deadline")
					} else if remaining := time.Until(deadline); remaining <= 0 || remaining > tc.limit {
						t.Errorf("request deadline is %s away, want at most %s", remaining, tc.limit)
					}
					return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(lookup.payload))}, nil
				})}
				http.DefaultClient = client
				lookup.run()
				if !called {
					t.Fatal("configured transport was not used")
				}
				if client.Timeout != tc.configured {
					t.Errorf("default client timeout changed to %s", client.Timeout)
				}
			})
		}
	}
}
