//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "webview_window_titlebar_accessory_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// addTitlebarAccessoryNow performs the actual native construction; it must
// run on the main thread. Called either immediately from AddTitlebarAccessory
// (window already running) or from installPendingTitlebarAccessories (as
// part of run(), for accessories requested before the window existed).
func addTitlebarAccessoryNow(w *WebviewWindow, options MacTitlebarAccessoryOptions, accessory *MacTitlebarAccessory) error {
	impl, ok := w.impl.(*macosWebviewWindow)
	if !ok {
		return nil // no-op on platforms where impl isn't the darwin type (shouldn't happen here)
	}

	resolvedURL := ""
	if options.URL != "" {
		resolved, err := assetserver.GetStartURL(options.URL)
		if err != nil {
			return err
		}
		resolvedURL = resolved
	}

	urlC := C.CString(resolvedURL)
	defer C.free(unsafe.Pointer(urlC))

	controller := C.titlebarAccessoryAdd(impl.nsWindow, C.int(options.Position),
		C.double(options.Width), C.double(options.Height), urlC)

	accessory.handleLock.Lock()
	accessory.nativeController = controller
	accessory.nativeWebview = C.titlebarAccessoryWebview(controller)
	accessory.handleLock.Unlock()

	return nil
}

func removeTitlebarAccessoryNow(w *WebviewWindow, controller unsafe.Pointer) {
	impl, ok := w.impl.(*macosWebviewWindow)
	if !ok {
		return
	}
	C.titlebarAccessoryRemove(impl.nsWindow, controller)
}

// installPendingTitlebarAccessories runs on the main thread as part of
// run(), after the window (and, if applicable, its split panes) exist.
func (w *macosWebviewWindow) installPendingTitlebarAccessories() {
	w.parent.titlebarAccessoriesLock.Lock()
	pending := w.parent.titlebarAccessoriesPending
	w.parent.titlebarAccessoriesPending = nil
	w.parent.titlebarAccessoriesLock.Unlock()

	for _, p := range pending {
		if err := addTitlebarAccessoryNow(w.parent, p.options, p.accessory); err != nil {
			w.parent.Error("AddTitlebarAccessory %q: %s", p.accessory.name, err)
		}
	}
}

func macTitlebarAccessoryWebviewSetURL(webview unsafe.Pointer, url string) {
	resolvedURL := url
	if resolved, err := assetserver.GetStartURL(url); err == nil {
		resolvedURL = resolved
	} else {
		globalApplication.error("MacTitlebarAccessory.SetURL: resolving URL %q: %s", url, err)
	}
	urlC := C.CString(resolvedURL)
	defer C.free(unsafe.Pointer(urlC))
	C.titlebarAccessoryWebviewSetURL(webview, urlC)
}

func macTitlebarAccessoryWebviewExecJS(webview unsafe.Pointer, js string) {
	jsC := C.CString(js)
	defer C.free(unsafe.Pointer(jsC))
	C.titlebarAccessoryWebviewExecJS(webview, jsC)
}
