package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The window ID header is converted to uint32 for the runtime request. Values
// that do not fit (negative or above MaxUint32) must be rejected rather than
// silently wrapped onto another window's ID.
func TestHTTPTransportRuntimeRejectsOutOfRangeWindowID(t *testing.T) {
	for _, id := range []string{"-1", "4294967296", "18446744073709551616", "abc"} {
		t.Run(id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/wails/runtime", strings.NewReader(`{"object":0,"method":0}`))
			req.Header.Set(webViewRequestHeaderWindowId, id)
			rec := httptest.NewRecorder()
			NewHTTPTransport().handleRuntimeRequest(rec, req)
			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422: %s", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), "windowId") {
				t.Fatalf("body = %q, want windowId decode error", rec.Body.String())
			}
		})
	}
}
