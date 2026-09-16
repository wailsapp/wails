//go:build darwin && !ios && !server && !wails_native

package application

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

func newWebViewAssetRequest(task unsafe.Pointer, windowID uint, windowName string) *webViewAssetRequest {
	return &webViewAssetRequest{Request: webview.NewRequest(task), windowId: windowID, windowName: windowName}
}

// cancelWebViewAssetRequest tells the asset server the platform has abandoned
// the scheme task so any in-flight handler can stop early.
func cancelWebViewAssetRequest(task unsafe.Pointer) {
	webview.CancelRequest(task)
}
