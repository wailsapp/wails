package assetserver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
)

func TestBundledRuntimeHEAD(t *testing.T) {
	server := &BundledAssetServer{}
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest(http.MethodHead, "/wails/runtime.js", nil))
	if response.Code != http.StatusOK || response.Body.Len() != 0 {
		t.Errorf("HEAD status = %d, body bytes = %d; want 200, 0", response.Code, response.Body.Len())
	}
	if got := response.Header().Get(HeaderContentLength); got != strconv.Itoa(len(bundledassets.RuntimeJS)) {
		t.Errorf("HEAD Content-Length = %q, want %d", got, len(bundledassets.RuntimeJS))
	}
	if got := response.Header().Get(HeaderContentType); got != "application/javascript" {
		t.Errorf("HEAD Content-Type = %q", got)
	}
}

func TestBundledAssetServerRoutes(t *testing.T) {
	server := &BundledAssetServer{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("application:" + r.URL.Path))
	})}
	for _, tt := range []struct {
		path   string
		status int
		body   string
	}{
		{"/wails/runtime.js", http.StatusOK, string(bundledassets.RuntimeJS)},
		{"/wails/missing.js", http.StatusNotFound, "404 page not found\n"},
		{"/wails/", http.StatusNotFound, "404 page not found\n"},
		{"/index.html", http.StatusOK, "application:/index.html"},
		{"/wailsish/file.js", http.StatusOK, "application:/wailsish/file.js"},
	} {
		t.Run(tt.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			if response.Code != tt.status || response.Body.String() != tt.body {
				t.Errorf("response status = %d, body length = %d; want %d, %d", response.Code, response.Body.Len(), tt.status, len(tt.body))
			}
			if request.URL.Path != tt.path {
				t.Errorf("request path changed to %q, want %q", request.URL.Path, tt.path)
			}
			if tt.path == "/wails/runtime.js" && response.Header().Get(HeaderContentType) != "application/javascript" {
				t.Error("runtime response lost its JavaScript content type")
			}
			if tt.path == "/wails/runtime.js" && response.Header().Get(HeaderContentLength) != strconv.Itoa(len(bundledassets.RuntimeJS)) {
				t.Error("GET runtime response has an incorrect Content-Length")
			}
		})
	}
}
