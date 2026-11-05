package assetserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	_ "unsafe"

	"github.com/stretchr/testify/assert"
)

func TestContentSniffing(t *testing.T) {
	longLead := strings.Repeat(" ", 512-6)

	tests := map[string]struct {
		Expect string
		Status int
		Header map[string][]string
		Body   []string
	}{
		"/simple": {
			Expect: "text/html; charset=utf-8",
			Body:   []string{"<html><body>Hello!</body></html>"},
		},
		"/split": {
			Expect: "text/html; charset=utf-8",
			Body: []string{
				"<html><body>Hello!",
				"</body></html>",
			},
		},
		"/lead/short/simple": {
			Expect: "text/html; charset=utf-8",
			Body: []string{
				"                                " + "<html><body>Hello!</body></html>",
			},
		},
		"/lead/short/split": {
			Expect: "text/html; charset=utf-8",
			Body: []string{
				"                                ",
				"<html><body>Hello!</body></html>",
			},
		},
		"/lead/long/simple": {
			Expect: "text/html; charset=utf-8",
			Body: []string{
				longLead + "<html><body>Hello!</body></html>",
			},
		},
		"/lead/long/split": {
			Expect: "text/html; charset=utf-8",
			Body: []string{
				longLead,
				"<html><body>Hello!</body></html>",
			},
		},
		"/lead/toolong/simple": {
			Expect: "text/plain; charset=utf-8",
			Body: []string{
				"Hello" + longLead + "<html><body>Hello!</body></html>",
			},
		},
		"/lead/toolong/split": {
			Expect: "text/plain; charset=utf-8",
			Body: []string{
				"Hello" + longLead,
				"<html><body>Hello!</body></html>",
			},
		},
		"/header": {
			Expect: "text/html; charset=utf-8",
			Status: http.StatusForbidden,
			Header: map[string][]string{
				"X-Custom": {"CustomValue"},
			},
			Body: []string{"<html><body>Hello!</body></html>"},
		},
		"/custom": {
			Expect: "text/plain;charset=utf-8",
			Header: map[string][]string{
				"Content-Type": {"text/plain;charset=utf-8"},
			},
			Body: []string{"<html><body>Hello!</body></html>"},
		},
	}

	srv, err := NewAssetServer(&Options{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			test, ok := tests[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			for key, values := range test.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			if test.Status != 0 {
				w.WriteHeader(test.Status)
			}

			for _, chunk := range test.Body {
				w.Write([]byte(chunk))
			}
		}),
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatal("AssetServer initialisation failed: ", err)
	}

	for path, test := range tests {
		t.Run(path[1:], func(t *testing.T) {
			res := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Fatal("http.NewRequest failed: ", err)
			}

			srv.ServeHTTP(res, req)

			expectedStatus := http.StatusOK
			if test.Status != 0 {
				expectedStatus = test.Status
			}
			if res.Code != expectedStatus {
				t.Errorf("Status code mismatch: want %d, got %d", expectedStatus, res.Code)
			}

			if ct := res.Header().Get("Content-Type"); ct != test.Expect {
				t.Errorf("Content type mismatch: want '%s', got '%s'", test.Expect, ct)
			}

			for key, values := range test.Header {
				assert.Equal(t, values, res.Header().Values(key), "Header '%s' mismatch", key)
			}

			assert.Equal(t, strings.Join(test.Body, ""), res.Body.String(), "Response body mismatch")
		})
	}
}

func TestIndexFallback(t *testing.T) {
	// Paths to try and whether a 404 should trigger a fallback.
	paths := map[string]bool{
		"":            true,
		"/":           true,
		"/index":      false,
		"/index.html": true,
		"/other":      false,
	}

	statuses := []int{
		http.StatusOK,
		http.StatusNotFound,
		http.StatusForbidden,
	}

	header := map[string][]string{
		"X-Custom": {"CustomValue"},
	}
	body := "<html><body>Hello!</body></html>"

	srv, err := NewAssetServer(&Options{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, values := range header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			status, err := strconv.Atoi(r.URL.Query().Get("status"))
			if err == nil && status != 0 && status != http.StatusOK {
				w.WriteHeader(status)
			}

			w.Write([]byte(body))
		}),
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatal("AssetServer initialisation failed: ", err)
	}

	for path, fallback := range paths {
		for _, status := range statuses {
			key := "<empty path>"
			if len(path) > 0 {
				key = path[1:]
			}

			t.Run(fmt.Sprintf("%s/status=%d", key, status), func(t *testing.T) {
				res := httptest.NewRecorder()

				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?status=%d", path, status), nil)
				if err != nil {
					t.Fatal("http.NewRequest failed: ", err)
				}

				srv.ServeHTTP(res, req)

				fallbackTriggered := false
				if status == http.StatusNotFound && fallback {
					status = http.StatusOK
					fallbackTriggered = true
				}

				if res.Code != status {
					t.Errorf("Status code mismatch: want %d, got %d", status, res.Code)
				}

				if fallbackTriggered {
					if body == res.Body.String() {
						t.Errorf("Fallback response has the same body as not found response")
					}
					return
				} else {
					for key, values := range header {
						assert.Equal(t, values, res.Header().Values(key), "Header '%s' mismatch", key)
					}

					assert.Equal(t, body, res.Body.String(), "Response body mismatch")
				}
			})
		}
	}
}
