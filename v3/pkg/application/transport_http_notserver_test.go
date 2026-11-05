//go:build !server

package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Webviews reach the runtime through the app's own asset handler, so the
// server-mode origin policy must not apply to them, and the WebKitGTK
// query-string fallback must keep working.
func TestHTTPTransportRuntimeRequestOutsideServerMode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wails/runtime?object=2", strings.NewReader(""))
	req.Header.Set("Origin", "wails://wails")
	rec := httptest.NewRecorder()
	NewHTTPTransport().handleRuntimeRequest(rec, req)

	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "missing method value") {
		t.Fatalf("status = %d, body = %q; want the query fallback to reach method validation", rec.Code, rec.Body.String())
	}
}
