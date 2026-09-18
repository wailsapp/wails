//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework UniformTypeIdentifiers -mmacosx-version-min=10.13

#include <stdlib.h>
#include "environment_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
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

func defaultHandlerCStrings(handler DefaultHandler) (contentType *C.char, scheme *C.char, free func()) {
	if handler.ContentType != "" {
		contentType = C.CString(handler.ContentType)
	}
	if handler.URLScheme != "" {
		scheme = C.CString(handler.URLScheme)
	}
	return contentType, scheme, func() {
		C.free(unsafe.Pointer(contentType))
		C.free(unsafe.Pointer(scheme))
	}
}

func platformSetDefaultHandler(handler DefaultHandler) error {
	id, ch := newCompletionWaiter()
	contentType, scheme, free := defaultHandlerCStrings(handler)
	defer free()
	InvokeSync(func() {
		C.wailsWorkspaceSetDefaultHandler(C.ulonglong(id), contentType, scheme)
	})
	if err := awaitCompletion(id, ch, "environment: SetDefaultHandler"); err != nil {
		return fmt.Errorf("environment: SetDefaultHandler for %s: %w", handler, err)
	}
	return nil
}

func platformDefaultHandler(handler DefaultHandler) (AppInfo, error) {
	contentType, scheme, free := defaultHandlerCStrings(handler)
	defer free()
	raw := C.wailsWorkspaceDefaultHandlerJSON(contentType, scheme)
	defer C.free(unsafe.Pointer(raw))
	text := C.GoString(raw)
	if text == "null" || text == "" {
		return AppInfo{}, fmt.Errorf("%w: no default application for %s", ErrApplicationNotFound, handler)
	}
	var info AppInfo
	if err := json.Unmarshal([]byte(text), &info); err != nil {
		return AppInfo{}, fmt.Errorf("environment: cannot decode the default handler for %s: %w", handler, err)
	}
	return info, nil
}
