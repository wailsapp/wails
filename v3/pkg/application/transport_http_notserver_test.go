//go:build !server

package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Rejected requests must never reach dispatch (no message processor is set, so
// they would panic). Accepted ones stop at the missing method check.
func TestHTTPTransportRuntimeRequestOutsideServerMode(t *testing.T) {
	const reachedDispatch = http.StatusUnprocessableEntity
	tests := []struct {
		name     string
		method   string
		target   string
		body     string
		headers  map[string]string
		wantCode int
	}{
		{"img-style GET with query", http.MethodGet, "/wails/runtime?object=2", "", map[string]string{"x-wails-client-id": "c"}, http.StatusMethodNotAllowed},
		{"no-cors text/plain POST", http.MethodPost, "/wails/runtime", `{"object":2,"method":2}`, map[string]string{"Content-Type": "text/plain"}, http.StatusForbidden},
		{"no-cors POST with query", http.MethodPost, "/wails/runtime?object=2&method=2", "", nil, http.StatusForbidden},
		{"foreign origin", http.MethodPost, "/wails/runtime", `{"object":2}`, map[string]string{"x-wails-client-id": "c", "Origin": "https://evil.example"}, http.StatusForbidden},
		{"cross-site fetch", http.MethodPost, "/wails/runtime", `{"object":2}`, map[string]string{"x-wails-client-id": "c", "Sec-Fetch-Site": "cross-site"}, http.StatusForbidden},
		{"runtime call, app origin", http.MethodPost, "/wails/runtime", `{"object":2}`, map[string]string{"x-wails-client-id": "c", "Origin": "wails://localhost"}, reachedDispatch},
		{"runtime call, no origin", http.MethodPost, "/wails/runtime", `{"object":2}`, map[string]string{"x-wails-client-id": "c"}, reachedDispatch},
		{"WebKitGTK 6 query fallback", http.MethodPost, "/wails/runtime?object=2", "", map[string]string{"x-wails-client-id": "c", "Origin": "null"}, reachedDispatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://localhost"+tt.target, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()
			NewHTTPTransport().handleRuntimeRequest(rec, req)
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d (%q), want %d", rec.Code, rec.Body.String(), tt.wantCode)
			}
		})
	}
}
