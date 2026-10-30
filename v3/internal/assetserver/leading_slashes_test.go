package assetserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestAssetPathsWithLeadingSlashes(t *testing.T) {
	server := newAssetFileServerFS(fstest.MapFS{
		"index.html": {Data: []byte("index page")},
		"asset.txt":  {Data: []byte("asset content")},
	})
	for _, prefix := range []string{"/", "//", "///"} {
		for _, tt := range []struct {
			path, body string
			status     int
		}{
			{"", "index page", http.StatusOK},
			{"asset.txt", "asset content", http.StatusOK},
			{"missing.txt", "", http.StatusNotFound},
		} {
			t.Run(prefix+tt.path, func(t *testing.T) {
				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, "http://localhost"+prefix+tt.path, nil)
				server.ServeHTTP(response, request)
				if response.Code != tt.status || response.Body.String() != tt.body {
					t.Errorf("status = %d, body = %q; want %d, %q", response.Code, response.Body.String(), tt.status, tt.body)
				}
			})
		}
	}
}
