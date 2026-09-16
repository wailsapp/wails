//go:build !wails_native

package application

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type webViewAssetRequest struct {
	Request    webview.Request
	windowId   uint
	windowName string
}

var _ webview.Request = &webViewAssetRequest{}

func (r *webViewAssetRequest) URL() (string, error)             { return r.Request.URL() }
func (r *webViewAssetRequest) Method() (string, error)          { return r.Request.Method() }
func (r *webViewAssetRequest) Body() (io.ReadCloser, error)     { return r.Request.Body() }
func (r *webViewAssetRequest) Response() webview.ResponseWriter { return r.Request.Response() }
func (r *webViewAssetRequest) Close() error                     { return r.Request.Close() }

// Context preserves native request cancellation through the header-injecting wrapper.
func (r *webViewAssetRequest) Context() context.Context {
	if contextual, ok := r.Request.(interface{ Context() context.Context }); ok {
		return contextual.Context()
	}
	return nil
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	header, err := r.Request.Header()
	if err != nil {
		return nil, err
	}
	result := header.Clone()
	result.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	if r.windowName != "" {
		result.Set(webViewRequestHeaderWindowName, r.windowName)
	}
	return result, nil
}

var webviewRequests = make(chan *webViewAssetRequest, 256)

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	defer handlePanic()
	url, _ := request.Request.URL()
	a.debug("handleWebViewRequest: Processing request", "url", url)
	a.assets.ServeWebViewRequest(request)
	a.debug("handleWebViewRequest: Request processing complete", "url", url)
}
