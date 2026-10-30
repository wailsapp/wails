package assetserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

type bodyWebViewRequest struct {
	*contextWebViewRequest
	body   io.ReadCloser
	header http.Header
}

func (*bodyWebViewRequest) Method() (string, error)        { return http.MethodPost, nil }
func (r *bodyWebViewRequest) Body() (io.ReadCloser, error) { return r.body, nil }
func (r *bodyWebViewRequest) Header() (http.Header, error) { return r.header, nil }

func TestWebViewProxyPreservesRequestBodies(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("X-Observed-Length", strconv.FormatInt(r.ContentLength, 10))
		_, _ = w.Write(data)
	}))
	defer backend.Close()
	target, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	server, err := NewAssetServer(&Options{Handler: httputil.NewSingleHostReverseProxy(target)})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name       string
		body       string
		length     string
		wantLength int64
	}{
		{"unknown", "payload", "", -1},
		{"declared", "payload", "7", 7},
		{"empty", "", "", 0},
		{"declared empty", "", "0", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			response := &contextWebViewResponse{}
			request := &bodyWebViewRequest{
				contextWebViewRequest: &contextWebViewRequest{ctx: context.Background(), rw: response},
				header:                make(http.Header),
			}
			if tt.body != "" {
				request.body = io.NopCloser(strings.NewReader(tt.body))
			}
			if tt.length != "" {
				request.header.Set(HeaderContentLength, tt.length)
			}
			server.processWebViewRequest(request)
			if response.Code() != http.StatusOK || response.body.String() != tt.body {
				t.Errorf("response status = %d, body = %q; want 200, %q", response.Code(), response.body.String(), tt.body)
			}
			if got := response.Header().Get("X-Observed-Length"); got != strconv.FormatInt(tt.wantLength, 10) {
				t.Errorf("backend length = %q, want %d", got, tt.wantLength)
			}
		})
	}
}
