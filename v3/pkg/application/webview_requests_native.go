//go:build wails_native

package application

import "unsafe"

// The Darwin application delegate retains the C-exported URL callback for ABI
// compatibility, but no native build includes a producer for this channel.
type webViewAssetRequest struct{}

func newWebViewAssetRequest(unsafe.Pointer, uint, string) *webViewAssetRequest {
	return &webViewAssetRequest{}
}

var webviewRequests = make(chan *webViewAssetRequest)
