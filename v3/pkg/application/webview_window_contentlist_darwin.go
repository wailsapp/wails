//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_contentlist_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

func macContentListApplySnapshot(list *MacContentList) {
	if list == nil {
		return
	}
	snapshot := list.snapshot()
	list.lock.RLock()
	pane := list.pane
	list.lock.RUnlock()
	if pane == nil {
		return
	}
	if pane.split.isInstalled() {
		list.registerRows()
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		applyMacContentListSnapshotToNative(handle, uint64(paneID), snapshot)
	})
}

func macContentListApplyRow(row *MacContentListRow) {
	if row == nil || row.list == nil {
		return
	}
	row.list.lock.RLock()
	pane := row.list.pane
	row.list.lock.RUnlock()
	if pane == nil {
		return
	}
	snapshot := snapshotMacContentListRow(row)
	spec, release := newMacContentListRowSpec(snapshot)
	defer release()
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewContentListUpdateRow(handle, paneID, C.ulonglong(snapshot.internalID), spec)
	})
}

func macContentListApplySelection(list *MacContentList, ids []uint64) {
	if list == nil {
		return
	}
	list.lock.RLock()
	pane := list.pane
	list.lock.RUnlock()
	if pane == nil {
		return
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		setMacContentListNativeSelection(handle, uint64(paneID), ids)
	})
}

func setMacContentListNativeSelection(handle unsafe.Pointer, paneID uint64, ids []uint64) {
	if len(ids) == 0 {
		C.splitViewContentListSetSelection(handle, C.ulonglong(paneID), nil, 0)
		return
	}
	buffer := (*C.ulonglong)(C.calloc(C.size_t(len(ids)), C.size_t(unsafe.Sizeof(C.ulonglong(0)))))
	defer C.free(unsafe.Pointer(buffer))
	for index, id := range ids {
		unsafe.Slice(buffer, len(ids))[index] = C.ulonglong(id)
	}
	C.splitViewContentListSetSelection(handle, C.ulonglong(paneID), buffer, C.int(len(ids)))
}

// applyMacContentListSnapshotToNative replaces the pane's native model and
// reloads the table with its selection preserved.
func applyMacContentListSnapshotToNative(handle unsafe.Pointer, paneID uint64, snapshot macContentListSnapshot) {
	C.splitViewContentListReset(handle, C.ulonglong(paneID))
	emptyTextC := C.CString(snapshot.emptyText)
	C.splitViewContentListSetOptions(handle, C.ulonglong(paneID), C.WailsContentListOptions{
		allowsMultipleSelection: C.bool(snapshot.allowsMultipleSelection),
		sortable:                C.bool(snapshot.sortable),
		sortColumn:              C.int(snapshot.sortColumn),
		sortAscending:           C.bool(snapshot.sortAscending),
		rowHeight:               C.double(snapshot.rowHeight),
		style:                   C.int(snapshot.style),
		alternatingRows:         C.bool(snapshot.alternatingRows),
		emptyText:               emptyTextC,
		headerVisible:           C.bool(snapshot.headerVisible),
	})
	C.free(unsafe.Pointer(emptyTextC))
	for _, column := range snapshot.columns {
		titleC := C.CString(column.Title)
		C.splitViewContentListAddColumn(handle, C.ulonglong(paneID), titleC,
			C.double(column.Width), C.double(column.MinWidth), C.double(column.MaxWidth),
			C.bool(column.Sortable), C.int(column.Alignment))
		C.free(unsafe.Pointer(titleC))
	}
	for _, row := range snapshot.rows {
		spec, release := newMacContentListRowSpec(row)
		C.splitViewContentListAddRow(handle, C.ulonglong(paneID), C.ulonglong(row.internalID), spec)
		release()
	}
	setMacContentListNativeSelection(handle, paneID, snapshot.selectedIDs)
	C.splitViewContentListReload(handle, C.ulonglong(paneID))
}

// newMacContentListRowSpec builds the borrowed C row description. The
// returned release function frees every string.
func newMacContentListRowSpec(row macContentListRowSnapshot) (C.WailsContentListRowSpec, func()) {
	cells := row.cells
	if cells == nil {
		cells = []string{}
	}
	cellsJSON, err := json.Marshal(cells)
	if err != nil {
		cellsJSON = []byte("[]")
	}
	strings := []*C.char{
		C.CString(row.title),
		C.CString(row.subtitle),
		C.CString(row.detail),
		C.CString(row.symbol),
		C.CString(row.tooltip),
		C.CString(string(cellsJSON)),
	}
	spec := C.WailsContentListRowSpec{
		title:     strings[0],
		subtitle:  strings[1],
		detail:    strings[2],
		symbol:    strings[3],
		tooltip:   strings[4],
		badge:     C.int(row.badge),
		cellsJSON: strings[5],
		disabled:  C.bool(row.disabled),
		hidden:    C.bool(row.hidden),
	}
	return spec, func() {
		for _, value := range strings {
			C.free(unsafe.Pointer(value))
		}
	}
}

//export processMacContentListSelectionChanged
func processMacContentListSelectionChanged(paneID C.ulonglong, rowIDs *C.ulonglong, count C.int) {
	event := macContentListSelectionEvent{paneID: uint64(paneID)}
	if rowIDs != nil && count > 0 {
		event.rowIDs = make([]uint64, 0, int(count))
		for _, id := range unsafe.Slice(rowIDs, int(count)) {
			event.rowIDs = append(event.rowIDs, uint64(id))
		}
	}
	macContentListSelectionEvents <- event
}

//export processMacContentListRowActivated
func processMacContentListRowActivated(paneID, rowID C.ulonglong) {
	macContentListActivateEvents <- macContentListActivateEvent{paneID: uint64(paneID), rowID: uint64(rowID)}
}

//export processMacContentListSortChanged
func processMacContentListSortChanged(paneID C.ulonglong, column C.int, ascending C.bool) {
	macContentListSortEvents <- macContentListSortEvent{
		paneID:    uint64(paneID),
		column:    int(column),
		ascending: bool(ascending),
	}
}

// processMacContentListContextMenu runs on the application thread while
// AppKit waits for the menu to show. It returns the NSMenu for the clicked
// row, or NULL when nothing should appear. The menu stays owned by its Go
// Menu.
//
//export processMacContentListContextMenu
func processMacContentListContextMenu(paneID, rowID C.ulonglong) unsafe.Pointer {
	menu := resolveMacContentListContextMenu(uint64(paneID), uint64(rowID))
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
