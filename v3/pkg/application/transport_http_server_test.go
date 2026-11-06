//go:build server

package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Accepted requests carry an empty JSON object, so they stop at the missing
// object check (422) without needing a message processor. Rejected requests
// must never reach dispatch: with no processor configured they would panic.
func TestServerModeRuntimeRequestPolicy(t *testing.T) {
	tests := []struct {
		name     string
		server   ServerOptions
		method   string
		target   string
		body     string
		headers  map[string]string
		wantCode int
	}{
		{
			name:     "GET with query parameters rejected",
			method:   http.MethodGet,
			target:   "/wails/runtime?object=2&method=2",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "POST with query parameters from cross origin rejected",
			method:   http.MethodPost,
			target:   "/wails/runtime?object=2&method=2",
			headers:  map[string]string{"Origin": "https://untrusted.example", "Content-Type": "text/plain"},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "cross origin POST rejected",
			method:   http.MethodPost,
			body:     `{"object":2,"method":2}`,
			headers:  map[string]string{"Origin": "https://untrusted.example"},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "opaque origin rejected",
			method:   http.MethodPost,
			body:     `{"object":2,"method":2}`,
			headers:  map[string]string{"Origin": "null"},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "cross site fetch metadata rejected",
			method:   http.MethodPost,
			body:     `{"object":2,"method":2}`,
			headers:  map[string]string{"Sec-Fetch-Site": "cross-site"},
			wantCode: http.StatusForbidden,
		},
		{
			name:   "cross origin chunk rejected",
			method: http.MethodPost,
			body:   `{"object":2,"method":2}`,
			headers: map[string]string{
				"Origin":         "https://untrusted.example",
				chunkIDHeader:    "chunk",
				chunkIndexHeader: "0",
				chunkTotalHeader: "1",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "same origin POST accepted",
			method:   http.MethodPost,
			body:     `{}`,
			headers:  map[string]string{"Origin": "http://example.com", "Sec-Fetch-Site": "same-origin"},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "non-browser POST accepted",
			method:   http.MethodPost,
			body:     `{}`,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "configured origin accepted",
			server: ServerOptions{
				WebSocketOriginPatterns: []string{"https://trusted.example"},
			},
			method:   http.MethodPost,
			body:     `{}`,
			headers:  map[string]string{"Origin": "https://trusted.example", "Sec-Fetch-Site": "cross-site"},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "configured pattern does not match other scheme",
			server: ServerOptions{
				WebSocketOriginPatterns: []string{"https://trusted.example"},
			},
			method:   http.MethodPost,
			body:     `{"object":2,"method":2}`,
			headers:  map[string]string{"Origin": "http://trusted.example"},
			wantCode: http.StatusForbidden,
		},
		{
			name: "all origins require explicit opt in",
			server: ServerOptions{
				WebSocketAllowAllOrigins: true,
			},
			method:   http.MethodPost,
			body:     `{}`,
			headers:  map[string]string{"Origin": "https://untrusted.example"},
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newServerTestApp(t)
			a.options.Server = tt.server
			previous := globalApplication
			globalApplication = a
			t.Cleanup(func() { globalApplication = previous })

			target := tt.target
			if target == "" {
				target = "/wails/runtime"
			}
			req := httptest.NewRequest(tt.method, target, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()
			next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Fatal("runtime request fell through to the asset handler")
			})
			NewHTTPTransport().Handler()(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantCode, rec.Body.String())
			}
		})
	}
}
