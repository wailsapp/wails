//go:build darwin

package webcontentsview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import "webcontentsview_darwin.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type macosWebContentsView struct {
	parent *WebContentsView
	nsView unsafe.Pointer
	nsWindow unsafe.Pointer
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	var cUserAgent *C.char
	if parent.options.WebPreferences.UserAgent != "" {
		cUserAgent = C.CString(parent.options.WebPreferences.UserAgent)
		defer C.free(unsafe.Pointer(cUserAgent))
	}

	prefs := C.WebContentsViewPreferences{
		devTools:                 C.bool(parent.options.WebPreferences.DevTools != application.Disabled),
		javascript:               C.bool(parent.options.WebPreferences.Javascript != application.Disabled),
		webSecurity:              C.bool(parent.options.WebPreferences.WebSecurity != application.Disabled),
		images:                   C.bool(parent.options.WebPreferences.Images != application.Disabled),
		plugins:                  C.bool(parent.options.WebPreferences.Plugins == application.Enabled),
		zoomFactor:               C.double(parent.options.WebPreferences.ZoomFactor),
		defaultFontSize:          C.int(parent.options.WebPreferences.DefaultFontSize),
		defaultMonospaceFontSize: C.int(parent.options.WebPreferences.DefaultMonospaceFontSize),
		minimumFontSize:          C.int(parent.options.WebPreferences.MinimumFontSize),
		userAgent:                cUserAgent,
	}

	if prefs.zoomFactor == 0 {
		prefs.zoomFactor = 1.0
	}

	var view = C.createWebContentsView(
		C.int(parent.options.Bounds.X), 
		C.int(parent.options.Bounds.Y), 
		C.int(parent.options.Bounds.Width), 
		C.int(parent.options.Bounds.Height),
		prefs,
	)
	
	result := &macosWebContentsView{
		parent: parent,
		nsView: view,
	}
	
	if parent.options.URL != "" {
		result.setURL(parent.options.URL)
	}
	
	return result
}

func (w *macosWebContentsView) setBounds(bounds application.Rect) {
	C.webContentsViewSetBounds(w.nsView, C.int(bounds.X), C.int(bounds.Y), C.int(bounds.Width), C.int(bounds.Height))
}

func (w *macosWebContentsView) setURL(url string) {
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))
	C.webContentsViewSetURL(w.nsView, cUrl)
}

func (w *macosWebContentsView) goBack() {
	C.webContentsViewGoBack(w.nsView)
}

func (w *macosWebContentsView) getURL() string {
	cUrl := C.webContentsViewGetURL(w.nsView)
	if cUrl == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cUrl))
	return C.GoString(cUrl)
}

func (w *macosWebContentsView) execJS(js string) {
	cJs := C.CString(js)
	defer C.free(unsafe.Pointer(cJs))
	C.webContentsViewExecJS(w.nsView, cJs)
}

func (w *macosWebContentsView) attach(window application.Window) {
	w.nsWindow = window.NativeWindow()
	if w.nsWindow != nil {
		C.windowAddWebContentsView(w.nsWindow, w.nsView)
	}
}

func (w *macosWebContentsView) detach() {
	if w.nsWindow != nil {
		C.windowRemoveWebContentsView(w.nsWindow, w.nsView)
		w.nsWindow = nil
	}
}

func (w *macosWebContentsView) nativeView() unsafe.Pointer {
	return w.nsView
}
