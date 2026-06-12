//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c -fmodules -fobjc-arc
#cgo LDFLAGS: -framework UIKit -framework LocalAuthentication -framework UserNotifications -framework AVFoundation -framework Security
#include <stdlib.h>
#include "mobile_features_ios.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

// This file exposes a set of "genuinely mobile" native capabilities to Wails
// applications running on iOS: share sheet, opening URLs, keeping the screen
// awake, the torch, safe-area insets, brightness, orientation, the status bar,
// biometric authentication, local notifications and Keychain-backed secure
// storage. Each maps to a small Objective-C bridge in mobile_features_ios.m.
//
// Results that arrive asynchronously (e.g. a biometric prompt) are delivered
// back to the frontend as custom events via iosEmitNativeEvent.

func cString(s string) (*C.char, func()) {
	c := C.CString(s)
	return c, func() { C.free(unsafe.Pointer(c)) }
}

// --- Phase A: one-way actions -------------------------------------------------

// IOSShare presents the native iOS share sheet (UIActivityViewController). The
// JSON payload may contain "text" and/or "url" keys.
func IOSShare(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_share(c)
}

// IOSOpenURL opens the given URL in the system browser (Safari).
func IOSOpenURL(url string) {
	c, free := cString(url)
	defer free()
	C.ios_open_url(c)
}

// IOSSetKeepAwake disables (true) or restores (false) the idle timer, keeping
// the screen on while the app is in the foreground.
func IOSSetKeepAwake(enabled bool) { C.ios_set_keep_awake(C.bool(enabled)) }

// IOSSetTorch turns the device torch (flashlight) on or off.
func IOSSetTorch(enabled bool) { C.ios_set_torch(C.bool(enabled)) }

// iosEmitNativeEvent is called from the Objective-C bridge to deliver an
// asynchronous result to the frontend as a custom event.
//
//export iosEmitNativeEvent
func iosEmitNativeEvent(cname *C.char, cjson *C.char) {
	name := C.GoString(cname)
	var data map[string]any
	if cjson != nil {
		if s := C.GoString(cjson); s != "" {
			_ = json.Unmarshal([]byte(s), &data)
		}
	}
	app := globalApplication
	if app == nil {
		return
	}
	app.Event.Emit(name, data)
}
