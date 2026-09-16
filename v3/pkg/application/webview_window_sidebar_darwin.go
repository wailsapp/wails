//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import "unsafe"

func macSidebarApplySnapshot(sidebar *MacSidebar) {
	if sidebar == nil {
		return
	}
	snapshot := sidebar.snapshot()
	sidebar.lock.RLock()
	pane := sidebar.pane
	sidebar.lock.RUnlock()
	if pane == nil {
		return
	}
	if pane.split.isInstalled() {
		sidebar.registerItems()
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		applyMacSidebarSnapshotToNative(handle, uint64(paneID), snapshot)
	})
}

func macSidebarApplySelection(sidebar *MacSidebar, item *MacSidebarItem) {
	if sidebar == nil {
		return
	}
	sidebar.lock.RLock()
	pane := sidebar.pane
	sidebar.lock.RUnlock()
	if pane == nil {
		return
	}
	var itemID uint64
	if item != nil {
		itemID = item.internalID
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewSidebarSetSelectedItem(handle, paneID, C.ulonglong(itemID))
	})
}

func applyMacSidebarSnapshotToNative(handle unsafe.Pointer, paneID uint64, snapshot macSidebarSnapshot) {
	C.splitViewSidebarReset(handle, C.ulonglong(paneID))
	C.splitViewSidebarSetOptions(handle, C.ulonglong(paneID),
		C.bool(snapshot.allowsMultipleSelection), C.bool(snapshot.reorderable))
	for _, entry := range snapshot.entries {
		if entry.item != nil {
			addMacSidebarItemToNative(handle, paneID, 0, *entry.item)
			continue
		}
		if entry.section == nil {
			continue
		}
		labelC := C.CString(entry.section.label)
		C.splitViewSidebarAddSection(handle, C.ulonglong(paneID), C.ulonglong(entry.section.internalID), labelC)
		C.free(unsafe.Pointer(labelC))
		for _, item := range entry.section.items {
			addMacSidebarItemToNative(handle, paneID, entry.section.internalID, item)
		}
	}
	C.splitViewSidebarSetSelectedItem(handle, C.ulonglong(paneID), C.ulonglong(snapshot.selectedItemID))
}

// addMacSidebarItemToNative adds one row beneath parentID (a section, a row,
// or 0 for the root) and then its nested rows.
func addMacSidebarItemToNative(handle unsafe.Pointer, paneID, parentID uint64, item macSidebarItemSnapshot) {
	labelC := C.CString(item.label)
	symbolC := C.CString(item.symbolName)
	tooltipC := C.CString(item.tooltip)
	accessoryC := C.CString(item.accessorySymbol)
	var tint RGBA
	hasTint := item.tintColor != nil
	if hasTint {
		tint = *item.tintColor
	}
	C.splitViewSidebarAddItem(handle,
		C.ulonglong(paneID),
		C.ulonglong(parentID),
		C.ulonglong(item.internalID),
		labelC,
		symbolC,
		tooltipC,
		C.bool(item.disabled),
		C.bool(item.hidden),
		C.bool(item.expanded),
		C.bool(item.editable),
		C.int(item.badge),
		accessoryC,
		C.bool(hasTint),
		C.int(tint.Red), C.int(tint.Green), C.int(tint.Blue), C.int(tint.Alpha))
	C.free(unsafe.Pointer(labelC))
	C.free(unsafe.Pointer(symbolC))
	C.free(unsafe.Pointer(tooltipC))
	C.free(unsafe.Pointer(accessoryC))
	for _, child := range item.children {
		addMacSidebarItemToNative(handle, paneID, item.internalID, child)
	}
}

//export processMacSidebarSelectionChanged
func processMacSidebarSelectionChanged(paneID C.ulonglong, itemIDs *C.ulonglong, count C.int) {
	event := macSidebarSelectionEvent{paneID: uint64(paneID)}
	if itemIDs != nil && count > 0 {
		event.itemIDs = make([]uint64, 0, int(count))
		for _, id := range unsafe.Slice(itemIDs, int(count)) {
			event.itemIDs = append(event.itemIDs, uint64(id))
		}
	}
	macSidebarSelectionEvents <- event
}

//export processMacSidebarItemExpanded
func processMacSidebarItemExpanded(itemID C.ulonglong, expanded C.bool) {
	macSidebarExpandedEvents <- macSidebarExpandedEvent{itemID: uint64(itemID), expanded: bool(expanded)}
}

//export processMacSidebarItemRenamed
func processMacSidebarItemRenamed(itemID C.ulonglong, label *C.char) {
	macSidebarRenameEvents <- macSidebarRenameEvent{itemID: uint64(itemID), label: C.GoString(label)}
}

//export processMacSidebarItemMoved
func processMacSidebarItemMoved(paneID, itemID, parentID C.ulonglong, index C.int) {
	macSidebarMoveEvents <- macSidebarMoveEvent{
		paneID:   uint64(paneID),
		itemID:   uint64(itemID),
		parentID: uint64(parentID),
		index:    int(index),
	}
}

// processMacSidebarContextMenu runs on the application thread while AppKit
// waits for the menu to show. It returns the NSMenu for the clicked row, or
// NULL when nothing should appear. The menu stays owned by its Go Menu.
//
//export processMacSidebarContextMenu
func processMacSidebarContextMenu(paneID, itemID C.ulonglong) unsafe.Pointer {
	menu := resolveMacSidebarContextMenu(uint64(paneID), uint64(itemID))
	if menu == nil {
		return nil
	}
	menu.Update()
	impl, ok := menu.impl.(*macosMenu)
	if !ok || impl == nil {
		return nil
	}
	return impl.nsMenu
}
