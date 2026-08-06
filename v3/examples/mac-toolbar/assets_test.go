package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestDaymarkBundledAssetsServeOnlyPrimaryWebContent(t *testing.T) {
	handler := application.BundledAssetFileServer(assets)
	tests := []struct {
		path string
		want string
	}{
		{path: "/", want: "Open Daymark"},
		{path: "/editor.html", want: "editor:ready"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://wails.localhost"+test.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("GET %s returned %d: %s", test.path, response.Code, response.Body.String())
			}
			if !strings.Contains(response.Body.String(), test.want) {
				t.Fatalf("GET %s did not serve the expected Daymark page", test.path)
			}
		})
	}
}

func TestDaymarkEditorContainsNoSidebarOrInspectorChrome(t *testing.T) {
	handler := application.BundledAssetFileServer(assets)
	request := httptest.NewRequest(http.MethodGet, "http://wails.localhost/editor.html", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	body := response.Body.String()
	for _, forbidden := range []string{
		`<aside`,
		`class="sidebar`,
		`class="inspector`,
		`id="details"`,
		`toolbar:details`,
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("primary WebView contains native-chrome substitute %q", forbidden)
		}
	}
}
