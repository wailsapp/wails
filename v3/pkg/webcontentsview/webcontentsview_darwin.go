//go:build darwin && !server

package webcontentsview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import "webcontentsview_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type macosWebContentsView struct {
	parent    *WebContentsView
	nsView    unsafe.Pointer
	nsWindow  unsafe.Pointer
	pendingJS []string
}

const nativeSnapshotCallbackTimeout = 30 * time.Second

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	return &macosWebContentsView{parent: parent}
}

func (w *macosWebContentsView) createView() {
	if w.nsView != nil {
		return
	}
	options := w.parent.optionsSnapshot()
	var cUserAgent *C.char
	if options.WebPreferences.UserAgent != "" {
		cUserAgent = C.CString(options.WebPreferences.UserAgent)
		defer C.free(unsafe.Pointer(cUserAgent))
	}
	prefs := C.WebContentsViewPreferences{
		devTools:                 C.bool(!options.WebPreferences.DevTools.IsSet() || options.WebPreferences.DevTools.Get()),
		javascript:               C.bool(!options.WebPreferences.Javascript.IsSet() || options.WebPreferences.Javascript.Get()),
		webSecurity:              C.bool(!options.WebPreferences.WebSecurity.IsSet() || options.WebPreferences.WebSecurity.Get()),
		images:                   C.bool(!options.WebPreferences.Images.IsSet() || options.WebPreferences.Images.Get()),
		plugins:                  C.bool(options.WebPreferences.Plugins.IsSet() && options.WebPreferences.Plugins.Get()),
		zoomFactor:               C.double(options.WebPreferences.ZoomFactor),
		defaultFontSize:          C.int(options.WebPreferences.DefaultFontSize),
		defaultMonospaceFontSize: C.int(options.WebPreferences.DefaultMonospaceFontSize),
		minimumFontSize:          C.int(options.WebPreferences.MinimumFontSize),
		userAgent:                cUserAgent,
	}
	if prefs.zoomFactor == 0 {
		prefs.zoomFactor = 1.0
	}
	w.nsView = C.createWebContentsView(C.int(options.Bounds.X), C.int(options.Bounds.Y), C.int(options.Bounds.Width), C.int(options.Bounds.Height), prefs)
}

func (w *macosWebContentsView) setBounds(bounds application.Rect) {
	if w.nsView != nil {
		application.InvokeSync(func() {
			C.webContentsViewSetBounds(w.nsView, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
		})
	}
}

func (w *macosWebContentsView) setURL(url string) {
	if w.nsView == nil {
		return
	}
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	application.InvokeSync(func() { C.webContentsViewSetURL(w.nsView, cURL) })
}

func (w *macosWebContentsView) setHTML(html string) {
	if w.nsView == nil {
		return
	}
	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))
	application.InvokeSync(func() { C.webContentsViewSetHTML(w.nsView, cHTML) })
}

func (w *macosWebContentsView) execJS(js string) {
	if w.nsView == nil {
		w.pendingJS = append(w.pendingJS, js)
		return
	}
	cJS := C.CString(js)
	defer C.free(unsafe.Pointer(cJS))
	application.InvokeSync(func() { C.webContentsViewExecJS(w.nsView, cJS) })
}

func (w *macosWebContentsView) goBack() {
	if w.nsView != nil {
		application.InvokeSync(func() { C.webContentsViewGoBack(w.nsView) })
	}
}

func (w *macosWebContentsView) getURL() string {
	if w.nsView == nil {
		return ""
	}
	var cURL *C.char
	application.InvokeSync(func() { cURL = C.webContentsViewGetURL(w.nsView) })
	if cURL == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cURL))
	return C.GoString(cURL)
}

func (w *macosWebContentsView) takeSnapshot() string {
	if w.nsView == nil {
		return ""
	}
	result := make(chan string, 1)
	id := registerSnapshotCallback(result)
	application.InvokeSync(func() { C.webContentsViewTakeSnapshot(w.nsView, C.uintptr_t(id)) })
	select {
	case snapshot := <-result:
		return snapshot
	case <-time.After(nativeSnapshotCallbackTimeout):
		unregisterSnapshotCallback(id)
		return ""
	}
}

func (w *macosWebContentsView) setVisible(visible bool) {
	if w.nsView != nil {
		application.InvokeSync(func() { C.webContentsViewSetVisible(w.nsView, C.bool(visible)) })
	}
}

func (w *macosWebContentsView) attach(window application.Window) {
	nativeWindow := window.NativeWindow()
	if nativeWindow == nil || w.nsWindow == nativeWindow {
		return
	}
	application.InvokeSync(func() {
		if w.nsWindow != nil {
			C.windowRemoveWebContentsView(w.nsWindow, w.nsView)
		}
		w.createView()
		options := w.parent.optionsSnapshot()
		w.nsWindow = nativeWindow
		C.windowAddWebContentsView(w.nsWindow, w.nsView)
		C.webContentsViewSetVisible(w.nsView, C.bool(w.parent.visibleSnapshot()))
		if options.URL != "" {
			w.setURL(options.URL)
		} else if options.HTML != "" {
			w.setHTML(options.HTML)
		}
		for _, js := range w.pendingJS {
			w.execJS(js)
		}
		w.pendingJS = nil
	})
}

func (w *macosWebContentsView) detach() {
	if w.nsWindow != nil {
		application.InvokeSync(func() { C.windowRemoveWebContentsView(w.nsWindow, w.nsView) })
		w.nsWindow = nil
	}
}

func (w *macosWebContentsView) destroy() {
	if w.nsView != nil {
		application.InvokeSync(func() { C.webContentsViewDestroy(w.nsView) })
		w.nsView = nil
		w.nsWindow = nil
	}
}

func (w *macosWebContentsView) nativeView() unsafe.Pointer { return w.nsView }

//export browserViewSnapshotCallback
func browserViewSnapshotCallback(callbackID C.uintptr_t, base64 *C.char) {
	snapshot := ""
	if base64 != nil {
		snapshot = C.GoString(base64)
	}
	dispatchSnapshotResult(uintptr(callbackID), snapshot)
}
