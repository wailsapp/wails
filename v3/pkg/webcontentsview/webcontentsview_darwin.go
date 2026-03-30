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
	"encoding/json"
	"errors"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type macosWebContentsView struct {
	parent   *WebContentsView
	nsView   unsafe.Pointer
	nsWindow unsafe.Pointer
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	cPageScript := C.CString(darwinPageAutomationScript)
	defer C.free(unsafe.Pointer(cPageScript))

	cAutomationScript := C.CString(darwinRuntimeAutomationScript)
	defer C.free(unsafe.Pointer(cAutomationScript))

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

	C.webContentsViewConfigureAutomation(view, C.uintptr_t(parent.id), cPageScript, cAutomationScript)

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

func (w *macosWebContentsView) takeSnapshot() string {
	ch := make(chan string, 1)
	id := registerSnapshotCallback(ch)
	application.InvokeSync(func() {
		C.webContentsViewTakeSnapshot(w.nsView, C.uintptr_t(id))
	})
	return <-ch
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

func (w *macosWebContentsView) automationEnsureReady() error {
	if w.nsView == nil {
		return ErrAutomationNotSupported
	}
	return nil
}

func (w *macosWebContentsView) automationNativeCapabilities() automationNativeCapabilities {
	return automationNativeCapabilities{
		PageRuntime:       true,
		AutomationRuntime: true,
		AsyncRuntime:      true,
		DOM:               true,
		Storage:           true,
		Cookies:           true,
		NetworkBasic:      true,
		NetworkProxy:      bool(C.webContentsViewAutomationSupportsProxyCapture()),
		Accessibility:     true,
		Inspection:        bool(C.webContentsViewAutomationSupportsInspection()),
		PDF:               true,
	}
}

func (w *macosWebContentsView) automationEvaluate(expression string, world automationExecutionWorld, awaitPromise bool) (automationRemoteObject, error) {
	cExpression := C.CString(expression)
	defer C.free(unsafe.Pointer(cExpression))

	var worldID C.int
	if world == automationExecutionWorldAutomation {
		worldID = 1
	}

	var result automationRemoteObject
	err := w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewAutomationEvaluate(w.nsView, cExpression, worldID, C.bool(awaitPromise), C.uintptr_t(callbackID))
		})
	}, &result)
	return result, err
}

func (w *macosWebContentsView) automationInvoke(method string, params json.RawMessage) (any, error) {
	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))

	paramJSON := string(params)
	if len(params) == 0 {
		paramJSON = "null"
	}
	cParams := C.CString(paramJSON)
	defer C.free(unsafe.Pointer(cParams))

	var result any
	err := w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewAutomationInvoke(w.nsView, cMethod, cParams, C.uintptr_t(callbackID))
		})
	}, &result)
	return result, err
}

func (w *macosWebContentsView) automationGetCookies() ([]AutomationCookie, error) {
	var result []AutomationCookie
	err := w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewGetCookies(w.nsView, C.uintptr_t(callbackID))
		})
	}, &result)
	return result, err
}

func (w *macosWebContentsView) automationSetCookie(cookie AutomationCookie) error {
	payload, err := json.Marshal(cookie)
	if err != nil {
		return err
	}
	cCookie := C.CString(string(payload))
	defer C.free(unsafe.Pointer(cCookie))

	return w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewSetCookie(w.nsView, cCookie, C.uintptr_t(callbackID))
		})
	}, nil)
}

func (w *macosWebContentsView) automationDeleteCookie(name, domain, path string) (bool, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cDomain := C.CString(domain)
	defer C.free(unsafe.Pointer(cDomain))

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var result struct {
		Deleted bool `json:"deleted"`
	}
	err := w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewDeleteCookie(w.nsView, cName, cDomain, cPath, C.uintptr_t(callbackID))
		})
	}, &result)
	return result.Deleted, err
}

func (w *macosWebContentsView) automationClearCookies() error {
	return w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewClearCookies(w.nsView, C.uintptr_t(callbackID))
		})
	}, nil)
}

func (w *macosWebContentsView) automationCreatePDF() (string, error) {
	var result string
	err := w.runAutomationCommand(func(callbackID uintptr) {
		application.InvokeSync(func() {
			C.webContentsViewCreatePDF(w.nsView, C.uintptr_t(callbackID))
		})
	}, &result)
	return result, err
}

func (w *macosWebContentsView) automationSetInspectable(enabled bool) error {
	if !bool(C.webContentsViewAutomationSupportsInspection()) {
		return ErrAutomationNotSupported
	}

	var ok bool
	application.InvokeSync(func() {
		ok = bool(C.webContentsViewAutomationSetInspectable(w.nsView, C.bool(enabled)))
	})
	if !ok {
		return errors.New("unable to update inspectable state")
	}
	return nil
}

func (w *macosWebContentsView) runAutomationCommand(invoke func(uintptr), target any) error {
	ch := make(chan automationCommandResult, 1)
	callbackID := registerAutomationCommandCallback(ch)
	invoke(callbackID)
	result := <-ch
	if result.err != "" {
		return errors.New(result.err)
	}
	if target == nil {
		return nil
	}
	if result.payload == "" {
		return nil
	}
	return json.Unmarshal([]byte(result.payload), target)
}

//export browserViewSnapshotCallback
func browserViewSnapshotCallback(callbackID C.uintptr_t, base64 *C.char) {
	id := uintptr(callbackID)
	str := ""
	if base64 != nil {
		str = C.GoString(base64)
	}
	dispatchSnapshotResult(id, str)
}

//export browserViewAutomationCommandCallback
func browserViewAutomationCommandCallback(callbackID C.uintptr_t, resultJSON *C.char, errMsg *C.char) {
	result := automationCommandResult{}
	if resultJSON != nil {
		result.payload = C.GoString(resultJSON)
	}
	if errMsg != nil {
		result.err = C.GoString(errMsg)
	}
	dispatchAutomationCommandResult(uintptr(callbackID), result)
}

//export browserViewAutomationEventCallback
func browserViewAutomationEventCallback(viewID C.uintptr_t, method *C.char, payloadJSON *C.char) {
	if method == nil {
		return
	}

	payload := ""
	if payloadJSON != nil {
		payload = C.GoString(payloadJSON)
	}
	dispatchAutomationEvent(uint(viewID), C.GoString(method), payload)
}
