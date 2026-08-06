//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_accessory_darwin.h"
*/
import "C"

import "unsafe"

func macAccessoryControllerKind(native unsafe.Pointer) (MacAccessoryViewControllerKind, error) {
	kind := MacAccessoryViewControllerKind(C.macAccessoryViewControllerKind(native))
	if kind == MacAccessoryViewControllerKindUnknown {
		return kind, ErrMacAccessoryControllerType
	}
	return kind, nil
}

func macAccessoryControllerSupportsScrollEdgeEffectStyle(native unsafe.Pointer) bool {
	return bool(C.macAccessoryViewControllerSupportsScrollEdgeEffectStyle(native))
}

func macAccessoryControllerSetScrollEdgeEffectStyle(native unsafe.Pointer, style MacScrollEdgeEffectStyle) error {
	switch C.macAccessoryViewControllerSetScrollEdgeEffectStyle(native, C.int(style)) {
	case C.WailsMacAccessoryStyleApplied:
		return nil
	case C.WailsMacAccessoryStyleUnavailable:
		return ErrMacScrollEdgeEffectStyleUnavailable
	default:
		return ErrMacAccessoryControllerType
	}
}

func macAccessoryControllerScrollEdgeEffectStyle(native unsafe.Pointer) (MacScrollEdgeEffectStyle, error) {
	value := int(C.macAccessoryViewControllerScrollEdgeEffectStyle(native))
	switch value {
	case int(MacScrollEdgeEffectStyleAutomatic), int(MacScrollEdgeEffectStyleSoft), int(MacScrollEdgeEffectStyleHard):
		return MacScrollEdgeEffectStyle(value), nil
	case C.WailsMacAccessoryStyleUnavailable:
		return MacScrollEdgeEffectStyleAutomatic, nil
	default:
		return MacScrollEdgeEffectStyleAutomatic, ErrMacAccessoryControllerType
	}
}
