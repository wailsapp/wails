//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
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
	if toolbar == nil {
		C.toolbarDetach(w.nsWindow)
		return nil
	}
	if err := validateToolbarItems(toolbar.Items); err != nil {
		return err
	}
	delegatePtr := C.toolbarNewAndAttach(w.nsWindow)
	for _, item := range toolbar.Items {
		addToolbarItem(delegatePtr, item)
	}
	return nil
}

func validateToolbarItems(items []MacToolbarItem) error {
	for _, item := range items {
		switch item.Kind {
		case ToolbarButton:
			if item.OnClick == nil {
				return fmt.Errorf("toolbar item %q: ToolbarButton requires OnClick", item.ID)
			}
		case ToolbarSearchField:
			if item.OnSearch == nil {
				return fmt.Errorf("toolbar item %q: ToolbarSearchField requires OnSearch", item.ID)
			}
		case ToolbarGroup:
			if len(item.Items) == 0 {
				return fmt.Errorf("toolbar item %q: ToolbarGroup requires at least one member item", item.ID)
			}
			for _, member := range item.Items {
				if member.OnClick == nil {
					return fmt.Errorf("toolbar item %q: group member %q requires OnClick", item.ID, member.ID)
				}
			}
		case ToolbarFlexibleSpace, ToolbarSidebarToggle:
			// no callback required
		}
	}
	return nil
}

func addToolbarItem(delegatePtr unsafe.Pointer, item MacToolbarItem) {
	idCStr := C.CString(item.ID)
	defer C.free(unsafe.Pointer(idCStr))
	labelCStr := C.CString(item.Label)
	defer C.free(unsafe.Pointer(labelCStr))

	switch item.Kind {
	case ToolbarButton:
		id := nextToolbarItemID()
		itemCopy := item
		addToToolbarItemMap(id, &itemCopy)

		symbolCStr := C.CString(item.SymbolName)
		defer C.free(unsafe.Pointer(symbolCStr))

		var tintR, tintG, tintB, tintA C.double
		hasTint := item.TintColor != nil
		if hasTint {
			tintR = C.double(float64(item.TintColor.Red) / 255.0)
			tintG = C.double(float64(item.TintColor.Green) / 255.0)
			tintB = C.double(float64(item.TintColor.Blue) / 255.0)
			tintA = C.double(float64(item.TintColor.Alpha) / 255.0)
		}
		C.toolbarAddButtonItem(delegatePtr, idCStr, C.uint(id), labelCStr, symbolCStr,
			C.bool(item.Bordered), C.bool(item.Prominent), C.bool(hasTint),
			tintR, tintG, tintB, tintA, C.int(item.BadgeCount))

	case ToolbarGroup:
		memberPtrs := make([]unsafe.Pointer, len(item.Items))
		for i, member := range item.Items {
			memberID := nextToolbarItemID()
			memberCopy := member
			addToToolbarItemMap(memberID, &memberCopy)

			memberIDStr := C.CString(member.ID)
			memberLabelStr := C.CString(member.Label)
			memberSymbolStr := C.CString(member.SymbolName)
			memberPtrs[i] = C.toolbarBuildButtonItemStandalone(memberIDStr, C.uint(memberID),
				memberLabelStr, memberSymbolStr, C.bool(member.Bordered))
			C.free(unsafe.Pointer(memberIDStr))
			C.free(unsafe.Pointer(memberLabelStr))
			C.free(unsafe.Pointer(memberSymbolStr))
		}
		C.toolbarAddGroupItem(delegatePtr, idCStr, labelCStr,
			(*unsafe.Pointer)(unsafe.Pointer(&memberPtrs[0])), C.int(len(item.Items)))

	case ToolbarSearchField:
		id := nextToolbarItemID()
		itemCopy := item
		addToToolbarItemMap(id, &itemCopy)
		C.toolbarAddSearchItem(delegatePtr, idCStr, C.uint(id), labelCStr)

	case ToolbarFlexibleSpace:
		C.toolbarAddFlexibleSpaceIdentifier(delegatePtr)

	case ToolbarSidebarToggle:
		C.toolbarAddSidebarToggleIdentifier(delegatePtr, idCStr)
	}
}
