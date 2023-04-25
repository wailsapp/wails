package assetserver

import (
	"io"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
)

// ProcessHTTPRequest processes the HTTP Request by faking a golang HTTP Server.
// The request will be finished with a StatusNotImplemented code if no handler has written to the response.
func (d *AssetServer) ProcessHTTPRequestLegacy(rw http.ResponseWriter, reqGetter func() (*http.Request, error)) {
	d.processWebViewRequest(&legacyRequest{reqGetter: reqGetter, rw: rw})
}

type legacyRequest struct {
	req *http.Request
	rw  http.ResponseWriter

	reqGetter func() (*http.Request, error)
}

func (r *legacyRequest) URL() (string, error) {
	req, err := r.request()
	if err != nil {
		return "", err
	}
	return req.URL.String(), nil
}

func (r *legacyRequest) Method() (string, error) {
	req, err := r.request()
	if err != nil {
		return "", err
	}
	return req.Method, nil
}

func (r *legacyRequest) Header() (http.Header, error) {
	req, err := r.request()
	if err != nil {
		return nil, err
	}
	return req.Header, nil
}

func (r *legacyRequest) Body() (io.ReadCloser, error) {
	req, err := r.request()
	if err != nil {
		return nil, err
	}
	return req.Body, nil
}

func (r legacyRequest) Response() webview.ResponseWriter {
	return &legacyRequestNoOpCloserResponseWriter{r.rw}
}

func (r legacyRequest) Close() error { return nil }

func (r *legacyRequest) request() (*http.Request, error) {
	if r.req != nil {
		return r.req, nil
	}

	req, err := r.reqGetter()
	if err != nil {
		return nil, err
	}
	r.req = req
	return req, nil
}

type legacyRequestNoOpCloserResponseWriter struct {
	http.ResponseWriter
}

func (*legacyRequestNoOpCloserResponseWriter) Finish() {}
