package setupwizard

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStaticAssetsHaveStableContentType guards the fix for the blank-logo bug:
// http.FileServer derives a type from the host MIME database, which on some hosts
// lacks .svg and falls back to sniffing "text/xml" — a type browsers refuse to
// render inside an <img>. setupRoutes pins the type explicitly, so serving must
// not depend on the host. favicon.svg is served under a stable (unhashed) name.
func TestStaticAssetsHaveStableContentType(t *testing.T) {
	w := New()
	mux := http.NewServeMux()
	w.setupRoutes(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cases := []struct {
		path string
		want string
	}{
		{"/favicon.svg", "image/svg+xml"},
		{"/", "text/html; charset=utf-8"},
	}
	for _, tc := range cases {
		resp, err := http.Get(srv.URL + tc.path)
		if err != nil {
			t.Fatalf("GET %s: %v", tc.path, err)
		}
		resp.Body.Close()
		if got := resp.Header.Get("Content-Type"); got != tc.want {
			t.Errorf("GET %s Content-Type = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestStaticContentType(t *testing.T) {
	cases := map[string]string{
		"/assets/wails-logo-black-text-abc123.svg": "image/svg+xml",
		"/assets/index-abc123.js":                  "text/javascript; charset=utf-8",
		"/assets/index-abc123.css":                 "text/css; charset=utf-8",
		"/index.html":                              "text/html; charset=utf-8",
		"/digital_wales_master.webp":               "image/webp",
		"/unknown.xyz":                             "",
	}
	for path, want := range cases {
		if got := staticContentType(path); got != want {
			t.Errorf("staticContentType(%q) = %q, want %q", path, got, want)
		}
	}
}
