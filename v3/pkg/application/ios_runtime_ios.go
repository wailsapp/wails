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

	json "github.com/goccy/go-json"
)

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

// Live mutations
func iosSetScrollEnabled(enabled bool) { C.ios_runtime_set_scroll_enabled(C.bool(enabled)) }
func iosSetBounceEnabled(enabled bool) { C.ios_runtime_set_bounce_enabled(C.bool(enabled)) }
func iosSetScrollIndicatorsEnabled(enabled bool) {
	C.ios_runtime_set_scroll_indicators_enabled(C.bool(enabled))
}
func iosSetBackForwardGesturesEnabled(enabled bool) {
	C.ios_runtime_set_back_forward_gestures_enabled(C.bool(enabled))
}
func iosSetLinkPreviewEnabled(enabled bool) { C.ios_runtime_set_link_preview_enabled(C.bool(enabled)) }
func iosSetInspectableEnabled(enabled bool) { C.ios_runtime_set_inspectable_enabled(C.bool(enabled)) }
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
