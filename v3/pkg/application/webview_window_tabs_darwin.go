//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_tabs_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func macWindowTabsAdd(nsWindow, other unsafe.Pointer, order MacTabOrder) error {
	var result C.int
	InvokeSync(func() {
		result = C.windowTabsAdd(nsWindow, other, C.int(order))
	})
	switch result {
	case C.WailsWindowTabsAdded:
		return nil
	case C.WailsWindowTabsInvalidOrder:
		return fmt.Errorf("unknown macOS tab order %d", order)
	default:
		return ErrMacWindowTabTargetRequired
	}
}

func macWindowTabsHasGroup(nsWindow unsafe.Pointer) bool {
	return InvokeSyncWithResult(func() bool {
		return bool(C.windowTabsHasGroup(nsWindow))
	})
}

func macWindowTabsGroupPointer(nsWindow unsafe.Pointer) unsafe.Pointer {
	return InvokeSyncWithResult(func() unsafe.Pointer {
		return C.windowTabsGroupPointer(nsWindow)
	})
}

func macWindowTabsGroupWindows(nsWindow unsafe.Pointer) []unsafe.Pointer {
	return InvokeSyncWithResult(func() []unsafe.Pointer {
		count := int(C.windowTabsGroupCount(nsWindow))
		if count == 0 {
			return nil
		}
		windows := make([]unsafe.Pointer, 0, count)
		for i := 0; i < count; i++ {
			if pointer := C.windowTabsGroupWindowAt(nsWindow, C.int(i)); pointer != nil {
				windows = append(windows, pointer)
			}
		}
		return windows
	})
}

func macWindowTabsSelectedWindow(nsWindow unsafe.Pointer) unsafe.Pointer {
	return InvokeSyncWithResult(func() unsafe.Pointer {
		return C.windowTabsSelectedWindow(nsWindow)
	})
}

func macWindowTabsSelectNext(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsSelectNext(nsWindow) })
}

func macWindowTabsSelectPrevious(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsSelectPrevious(nsWindow) })
}

func macWindowTabsSelect(nsWindow, target unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsSelect(nsWindow, target) })
}

func macWindowTabsIsTabBarVisible(nsWindow unsafe.Pointer) bool {
	return InvokeSyncWithResult(func() bool {
		return bool(C.windowTabsIsTabBarVisible(nsWindow))
	})
}

func macWindowTabsToggleTabBar(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsToggleTabBar(nsWindow) })
}

func macWindowTabsIsOverviewVisible(nsWindow unsafe.Pointer) bool {
	return InvokeSyncWithResult(func() bool {
		return bool(C.windowTabsIsOverviewVisible(nsWindow))
	})
}

func macWindowTabsToggleOverview(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsToggleOverview(nsWindow) })
}

func macWindowTabsIdentifier(nsWindow unsafe.Pointer) string {
	return InvokeSyncWithResult(func() string {
		identifier := C.windowTabsIdentifier(nsWindow)
		if identifier == nil {
			return ""
		}
		defer C.free(unsafe.Pointer(identifier))
		return C.GoString(identifier)
	})
}

func macWindowTabsMoveToNewWindow(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsMoveToNewWindow(nsWindow) })
}

func macWindowTabsMergeAllWindows(nsWindow unsafe.Pointer) {
	InvokeSync(func() { C.windowTabsMergeAllWindows(nsWindow) })
}

func macWindowTabsSetTitle(nsWindow unsafe.Pointer, title string) {
	titleC := C.CString(title)
	defer C.free(unsafe.Pointer(titleC))
	InvokeSync(func() { C.windowTabsSetTitle(nsWindow, titleC) })
}

func macWindowTabsSetToolTip(nsWindow unsafe.Pointer, tooltip string) {
	tooltipC := C.CString(tooltip)
	defer C.free(unsafe.Pointer(tooltipC))
	InvokeSync(func() { C.windowTabsSetToolTip(nsWindow, tooltipC) })
}
