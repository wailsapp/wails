//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c -fmodules -fobjc-arc
#cgo LDFLAGS: -framework UIKit
#include <stdlib.h>
#include "application_ios.h"
*/
import "C"
import (
	"unsafe"

	"encoding/json"
)

// iosHapticsImpact triggers an iOS haptic impact using the provided style.
// The style parameter specifies the impact style name understood by the native haptic engine.
func iosHapticsImpact(style string) {
	cstr := C.CString(style)
	defer C.free(unsafe.Pointer(cstr))
	C.ios_haptics_impact(cstr)
}

type deviceInfo struct {
	Model         string `json:"model"`
	SystemName    string `json:"systemName"`
	SystemVersion string `json:"systemVersion"`
	IsSimulator   bool   `json:"isSimulator"`
}

func iosDeviceInfo() deviceInfo {
	ptr := C.ios_device_info_json()
	if ptr == nil {
		return deviceInfo{}
	}
	defer C.free(unsafe.Pointer(ptr))
	goStr := C.GoString(ptr)
	var out deviceInfo
	_ = json.Unmarshal([]byte(goStr), &out)
	return out
}

// iosSetScrollEnabled sets whether scrolling is enabled in the iOS runtime.
func iosSetScrollEnabled(enabled bool) { C.ios_runtime_set_scroll_enabled(C.bool(enabled)) }
// iosSetBounceEnabled sets whether scroll bounce (rubber-band) behavior is enabled at runtime.
// If enabled is true, scrollable content will bounce when pulled past its edges; if false, that bounce is disabled.
func iosSetBounceEnabled(enabled bool) { C.ios_runtime_set_bounce_enabled(C.bool(enabled)) }
// iosSetScrollIndicatorsEnabled configures whether the iOS runtime shows scroll indicators.
// The enabled parameter controls visibility: true shows indicators, false hides them.
func iosSetScrollIndicatorsEnabled(enabled bool) {
	C.ios_runtime_set_scroll_indicators_enabled(C.bool(enabled))
}
// iosSetBackForwardGesturesEnabled enables back-forward navigation gestures when enabled is true and disables them when enabled is false.
func iosSetBackForwardGesturesEnabled(enabled bool) {
	C.ios_runtime_set_back_forward_gestures_enabled(C.bool(enabled))
}
// iosSetLinkPreviewEnabled sets whether link previews are enabled in the iOS runtime.
// Pass true to enable link previews, false to disable them.
func iosSetLinkPreviewEnabled(enabled bool) { C.ios_runtime_set_link_preview_enabled(C.bool(enabled)) }
// iosSetInspectableEnabled sets whether runtime web content inspection is enabled.
// When enabled is true the runtime allows inspection of web content; when false inspection is disabled.
func iosSetInspectableEnabled(enabled bool) { C.ios_runtime_set_inspectable_enabled(C.bool(enabled)) }
// iosSetCustomUserAgent sets the runtime's custom User-Agent string.
// If ua is an empty string, the custom User-Agent is cleared.
func iosSetCustomUserAgent(ua string) {
	var cstr *C.char
	if ua != "" {
		cstr = C.CString(ua)
		defer C.free(unsafe.Pointer(cstr))
	}
	C.ios_runtime_set_custom_user_agent(cstr)
}

// Native tabs
func iosSetNativeTabsEnabled(enabled bool) { C.ios_native_tabs_set_enabled(C.bool(enabled)) }
func iosNativeTabsIsEnabled() bool         { return bool(C.ios_native_tabs_is_enabled()) }
func iosSelectNativeTab(index int)         { C.ios_native_tabs_select_index(C.int(index)) }