//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_darwin.h"
#include "webview_window_toolbar_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//export processToolbarItemClick
func processToolbarItemClick(itemID C.uint) {
	toolbarItemClicked <- uint(itemID)
}

//export processToolbarSearch
func processToolbarSearch(itemID C.uint, query *C.char) {
	toolbarSearchTriggered <- toolbarSearchEvent{itemID: uint(itemID), query: C.GoString(query)}
}

func (w *macosWebviewWindow) setToolbar(toolbar *MacToolbar) error {
	if w.activeToolbar != nil {
		clearMacToolbarState(w.activeToolbar)
		w.activeToolbar = nil
	}
	if toolbar == nil {
		C.toolbarDetach(w.nsWindow)
		return nil
	}
	delegatePtr := C.toolbarNewAndAttach(w.nsWindow)
	if delegatePtr == nil {
		return fmt.Errorf("failed to create native toolbar")
	}
	var itemIDs []uint
	for _, item := range toolbar.items {
		itemIDs = append(itemIDs, addToolbarItem(delegatePtr, item)...)
	}
	titleBar := w.parent.options.Mac.TitleBar
	C.toolbarReload(w.nsWindow, delegatePtr, C.int(titleBar.ToolbarStyle))
	toolbar.stateLock.Lock()
	toolbar.state.native = delegatePtr
	toolbar.state.itemIDs = itemIDs
	toolbar.stateLock.Unlock()
	// The toolbar may be attached after the window's titlebar options were
	// applied (especially when SetToolbar was called before app.Run()). Apply
	// the presentation settings again now that AppKit has a real toolbar.
	C.windowSetToolbarStyle(w.nsWindow, C.int(titleBar.ToolbarStyle))
	C.windowSetShowToolbarWhenFullscreen(w.nsWindow, C.bool(titleBar.ShowToolbarWhenFullscreen))
	C.windowSetHideToolbarSeparator(w.nsWindow, C.bool(titleBar.HideToolbarSeparator))
	w.activeToolbar = toolbar
	return nil
}

func (w *macosWebviewWindow) refreshToolbarAfterShow() {
	w.parent.toolbarLock.RLock()
	toolbar := w.parent.toolbar
	w.parent.toolbarLock.RUnlock()
	if toolbar == nil {
		return
	}
	if err := w.setToolbar(toolbar); err != nil {
		w.parent.Error("SetToolbar: %s", err)
	}
}

func addToolbarItem(delegatePtr unsafe.Pointer, item *MacToolbarItem) []uint {
	idCStr := C.CString(item.identifier)
	defer C.free(unsafe.Pointer(idCStr))
	labelCStr := C.CString(item.label)
	defer C.free(unsafe.Pointer(labelCStr))
	tooltipCStr := C.CString(item.tooltip)
	defer C.free(unsafe.Pointer(tooltipCStr))
	var itemIDs []uint

	switch item.kind {
	case ToolbarButton:
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)

		symbolCStr := C.CString(item.symbolName)
		defer C.free(unsafe.Pointer(symbolCStr))

		var tintR, tintG, tintB, tintA C.double
		hasTint := item.tintColor != nil
		if hasTint {
			tintR = C.double(float64(item.tintColor.Red) / 255.0)
			tintG = C.double(float64(item.tintColor.Green) / 255.0)
			tintB = C.double(float64(item.tintColor.Blue) / 255.0)
			tintA = C.double(float64(item.tintColor.Alpha) / 255.0)
		}
		C.toolbarAddButtonItem(delegatePtr, idCStr, C.uint(id), labelCStr, symbolCStr,
			tooltipCStr, C.bool(item.bordered), C.bool(item.prominent),
			C.bool(item.disabled), C.bool(item.hidden), C.bool(hasTint),
			tintR, tintG, tintB, tintA, C.int(item.badgeCount))

	case ToolbarGroup:
		memberPtrs := make([]unsafe.Pointer, len(item.items))
		for i, member := range item.items {
			memberID := nextToolbarNativeID()
			itemIDs = append(itemIDs, memberID)
			addToToolbarItemMap(memberID, member)

			memberIDStr := C.CString(member.identifier)
			memberLabelStr := C.CString(member.label)
			memberSymbolStr := C.CString(member.symbolName)
			memberTooltipStr := C.CString(member.tooltip)
			memberPtrs[i] = C.toolbarBuildButtonItemStandalone(memberIDStr, C.uint(memberID),
				memberLabelStr, memberSymbolStr, memberTooltipStr,
				C.bool(member.bordered), C.bool(member.disabled), C.bool(member.hidden))
			C.free(unsafe.Pointer(memberIDStr))
			C.free(unsafe.Pointer(memberLabelStr))
			C.free(unsafe.Pointer(memberSymbolStr))
			C.free(unsafe.Pointer(memberTooltipStr))
		}
		C.toolbarAddGroupItem(delegatePtr, idCStr, labelCStr,
			(*unsafe.Pointer)(unsafe.Pointer(&memberPtrs[0])), C.int(len(item.items)),
			C.int(item.selectionMode), C.int(item.selectedIndex))

	case ToolbarSearchField:
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)
		C.toolbarAddSearchItem(delegatePtr, idCStr, C.uint(id), labelCStr, tooltipCStr, C.bool(item.disabled), C.bool(item.hidden))

	case ToolbarFlexibleSpace:
		C.toolbarAddFlexibleSpaceIdentifier(delegatePtr)

	case ToolbarSidebarToggle:
		C.toolbarAddSidebarToggleIdentifier(delegatePtr, idCStr)
	case ToolbarInspectorToggle:
		C.toolbarAddInspectorToggleIdentifier(delegatePtr, idCStr)
	}
	return itemIDs
}

func clearMacToolbarState(toolbar *MacToolbar) {
	toolbar.stateLock.RLock()
	state := toolbar.state
	toolbar.stateLock.RUnlock()
	if state == nil {
		return
	}
	state.lock.Lock()
	ids := state.itemIDs
	state.native = nil
	state.itemIDs = nil
	state.lock.Unlock()
	for _, id := range ids {
		removeFromToolbarItemMap(id)
	}
}

func macToolbarItemSetLabel(native unsafe.Pointer, id, label string) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))
	C.toolbarItemSetLabel(native, idC, labelC)
}

func macToolbarItemSetEnabled(native unsafe.Pointer, id string, enabled bool) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarItemSetEnabled(native, idC, C.bool(enabled))
}

func macToolbarItemSetHidden(native unsafe.Pointer, id string, hidden bool) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarItemSetHidden(native, idC, C.bool(hidden))
}

func macToolbarItemSetBadgeCount(native unsafe.Pointer, id string, count int) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarItemSetBadgeCount(native, idC, C.int(count))
}

func macToolbarGroupSetSelectedIndex(native unsafe.Pointer, id string, index int) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarGroupSetSelectedIndex(native, idC, C.int(index))
}
