package assetserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceRoutingPreservesEscapedPathAndRequest(t *testing.T) {
	server, err := NewAssetServer(&Options{Handler: http.NotFoundHandler()})
	if err != nil {
		t.Fatal(err)
	}
	hits := 0
	server.AttachServiceHandler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Path != "/item/name" || r.URL.RawPath != "/item%2Fname" || r.URL.EscapedPath() != "/item%2Fname" {
			t.Errorf("service URL path=%q, raw=%q, escaped=%q", r.URL.Path, r.URL.RawPath, r.URL.EscapedPath())
		}
		if r.RequestURI != "/api/item%2Fname?key=value" || r.URL.RawQuery != "key=value" {
			t.Error("service routing changed the request URI or query")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/item%2Fname?key=value", nil)
	original := *request.URL
	for range 2 {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		if response.Code != http.StatusAccepted {
			t.Errorf("status = %d, want 202", response.Code)
		}
		if *request.URL != original {
			t.Error("service routing modified the original request URL")
		}
	}
	if hits != 2 {
		t.Errorf("service handler called %d times, want 2", hits)
	}
}
