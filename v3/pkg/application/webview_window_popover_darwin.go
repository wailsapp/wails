//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_popover_darwin.h"
#include "webview_window_accessory_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func macPopoversSupported() bool { return true }

// macPopoverEnsureNative creates the NSPopover on the application thread
// the first time it is needed, building the content strip from the options'
// accessory. It returns the native handle.
func macPopoverEnsureNative(p *MacPopover) (unsafe.Pointer, error) {
	return InvokeSyncWithResultAndError(func() (unsafe.Pointer, error) {
		p.lock.Lock()
		defer p.lock.Unlock()
		if p.destroyed {
			return nil, ErrMacPopoverDestroyed
		}
		if p.native != nil {
			return p.native, nil
		}
		var strip unsafe.Pointer
		if content := p.options.Content; content != nil {
			if err := content.claim(macAccessoryTarget{window: p}); err != nil {
				return nil, ErrMacPopoverContentAttached
			}
			controller, err := macPopoverBuildStrip(content)
			if err != nil {
				content.releaseClaim()
				return nil, err
			}
			strip = controller
		}
		native := C.popoverCreate(C.ulonglong(p.id), C.double(p.options.Width), C.double(p.options.Height),
			C.int(p.options.Behavior), C.bool(!p.options.DisableAnimation), strip)
		if native == nil {
			return nil, errors.New("failed to create the native popover")
		}
		// The accessory keeps the reference macAccessoryCreate returned and
		// releases it in MacAccessory.Remove; the content controller holds
		// its own, released with the popover in popoverRelease.
		p.native = native
		return native, nil
	})
}

// macPopoverBuildStrip creates the accessory's native controller and
// installs its controls, the same steps attachMacAccessoryNative performs
// for a titlebar host, without attaching it to a window. The caller owns
// one reference to the returned controller.
func macPopoverBuildStrip(accessory *MacAccessory) (unsafe.Pointer, error) {
	snapshot := accessory.snapshot()
	controller := C.macAccessoryCreate(C.int(C.WailsMacAccessoryKindTitlebar), C.int(C.WailsMacAccessoryLayoutBottom), C.double(accessory.resolvedHeight()))
	if controller == nil {
		return nil, errors.New("failed to create the native popover content")
	}
	registerMacAccessoryControls(accessory)
	for _, control := range snapshot.controls {
		addMacAccessoryNativeControl(controller, control.snapshot())
	}
	C.macAccessorySetHidden(controller, C.bool(snapshot.hidden))
	accessory.lock.Lock()
	accessory.native = controller
	accessory.controller = &MacAccessoryViewController{native: controller, kind: MacAccessoryViewControllerKindTitlebar}
	accessory.lock.Unlock()
	for _, control := range snapshot.controls {
		applyMacAccessoryControlLatestState(controller, control.snapshot())
	}
	return controller, nil
}

func macPopoverShowInWindow(native, nsWindow unsafe.Pointer, rect Rect, edge MacRectEdge) error {
	result := InvokeSyncWithResult(func() C.int {
		return C.popoverShowInWindow(native, nsWindow, C.double(rect.X), C.double(rect.Y),
			C.double(rect.Width), C.double(rect.Height), C.int(edge.nsRectEdge()))
	})
	return macPopoverShowResult(result)
}

// macPopoverShowFromToolbarItem runs on the application thread inside
// MacToolbarItem.update, so it calls C directly.
func macPopoverShowFromToolbarItem(native, nsWindow unsafe.Pointer, identifier string, edge MacRectEdge) error {
	identifierC := C.CString(identifier)
	defer C.free(unsafe.Pointer(identifierC))
	return macPopoverShowResult(C.popoverShowFromToolbarItem(native, nsWindow, identifierC, C.int(edge.nsRectEdge())))
}

func macPopoverShowFromStatusItem(native, statusItem unsafe.Pointer, edge MacRectEdge) error {
	result := InvokeSyncWithResult(func() C.int {
		return C.popoverShowFromStatusItem(native, statusItem, C.int(edge.nsRectEdge()))
	})
	return macPopoverShowResult(result)
}

func macPopoverShowResult(result C.int) error {
	switch result {
	case C.WailsPopoverShown:
		return nil
	case C.WailsPopoverAnchorUnavailable:
		return ErrMacPopoverAnchorUnavailable
	default:
		return ErrMacPopoverAnchorRequired
	}
}

func macPopoverClose(native unsafe.Pointer) {
	InvokeSync(func() { C.popoverClose(native) })
}

func macPopoverIsShown(native unsafe.Pointer) bool {
	return InvokeSyncWithResult(func() bool { return bool(C.popoverIsShown(native)) })
}

func macPopoverSetContentSize(native unsafe.Pointer, width, height float64) {
	InvokeSync(func() { C.popoverSetContentSize(native, C.double(width), C.double(height)) })
}

func macPopoverSetBehavior(native unsafe.Pointer, behavior MacPopoverBehavior) {
	InvokeSync(func() { C.popoverSetBehavior(native, C.int(behavior)) })
}

// macPopoverRelease closes and releases the native popover, if it exists.
func macPopoverRelease(p *MacPopover) {
	p.lock.RLock()
	created := p.native != nil
	p.lock.RUnlock()
	if !created {
		// Never shown: nothing native exists, and there may be no
		// application to dispatch through.
		return
	}
	InvokeSync(func() {
		p.lock.Lock()
		native := p.native
		p.native = nil
		p.lock.Unlock()
		if native != nil {
			C.popoverRelease(native)
		}
	})
}

// macSystemTrayStatusItem returns the tray's NSStatusItem, or nil before
// the tray runs.
func macSystemTrayStatusItem(tray *SystemTray) unsafe.Pointer {
	impl, ok := tray.impl.(*macosSystemTray)
	if !ok || impl == nil {
		return nil
	}
	return InvokeSyncWithResult(func() unsafe.Pointer { return impl.nsStatusItem })
}

//export processMacPopoverClosed
func processMacPopoverClosed(popoverID C.ulonglong) {
	macPopoverClosed <- uint64(popoverID)
}
