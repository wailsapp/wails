package assetserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFallbackResponseWriterHeaderDeletions(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.Header().Set("X-Remove", "original")
			recorder.Header().Set("X-Keep", "preserved")
			fallbackCalled := false
			writer := &fallbackResponseWriter{
				rw:  recorder,
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				fallback: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fallbackCalled = true
					w.WriteHeader(http.StatusOK)
				}),
			}
			writer.Header().Del("X-Remove")
			writer.Header().Set("X-New", "added")
			writer.WriteHeader(status)

			header := recorder.Result().Header
			wantRemoved, wantNew := "", "added"
			if status == http.StatusNotFound {
				wantRemoved, wantNew = "original", ""
			}
			if header.Get("X-Remove") != wantRemoved || header.Get("X-New") != wantNew || header.Get("X-Keep") != "preserved" {
				t.Errorf("unexpected response headers: %v", header)
			}
			if fallbackCalled != (status == http.StatusNotFound) {
				t.Errorf("fallbackCalled = %v for status %d", fallbackCalled, status)
			}
		})
	}
}
