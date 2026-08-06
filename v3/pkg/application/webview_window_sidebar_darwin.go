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

func addMacSidebarItemToNative(handle unsafe.Pointer, paneID, sectionID uint64, item macSidebarItemSnapshot) {
	labelC := C.CString(item.label)
	symbolC := C.CString(item.symbolName)
	tooltipC := C.CString(item.tooltip)
	C.splitViewSidebarAddItem(handle,
		C.ulonglong(paneID),
		C.ulonglong(sectionID),
		C.ulonglong(item.internalID),
		labelC,
		symbolC,
		tooltipC,
		C.bool(item.disabled),
		C.bool(item.hidden))
	C.free(unsafe.Pointer(labelC))
	C.free(unsafe.Pointer(symbolC))
	C.free(unsafe.Pointer(tooltipC))
}
