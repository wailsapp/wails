//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_sheet_darwin.h"
*/
import "C"

import "unsafe"

func macSheetsSupported() bool { return true }

func macSheetBegin(parent, sheet unsafe.Pointer, critical bool) error {
	result := InvokeSyncWithResult(func() C.int {
		return C.windowSheetBegin(parent, sheet, C.bool(critical))
	})
	switch result {
	case C.WailsWindowSheetBegan:
		return nil
	case C.WailsWindowSheetSelf:
		return ErrMacSheetSelf
	default:
		return ErrMacSheetRequired
	}
}

func macSheetEnd(sheet unsafe.Pointer, code int) {
	InvokeSync(func() { C.windowSheetEnd(sheet, C.long(code)) })
}

func macSheetParent(sheet unsafe.Pointer) unsafe.Pointer {
	return InvokeSyncWithResult(func() unsafe.Pointer {
		return C.windowSheetParent(sheet)
	})
}

func macSheetAttached(parent unsafe.Pointer) unsafe.Pointer {
	return InvokeSyncWithResult(func() unsafe.Pointer {
		return C.windowSheetAttached(parent)
	})
}

//export processMacSheetEnded
func processMacSheetEnded(sheet unsafe.Pointer, code C.long) {
	macSheetEnded <- macSheetEndEvent{sheet: sheet, code: int(code)}
}
