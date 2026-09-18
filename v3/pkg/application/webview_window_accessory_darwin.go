//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_accessory_darwin.h"
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

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

// Native control strips (MacAccessory).

const macAccessoriesSupported = true

//export processMacAccessoryClicked
func processMacAccessoryClicked(controlID C.ulonglong) {
	macAccessoryEvents <- macAccessoryEvent{controlID: uint64(controlID), kind: macAccessoryEventClick}
}

//export processMacAccessorySearched
func processMacAccessorySearched(controlID C.ulonglong, query *C.char) {
	macAccessoryEvents <- macAccessoryEvent{
		controlID: uint64(controlID),
		kind:      macAccessoryEventSearch,
		text:      C.GoString(query),
	}
}

//export processMacAccessorySegmentChanged
func processMacAccessorySegmentChanged(controlID C.ulonglong, selected C.int) {
	macAccessoryEvents <- macAccessoryEvent{
		controlID: uint64(controlID),
		kind:      macAccessoryEventSelection,
		index:     int(selected),
	}
}

// macAccessoryNSWindow returns the host window's NSWindow, or nil while it
// does not exist. It must run on the application thread, where the pointer
// is written.
func macAccessoryNSWindow(window macAccessoryWindow) unsafe.Pointer {
	switch host := window.(type) {
	case *WebviewWindow:
		impl, ok := host.impl.(*macosWebviewWindow)
		if !ok || impl == nil {
			return nil
		}
		return impl.nsWindow
	case *NativeWindow:
		impl := host.implementation()
		if impl == nil {
			return nil
		}
		return impl.nativeWindow()
	}
	return nil
}

// macAccessoryAttachTitlebar attaches immediately when the NSWindow exists
// and otherwise queues the accessory for the window's creation path. The
// decision is made on the application thread so it cannot interleave with
// native window creation, which also runs there.
func macAccessoryAttachTitlebar(window macAccessoryWindow, accessory *MacAccessory) error {
	var err error
	InvokeSync(func() {
		nsWindow := macAccessoryNSWindow(window)
		if nsWindow == nil {
			queueMacAccessory(window, accessory)
			return
		}
		err = attachMacAccessoryNative(accessory, nsWindow, nil, false)
	})
	return err
}

// macAccessoryAttachPane attaches to an installed pane's NSSplitViewItem.
func macAccessoryAttachPane(pane *MacSplitPane, accessory *MacAccessory, top bool) error {
	var err error
	attempted := false
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		attempted = true
		item := C.splitViewPaneItem(handle, paneID)
		if item == nil {
			err = fmt.Errorf("the split pane has no native split item")
			return
		}
		err = attachMacAccessoryNative(accessory, nil, item, top)
	})
	if !attempted && err == nil {
		err = ErrMacAccessoryPaneRequired
	}
	return err
}

// attachMacAccessoryNative builds the native controller and its controls
// from a snapshot, attaches it, and commits the handle. It runs on the
// application thread. Every failure path leaves no native state behind.
func attachMacAccessoryNative(accessory *MacAccessory, nsWindow, splitItem unsafe.Pointer, top bool) error {
	snapshot := accessory.snapshot()
	kind := MacAccessoryViewControllerKindTitlebar
	nativeKind := C.int(C.WailsMacAccessoryKindTitlebar)
	if splitItem != nil {
		if !bool(C.macAccessorySplitItemAccessoriesSupported()) {
			return ErrMacSplitItemAccessoryUnavailable
		}
		kind = MacAccessoryViewControllerKindSplitItem
		nativeKind = C.int(C.WailsMacAccessoryKindSplitItem)
	}
	controller := C.macAccessoryCreate(nativeKind, C.int(snapshot.layout), C.double(snapshot.height))
	if controller == nil {
		return fmt.Errorf("failed to create the native accessory view controller")
	}

	registerMacAccessoryControls(accessory)
	for _, control := range snapshot.controls {
		addMacAccessoryNativeControl(controller, control.snapshot())
	}
	C.macAccessorySetHidden(controller, C.bool(snapshot.hidden))
	if snapshot.fullScreenMinHeight > 0 {
		C.macAccessorySetFullScreenMinHeight(controller, C.double(snapshot.fullScreenMinHeight))
	}
	if snapshot.adjustsSizeSet {
		C.macAccessorySetAutomaticallyAdjustsSize(controller, C.bool(snapshot.adjustsSize))
	}
	C.macAccessorySetAppliesContentInsets(controller, C.bool(snapshot.contentInsets))

	attached := false
	if splitItem != nil {
		attached = bool(C.macAccessoryAttachToSplitItem(splitItem, controller, C.bool(top)))
	} else {
		attached = bool(C.macAccessoryAttachToWindow(nsWindow, controller))
	}
	if !attached {
		unregisterMacAccessoryControls(accessory)
		C.macAccessoryRelease(controller)
		return fmt.Errorf("failed to attach the native accessory")
	}

	accessory.lock.Lock()
	accessory.native = controller
	accessory.controller = &MacAccessoryViewController{native: controller, kind: kind}
	accessory.lock.Unlock()

	// A setter may have run between the snapshot and the commit; replay the
	// latest control state so nothing is lost.
	for _, control := range snapshot.controls {
		applyMacAccessoryControlLatestState(controller, control.snapshot())
	}
	accessory.applyScrollEdgeStyle()
	return nil
}

func addMacAccessoryNativeControl(controller unsafe.Pointer, snapshot macAccessoryControlSnapshot) {
	text := C.CString(snapshot.text)
	defer C.free(unsafe.Pointer(text))
	symbol := C.CString(snapshot.symbol)
	defer C.free(unsafe.Pointer(symbol))
	tooltip := C.CString(snapshot.tooltip)
	defer C.free(unsafe.Pointer(tooltip))
	id := C.ulonglong(snapshot.id)
	width := C.double(snapshot.width)
	switch snapshot.kind {
	case MacAccessoryControlSearch:
		placeholder := C.CString(snapshot.placeholder)
		defer C.free(unsafe.Pointer(placeholder))
		C.macAccessoryAddSearch(controller, id, placeholder, text, C.bool(snapshot.incremental),
			tooltip, C.bool(snapshot.disabled), C.bool(snapshot.hidden), width)
	case MacAccessoryControlSegmented:
		labels := C.CString(macAccessoryStringsJSON(snapshot.segments))
		defer C.free(unsafe.Pointer(labels))
		symbols := C.CString(macAccessoryStringsJSON(snapshot.segmentSymbols))
		defer C.free(unsafe.Pointer(symbols))
		C.macAccessoryAddSegmented(controller, id, labels, symbols, C.int(snapshot.selected),
			tooltip, C.bool(snapshot.disabled), C.bool(snapshot.hidden), width)
	case MacAccessoryControlButton:
		C.macAccessoryAddButton(controller, id, text, symbol, nil,
			tooltip, C.bool(snapshot.disabled), C.bool(snapshot.hidden), width)
	case MacAccessoryControlMenuButton:
		menu := macToolbarNativeMenu(snapshot.menu)
		if menu == nil {
			// A menu button without a menu is still a visible, inert button
			// so the strip keeps its layout; SetMenu can supply one later.
			menu = macAccessoryEmptyMenu()
		}
		C.macAccessoryAddButton(controller, id, text, symbol, menu,
			tooltip, C.bool(snapshot.disabled), C.bool(snapshot.hidden), width)
	case MacAccessoryControlLabel:
		C.macAccessoryAddLabel(controller, id, text, symbol, tooltip, C.bool(snapshot.hidden), width)
	case MacAccessoryControlFlexibleSpace:
		C.macAccessoryAddFlexibleSpace(controller, id)
	case MacAccessoryControlNativeView:
		C.macAccessoryAddNativeView(controller, id, snapshot.nativeView, C.bool(snapshot.hidden), width)
	}
}

// macAccessoryEmptyMenu is the placeholder menu shared by menu buttons
// created without a Menu.
var macAccessoryPlaceholderMenu *Menu

func macAccessoryEmptyMenu() unsafe.Pointer {
	if macAccessoryPlaceholderMenu == nil {
		macAccessoryPlaceholderMenu = NewMenu()
	}
	return macToolbarNativeMenu(macAccessoryPlaceholderMenu)
}

// applyMacAccessoryControlLatestState replays the mutable properties that a
// setter could have changed while the native control was being built.
func applyMacAccessoryControlLatestState(controller unsafe.Pointer, snapshot macAccessoryControlSnapshot) {
	id := snapshot.id
	macAccessoryControlSetEnabled(controller, id, !snapshot.disabled)
	macAccessoryControlSetHidden(controller, id, snapshot.hidden)
	macAccessoryControlSetTooltip(controller, id, snapshot.tooltip)
	switch snapshot.kind {
	case MacAccessoryControlSearch:
		macAccessoryControlSetPlaceholder(controller, id, snapshot.placeholder)
		macAccessoryControlSetIncremental(controller, id, snapshot.incremental)
		macAccessoryControlSetText(controller, id, snapshot.text)
	case MacAccessoryControlSegmented:
		macAccessoryControlSetSegments(controller, id, snapshot.segments, snapshot.segmentSymbols, snapshot.selected)
	case MacAccessoryControlButton, MacAccessoryControlLabel:
		macAccessoryControlSetText(controller, id, snapshot.text)
		macAccessoryControlSetSymbol(controller, id, snapshot.symbol)
	case MacAccessoryControlMenuButton:
		macAccessoryControlSetText(controller, id, snapshot.text)
		macAccessoryControlSetSymbol(controller, id, snapshot.symbol)
		if snapshot.menu != nil {
			macAccessoryControlSetMenu(controller, id, snapshot.menu)
		}
	}
	if snapshot.kind != MacAccessoryControlFlexibleSpace {
		macAccessoryControlSetWidth(controller, id, snapshot.width)
	}
}

// macAccessoryDetachNative removes the native controller on the application
// thread. It is a no-op for accessories that were never attached.
func macAccessoryDetachNative(accessory *MacAccessory) {
	if accessory == nil || globalApplication == nil {
		return
	}
	accessory.lock.RLock()
	attached := accessory.native != nil
	accessory.lock.RUnlock()
	if !attached {
		return
	}
	InvokeSync(func() {
		accessory.lock.Lock()
		native := accessory.native
		target := accessory.target
		accessory.native = nil
		accessory.controller = nil
		accessory.lock.Unlock()
		if native == nil {
			return
		}
		if target.pane != nil {
			// The split item is looked up through the live split so a pane
			// torn down with its window is never dereferenced; the controller
			// is simply released in that case.
			macSplitPaneWithNative(target.pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
				C.macAccessoryDetachFromSplitItem(C.splitViewPaneItem(handle, paneID), native)
			})
		} else {
			C.macAccessoryDetachFromWindow(native)
		}
		unregisterMacAccessoryControls(accessory)
		C.macAccessoryRelease(native)
	})
}

// flushMacTitlebarAccessories attaches every accessory queued for window.
// It runs on the application thread once the NSWindow exists.
func flushMacTitlebarAccessories(window macAccessoryWindow, nsWindow unsafe.Pointer) {
	if nsWindow == nil {
		return
	}
	for _, accessory := range takePendingMacAccessories(window) {
		if err := attachMacAccessoryNative(accessory, nsWindow, nil, false); err != nil {
			accessory.releaseClaim()
			window.Error("AddTitlebarAccessory: %s", err)
		}
	}
}

// flushMacSplitPaneAccessories attaches every accessory queued for the
// panes of an installed split view. It runs on the application thread.
func flushMacSplitPaneAccessories(split *MacSplitView) {
	if split == nil || !split.isInstalled() {
		return
	}
	for _, pane := range split.paneSnapshot() {
		for _, accessory := range takePendingMacAccessories(pane.MacSplitPane) {
			accessory.lock.RLock()
			top := accessory.target.top
			accessory.lock.RUnlock()
			if err := macAccessoryAttachPane(pane.MacSplitPane, accessory, top); err != nil {
				accessory.releaseClaim()
				reportMacSplitError(split.ownerWindow(), "pane accessory: %s", err)
			}
		}
	}
}

// flushPendingMacAccessories is called from macosWebviewWindow.run after the
// split view and toolbar are installed.
func (w *macosWebviewWindow) flushPendingMacAccessories() {
	flushMacTitlebarAccessories(w.parent, w.nsWindow)
	flushMacSplitPaneAccessories(w.activeSplitView)
}

// macAccessoryFlushNativeWindow runs after NativeWindow scheduled its native
// creation. Main-thread work is delivered in order, so by the time this
// closure runs the window and its split view exist (or creation failed and
// the accessories stay queued).
func macAccessoryFlushNativeWindow(window *NativeWindow) {
	if globalApplication == nil {
		return
	}
	InvokeSync(func() {
		impl := window.implementation()
		if impl == nil {
			return
		}
		flushMacTitlebarAccessories(window, impl.nativeWindow())
		window.lock.RLock()
		split := window.split
		window.lock.RUnlock()
		flushMacSplitPaneAccessories(split)
	})
}

func macAccessoryStringsJSON(values []string) string {
	if values == nil {
		values = []string{}
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func macAccessorySetHidden(native unsafe.Pointer, hidden bool) {
	C.macAccessorySetHidden(native, C.bool(hidden))
}

func macAccessorySetHeight(native unsafe.Pointer, height float64) {
	C.macAccessorySetHeight(native, C.double(height))
}

func macAccessorySetFullScreenMinHeight(native unsafe.Pointer, height float64) {
	C.macAccessorySetFullScreenMinHeight(native, C.double(height))
}

func macAccessorySetAutomaticallyAdjustsSize(native unsafe.Pointer, adjusts bool) {
	C.macAccessorySetAutomaticallyAdjustsSize(native, C.bool(adjusts))
}

func macAccessorySetAppliesContentInsets(native unsafe.Pointer, applies bool) {
	C.macAccessorySetAppliesContentInsets(native, C.bool(applies))
}

func macAccessoryControlSetEnabled(native unsafe.Pointer, id uint64, enabled bool) {
	C.macAccessoryControlSetEnabled(native, C.ulonglong(id), C.bool(enabled))
}

func macAccessoryControlSetHidden(native unsafe.Pointer, id uint64, hidden bool) {
	C.macAccessoryControlSetHidden(native, C.ulonglong(id), C.bool(hidden))
}

func macAccessoryControlSetTooltip(native unsafe.Pointer, id uint64, tooltip string) {
	value := C.CString(tooltip)
	defer C.free(unsafe.Pointer(value))
	C.macAccessoryControlSetTooltip(native, C.ulonglong(id), value)
}

func macAccessoryControlSetWidth(native unsafe.Pointer, id uint64, width float64) {
	C.macAccessoryControlSetWidth(native, C.ulonglong(id), C.double(width))
}

func macAccessoryControlSetText(native unsafe.Pointer, id uint64, text string) {
	value := C.CString(text)
	defer C.free(unsafe.Pointer(value))
	C.macAccessoryControlSetText(native, C.ulonglong(id), value)
}

func macAccessoryControlSetSymbol(native unsafe.Pointer, id uint64, symbol string) {
	value := C.CString(symbol)
	defer C.free(unsafe.Pointer(value))
	C.macAccessoryControlSetSymbol(native, C.ulonglong(id), value)
}

func macAccessoryControlSetPlaceholder(native unsafe.Pointer, id uint64, placeholder string) {
	value := C.CString(placeholder)
	defer C.free(unsafe.Pointer(value))
	C.macAccessoryControlSetPlaceholder(native, C.ulonglong(id), value)
}

func macAccessoryControlSetIncremental(native unsafe.Pointer, id uint64, incremental bool) {
	C.macAccessoryControlSetIncremental(native, C.ulonglong(id), C.bool(incremental))
}

func macAccessoryControlSetSegments(native unsafe.Pointer, id uint64, labels, symbols []string, selected int) {
	labelsC := C.CString(macAccessoryStringsJSON(labels))
	defer C.free(unsafe.Pointer(labelsC))
	symbolsC := C.CString(macAccessoryStringsJSON(symbols))
	defer C.free(unsafe.Pointer(symbolsC))
	C.macAccessoryControlSetSegments(native, C.ulonglong(id), labelsC, symbolsC, C.int(selected))
}

func macAccessoryControlSetSelectedSegment(native unsafe.Pointer, id uint64, selected int) {
	C.macAccessoryControlSetSelectedSegment(native, C.ulonglong(id), C.int(selected))
}

func macAccessoryControlSetMenu(native unsafe.Pointer, id uint64, menu *Menu) {
	C.macAccessoryControlSetMenu(native, C.ulonglong(id), macToolbarNativeMenu(menu))
}
