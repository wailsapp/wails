//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -mmacosx-version-min=10.13

#include <stdlib.h>
#include "environment_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func platformAccessibilitySettings() AccessibilitySettings {
	native := C.wailsAccessibilitySettings()
	return AccessibilitySettings{
		ReduceMotion:              bool(native.reduceMotion),
		ReduceTransparency:        bool(native.reduceTransparency),
		IncreaseContrast:          bool(native.increaseContrast),
		DifferentiateWithoutColor: bool(native.differentiateWithoutColor),
		InvertColors:              bool(native.invertColors),
		VoiceOverEnabled:          bool(native.voiceOverEnabled),
		SwitchControlEnabled:      bool(native.switchControlEnabled),
	}
}

func platformKeyboardLayout() KeyboardLayout {
	var layout KeyboardLayout
	raw := InvokeSyncWithResult(func() string {
		cJSON := C.wailsKeyboardLayoutJSON()
		defer C.free(unsafe.Pointer(cJSON))
		return C.GoString(cJSON)
	})
	_ = json.Unmarshal([]byte(raw), &layout)
	if layout.Languages == nil {
		layout.Languages = []string{}
	}
	return layout
}

func platformLocale() LocaleInfo {
	var info LocaleInfo
	cJSON := C.wailsLocaleJSON()
	defer C.free(unsafe.Pointer(cJSON))
	_ = json.Unmarshal([]byte(C.GoString(cJSON)), &info)
	if info.Preferred == nil {
		info.Preferred = []string{}
	}
	return info
}

//export environmentNotification
func environmentNotification(kind C.int) {
	switch kind {
	case 0:
		applicationEvents <- newApplicationEvent(events.Mac.ApplicationDidChangeAccessibilitySettings)
		applicationEvents <- newApplicationEvent(events.Common.AccessibilitySettingsChanged)
	case 1:
		applicationEvents <- newApplicationEvent(events.Mac.ApplicationDidChangeKeyboardLayout)
	case 2:
		applicationEvents <- newApplicationEvent(events.Mac.ApplicationDidChangeLocale)
	}
}
