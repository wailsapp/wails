package github

import (
	"net/http"
	"strings"
	"testing"
)

func TestReleaseNotesResponses(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		payload string
		want    string
	}{
		{"success", 200, `{"body":"Useful release notes"}`, "Useful release notes"},
		{"missing_release", 404, `{"message":"Not Found"}`, "No release notes found"},
		{"rate_limit", 403, `{"message":"API rate limit exceeded"}`, "HTTP 403 (Forbidden)"},
		{"too_many_requests", 429, `{"message":"Too many requests"}`, "HTTP 429 (Too Many Requests)"},
		{"server_error", 500, `<html>Server error</html>`, "HTTP 500 (Internal Server Error)"},
		{"error_with_body", 502, `{"body":"Not actual release notes"}`, "HTTP 502 (Bad Gateway)"},
		{"missing_body", 200, `{}`, "No release notes found"},
		{"null_body", 200, `{"body":null}`, "No release notes found"},
		{"invalid_body", 200, `{"body":123}`, "Unable to retrieve release notes"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			previous := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = previous })
			body := &trackedResponseBody{Reader: strings.NewReader(tc.payload)}
			http.DefaultClient = &http.Client{Transport: tagTransport(func(r *http.Request) (*http.Response, error) {
				if got := r.URL.String(); got != "https://api.github.com/repos/wailsapp/wails/releases/tags/v3.0.0" {
					t.Errorf("unexpected request: %s", got)
				}
				return &http.Response{StatusCode: tc.status, Body: body, Header: make(http.Header)}, nil
			})}
			got := GetReleaseNotes("v3.0.0", true)
			if !strings.Contains(got, tc.want) {
				t.Errorf("got %q, want %q", got, tc.want)
			}
			if body.closeCalls != 1 {
				t.Errorf("response body closed %d times, want 1", body.closeCalls)
			}
		})
	}
}
