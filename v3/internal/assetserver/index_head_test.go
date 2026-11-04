package assetserver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFallbackIndexHEAD(t *testing.T) {
	server, err := NewAssetServer(&Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/index.html"} {
		for _, language := range []string{"", "en", "zh", "unknown"} {
			t.Run(path+"/"+language, func(t *testing.T) {
				get := httptest.NewRecorder()
				getRequest := httptest.NewRequest(http.MethodGet, path, nil)
				getRequest.Header.Set(HeaderAcceptLanguage, language)
				server.ServeHTTP(get, getRequest)
				head := httptest.NewRecorder()
				headRequest := getRequest.Clone(getRequest.Context())
				headRequest.Method = http.MethodHead
				server.ServeHTTP(head, headRequest)
				if get.Code != http.StatusOK || get.Body.Len() == 0 {
					t.Fatalf("GET status = %d, body length = %d", get.Code, get.Body.Len())
				}
				if head.Code != get.Code || head.Body.Len() != 0 {
					t.Errorf("HEAD status = %d, body length = %d; want %d, 0", head.Code, head.Body.Len(), get.Code)
				}
				if got := head.Header().Get(HeaderContentLength); got != strconv.Itoa(get.Body.Len()) {
					t.Errorf("HEAD Content-Length = %q, want %d", got, get.Body.Len())
				}
				if got := head.Header().Get(HeaderContentType); got != get.Header().Get(HeaderContentType) {
					t.Errorf("HEAD Content-Type = %q, want %q", got, get.Header().Get(HeaderContentType))
				}
			})
		}
	}
}
