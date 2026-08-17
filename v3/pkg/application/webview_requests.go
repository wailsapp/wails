//go:build !wails_native

package application

import (
	"io"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type webViewAssetRequest struct {
	Request    webview.Request
	windowId   uint
	windowName string
}

var _ webview.Request = &webViewAssetRequest{}

func newWebViewAssetRequest(task unsafe.Pointer, windowID uint, windowName string) *webViewAssetRequest {
	return &webViewAssetRequest{Request: webview.NewRequest(task), windowId: windowID, windowName: windowName}
}

func (r *webViewAssetRequest) URL() (string, error)             { return r.Request.URL() }
func (r *webViewAssetRequest) Method() (string, error)          { return r.Request.Method() }
func (r *webViewAssetRequest) Body() (io.ReadCloser, error)     { return r.Request.Body() }
func (r *webViewAssetRequest) Response() webview.ResponseWriter { return r.Request.Response() }
func (r *webViewAssetRequest) Close() error                     { return r.Request.Close() }

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
