// Package webview implements Go bindings to https://github.com/zserge/webview C library.
// It is a modified version of webview.go from that repository
// Bindings closely repeat the C APIs and include both, a simplified
// single-function API to just open a full-screen webview window, and a more
// advanced and featureful set of APIs, including Go-to-JavaScript bindings.
//
// The library uses gtk-webkit, Cocoa/Webkit and MSHTML (IE8..11) as a browser
// engine and supports Linux, MacOS and Windows 7..10 respectively.
//
package webview

/*
#cgo linux openbsd freebsd CFLAGS: -DWEBVIEW_GTK=1 -Wno-deprecated-declarations
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0

#cgo windows CFLAGS: -DWEBVIEW_WINAPI=1 -std=c99
#cgo windows LDFLAGS: -lole32 -lcomctl32 -loleaut32 -luuid -lgdi32

#cgo darwin CFLAGS: -DWEBVIEW_COCOA=1 -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework WebKit

#include <stdlib.h>
#include <stdint.h>
#define WEBVIEW_STATIC
#define WEBVIEW_IMPLEMENTATION
#include "webview.h"

extern void _webviewExternalInvokeCallback(void *, void *);

static inline void CgoWebViewFree(void *w) {
	free((void *)((struct webview *)w)->title);
	free((void *)((struct webview *)w)->url);
	free(w);
}

static inline void *CgoWebViewCreate(int width, int height, char *title, char *url, int resizable, int debug) {
	struct webview *w = (struct webview *) calloc(1, sizeof(*w));
	w->width = width;
	w->height = height;
	w->title = title;
	w->url = url;
	w->resizable = resizable;
	w->debug = debug;
	w->external_invoke_cb = (webview_external_invoke_cb_t) _webviewExternalInvokeCallback;
	if (webview_init(w) != 0) {
		CgoWebViewFree(w);
		return NULL;
	}
	return (void *)w;
}

static inline int CgoWebViewLoop(void *w, int blocking) {
	return webview_loop((struct webview *)w, blocking);
}

static inline void CgoWebViewTerminate(void *w) {
	webview_terminate((struct webview *)w);
}

static inline void CgoWebViewExit(void *w) {
	webview_exit((struct webview *)w);
}

static inline void CgoWebViewSetTitle(void *w, char *title) {
	webview_set_title((struct webview *)w, title);
}

static inline void CgoWebViewFocus(void *w) {
	webview_focus((struct webview *)w);
}

static inline void CgoWebViewMinSize(void *w, int width, int height) {
	webview_minsize((struct webview *)w, width, height);
}

static inline void CgoWebViewMaxSize(void *w, int width, int height) {
	webview_maxsize((struct webview *)w, width, height);
}

static inline void CgoWebViewSetFullscreen(void *w, int fullscreen) {
	webview_set_fullscreen((struct webview *)w, fullscreen);
}

static inline void CgoWebViewSetColor(void *w, uint8_t r, uint8_t g, uint8_t b, uint8_t a) {
	webview_set_color((struct webview *)w, r, g, b, a);
}

static inline void CgoDialog(void *w, int dlgtype, int flags,
char *title, char *arg, char *res, size_t ressz, char *filter) {
	webview_dialog(w, dlgtype, flags,
	(const char*)title, (const char*) arg, res, ressz, filter);
}

static inline int CgoWebViewEval(void *w, char *js) {
	return webview_eval((struct webview *)w, js);
}

static inline void CgoWebViewInjectCSS(void *w, char *css) {
	webview_inject_css((struct webview *)w, css);
}

extern void _webviewDispatchGoCallback(void *);
static inline void _webview_dispatch_cb(struct webview *w, void *arg) {
	_webviewDispatchGoCallback(arg);
}
static inline void CgoWebViewDispatch(void *w, uintptr_t arg) {
	webview_dispatch((struct webview *)w, _webview_dispatch_cb, (void *)arg);
}
*/
import "C"
import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

func init() {
	// Ensure that main.main is called from the main thread
	runtime.LockOSThread()
}

// Open is a simplified API to open a single native window with a full-size webview in
// it. It can be helpful if you want to communicate with the core app using XHR
// or WebSockets (as opposed to using JavaScript bindings).
//
// Window appearance can be customized using title, width, height and resizable parameters.
// URL must be provided and can user either a http or https protocol, or be a
// local file:// URL. On some platforms "data:" URLs are also supported
// (Linux/MacOS).
func Open(title, url string, w, h int, resizable bool) error {
	titleStr := C.CString(title)
	defer C.free(unsafe.Pointer(titleStr))
	urlStr := C.CString(url)
	defer C.free(unsafe.Pointer(urlStr))
	resize := C.int(0)
	if resizable {
		resize = C.int(1)
	}

	r := C.webview(titleStr, urlStr, C.int(w), C.int(h), resize)
	if r != 0 {
		return errors.New("failed to create webview")
	}
	return nil
}

// ExternalInvokeCallbackFunc is a function type that is called every time
// "window.external.invoke()" is called from JavaScript. Data is the only
// obligatory string parameter passed into the "invoke(data)" function from
// JavaScript. To pass more complex data serialized JSON or base64 encoded
// string can be used.
type ExternalInvokeCallbackFunc func(w WebView, data string)

// Settings is a set of parameters to customize the initial WebView appearance
// and behavior. It is passed into the webview.New() constructor.
type Settings struct {
	// WebView main window title
	Title string
	// URL to open in a webview
	URL string
	// Window width in pixels
	Width int
	// Window height in pixels
	Height int
	// Allows/disallows window resizing
	Resizable bool
	// Enable debugging tools (Linux/BSD/MacOS, on Windows use Firebug)
	Debug bool
	// A callback that is executed when JavaScript calls "window.external.invoke()"
	ExternalInvokeCallback ExternalInvokeCallbackFunc
}

// WebView is an interface that wraps the basic methods for controlling the UI
// loop, handling multithreading and providing JavaScript bindings.
type WebView interface {
	// Run() starts the main UI loop until the user closes the webview window or
	// Terminate() is called.
	Run()
	// Loop() runs a single iteration of the main UI.
	Loop(blocking bool) bool
	// SetTitle() changes window title. This method must be called from the main
	// thread only. See Dispatch() for more details.
	SetTitle(title string)

	// Focus() puts the main window into focus
	Focus()

	// SetMinSize() sets the minimum size of the window
	SetMinSize(width, height int)

	// SetMaxSize() sets the maximum size of the window
	SetMaxSize(width, height int)

	// SetFullscreen() controls window full-screen mode. This method must be
	// called from the main thread only. See Dispatch() for more details.
	SetFullscreen(fullscreen bool)
	// SetColor() changes window background color. This method must be called from
	// the main thread only. See Dispatch() for more details.
	SetColor(r, g, b, a uint8)
	// Eval() evaluates an arbitrary JS code inside the webview. This method must
	// be called from the main thread only. See Dispatch() for more details.
	Eval(js string) error
	// InjectJS() injects an arbitrary block of CSS code using the JS API. This
	// method must be called from the main thread only. See Dispatch() for more
	// details.
	InjectCSS(css string)
	// Dialog() opens a system dialog of the given type and title. String
	// argument can be provided for certain dialogs, such as alert boxes. For
	// alert boxes argument is a message inside the dialog box.
	Dialog(dlgType DialogType, flags int, title string, arg string, filter string) string
	// Terminate() breaks the main UI loop. This method must be called from the main thread
	// only. See Dispatch() for more details.
	Terminate()
	// Dispatch() schedules some arbitrary function to be executed on the main UI
	// thread. This may be helpful if you want to run some JavaScript from
	// background threads/goroutines, or to terminate the app.
	Dispatch(func())
	// Exit() closes the window and cleans up the resources. Use Terminate() to
	// forcefully break out of the main UI loop.
	Exit()
}

// DialogType is an enumeration of all supported system dialog types
type DialogType int

const (
	// DialogTypeOpen is a system file open dialog
	DialogTypeOpen DialogType = iota
	// DialogTypeSave is a system file save dialog
	DialogTypeSave
	// DialogTypeAlert is a system alert dialog (message box)
	DialogTypeAlert
)

const (
	// DialogFlagFile is a normal file picker dialog
	DialogFlagFile = C.WEBVIEW_DIALOG_FLAG_FILE
	// DialogFlagDirectory is an open directory dialog
	DialogFlagDirectory = C.WEBVIEW_DIALOG_FLAG_DIRECTORY
	// DialogFlagInfo is an info alert dialog
	DialogFlagInfo = C.WEBVIEW_DIALOG_FLAG_INFO
	// DialogFlagWarning is a warning alert dialog
	DialogFlagWarning = C.WEBVIEW_DIALOG_FLAG_WARNING
	// DialogFlagError is an error dialog
	DialogFlagError = C.WEBVIEW_DIALOG_FLAG_ERROR
)

var (
	m     sync.Mutex
	index uintptr
	fns   = map[uintptr]func(){}
	cbs   = map[WebView]ExternalInvokeCallbackFunc{}
)

type webview struct {
	w unsafe.Pointer
}

var _ WebView = &webview{}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// NewWebview creates and opens a new webview window using the given settings. The
// returned object implements the WebView interface. This function returns nil
// if a window can not be created.
func NewWebview(settings Settings) WebView {
	if settings.Width == 0 {
		settings.Width = 640
	}
	if settings.Height == 0 {
		settings.Height = 480
	}
	if settings.Title == "" {
		settings.Title = "WebView"
	}
	w := &webview{}
	w.w = C.CgoWebViewCreate(C.int(settings.Width), C.int(settings.Height),
		C.CString(settings.Title), C.CString(settings.URL),
		C.int(boolToInt(settings.Resizable)), C.int(boolToInt(settings.Debug)))
	m.Lock()
	if settings.ExternalInvokeCallback != nil {
		cbs[w] = settings.ExternalInvokeCallback
	} else {
		cbs[w] = func(w WebView, data string) {}
	}
	m.Unlock()
	return w
}

func (w *webview) Loop(blocking bool) bool {
	block := C.int(0)
	if blocking {
		block = 1
	}
	return C.CgoWebViewLoop(w.w, block) == 0
}

func (w *webview) Run() {
	for w.Loop(true) {
	}
}

func (w *webview) Exit() {
	C.CgoWebViewExit(w.w)
}

func (w *webview) Dispatch(f func()) {
	m.Lock()
	for ; fns[index] != nil; index++ {
	}
	fns[index] = f
	m.Unlock()
	C.CgoWebViewDispatch(w.w, C.uintptr_t(index))
}

func (w *webview) SetTitle(title string) {
	p := C.CString(title)
	defer C.free(unsafe.Pointer(p))
	C.CgoWebViewSetTitle(w.w, p)
}

func (w *webview) SetColor(r, g, b, a uint8) {
	C.CgoWebViewSetColor(w.w, C.uint8_t(r), C.uint8_t(g), C.uint8_t(b), C.uint8_t(a))
}

func (w *webview) Focus() {
	C.CgoWebViewFocus(w.w)
}

func (w *webview) SetMinSize(width, height int) {
	C.CgoWebViewMinSize(w.w, C.int(width), C.int(height))
}

func (w *webview) SetMaxSize(width, height int) {
	C.CgoWebViewMaxSize(w.w, C.int(width), C.int(height))
}

func (w *webview) SetFullscreen(fullscreen bool) {
	C.CgoWebViewSetFullscreen(w.w, C.int(boolToInt(fullscreen)))
}

func (w *webview) Dialog(dlgType DialogType, flags int, title string, arg string, filter string) string {
	const maxPath = 4096
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))
	argPtr := C.CString(arg)
	defer C.free(unsafe.Pointer(argPtr))
	resultPtr := (*C.char)(C.calloc((C.size_t)(unsafe.Sizeof((*C.char)(nil))), (C.size_t)(maxPath)))
	defer C.free(unsafe.Pointer(resultPtr))
	filterPtr := C.CString(filter)
	defer C.free(unsafe.Pointer(filterPtr))
	C.CgoDialog(w.w, C.int(dlgType), C.int(flags), titlePtr,
		argPtr, resultPtr, C.size_t(maxPath), filterPtr)
	return C.GoString(resultPtr)
}

func (w *webview) Eval(js string) error {
	p := C.CString(js)
	defer C.free(unsafe.Pointer(p))
	switch C.CgoWebViewEval(w.w, p) {
	case -1:
		return errors.New("evaluation failed")
	}
	return nil
}

func (w *webview) InjectCSS(css string) {
	p := C.CString(css)
	defer C.free(unsafe.Pointer(p))
	C.CgoWebViewInjectCSS(w.w, p)
}

func (w *webview) Terminate() {
	C.CgoWebViewTerminate(w.w)
}

//export _webviewDispatchGoCallback
func _webviewDispatchGoCallback(index unsafe.Pointer) {
	var f func()
	m.Lock()
	f = fns[uintptr(index)]
	delete(fns, uintptr(index))
	m.Unlock()
	if f != nil {
		f()
	}
}

//export _webviewExternalInvokeCallback
func _webviewExternalInvokeCallback(w unsafe.Pointer, data unsafe.Pointer) {
	m.Lock()
	var (
		cb ExternalInvokeCallbackFunc
		wv WebView
	)
	for wv, cb = range cbs {
		if wv.(*webview).w == w {
			break
		}
	}
	m.Unlock()
	if cb != nil {
		cb(wv, C.GoString((*C.char)(data)))
	}
}
