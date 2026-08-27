//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "mac_window_chrome_darwin.h"
#include "webview_window_toolbar_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/base64"
	"encoding/json"
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

//export processToolbarShareResult
func processToolbarShareResult(itemID C.uint, service *C.char, errorMessage *C.char) {
	toolbarShareCompleted <- toolbarShareEvent{
		itemID:  uint(itemID),
		service: C.GoString(service),
		err:     C.GoString(errorMessage),
	}
}

//export processToolbarShareData
func processToolbarShareData(providerID C.uint, contentType *C.char) *C.char {
	type response struct {
		Data  string `json:"data"`
		Error string `json:"error,omitempty"`
	}
	data, err := handleToolbarShareData(uint(providerID), MacShareContentType(C.GoString(contentType)))
	payload := response{}
	if err != nil {
		payload.Error = err.Error()
	} else {
		payload.Data = base64.StdEncoding.EncodeToString(data)
	}
	encoded, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return C.CString(`{"error":"failed to encode share data"}`)
	}
	return C.CString(string(encoded))
}

//export processToolbarShareProviderRelease
func processToolbarShareProviderRelease(providerID C.uint) {
	releaseToolbarShareProvider(uint(providerID))
}

func (w *macosWebviewWindow) setToolbar(toolbar *MacToolbar) error {
	attached, err := attachMacToolbar(
		w.parent,
		w.nsWindow,
		w.activeToolbar,
		toolbar,
		w.hasSidebarSplitLayout(),
		w.hasInspectorSplitLayout(),
		w.parent.options.Mac.TitleBar,
	)
	if err == nil {
		w.activeToolbar = attached
	}
	return err
}

func attachMacToolbar(owner macToolbarWindow, nsWindow unsafe.Pointer, previousToolbar, toolbar *MacToolbar,
	hasSidebar, hasInspector bool, titleBar MacTitleBar) (*MacToolbar, error) {
	if toolbar == nil {
		C.toolbarDetach(nsWindow)
		if previousToolbar != nil {
			clearMacToolbarState(previousToolbar, owner, true)
		}
		return nil, nil
	}

	// A tracking separator is meaningful only above a native sidebar divider.
	// Reject the attachment instead of installing a misplaced decorative
	// item; the split view is installed before toolbars, so a pending layout
	// is already in the window by the time a stashed toolbar attaches.
	if toolbar.hasSidebarTrackingSeparator() && !hasSidebar {
		return previousToolbar, fmt.Errorf("a sidebar tracking separator requires a split view with a sidebar; call SetSplitView before the window is shown")
	}
	if toolbar.hasInspectorChrome() && !hasInspector {
		return previousToolbar, fmt.Errorf("inspector toolbar items require a split view with an inspector; call SetSplitView before the window is shown")
	}

	identifier := C.CString(toolbar.identifier)
	handle := C.toolbarCreate(identifier)
	C.free(unsafe.Pointer(identifier))
	if handle == nil {
		return previousToolbar, fmt.Errorf("failed to create native toolbar")
	}

	var itemIDs []uint
	committed := false
	defer func() {
		if !committed {
			releaseMacToolbarResources(handle, itemIDs)
		}
	}()

	for _, item := range toolbar.itemSnapshot() {
		itemIDs = append(itemIDs, addToolbarItem(handle, item)...)
	}
	toolbar.stateLock.RLock()
	displayMode := toolbar.displayMode
	toolbar.stateLock.RUnlock()
	macToolbarSetDisplayMode(handle, displayMode)

	toolbar.stateLock.Lock()
	if toolbar.state == nil || toolbar.state.window != owner {
		toolbar.stateLock.Unlock()
		return previousToolbar, fmt.Errorf("toolbar ownership changed while attaching")
	}
	previousNative := toolbar.state.native
	previousItemIDs := toolbar.state.itemIDs
	// Commit and attach while ownership is locked. Candidate construction is
	// already complete, so every failure path above leaves the previous native
	// toolbar untouched.
	C.toolbarAttach(nsWindow, handle, C.int(titleBar.ToolbarStyle))
	toolbar.state.native = handle
	toolbar.state.itemIDs = itemIDs
	toolbar.stateLock.Unlock()

	committed = true

	if previousToolbar == toolbar {
		releaseMacToolbarResources(previousNative, previousItemIDs)
	} else if previousToolbar != nil {
		clearMacToolbarState(previousToolbar, owner, true)
	}

	// An item may be changed concurrently while the detached candidate is
	// being built. Reapply the latest snapshots after committing the native
	// handle so a setter that ran just before installation cannot be lost.
	applyMacToolbarLatestState(toolbar)

	// Apply each presentation preference explicitly after the real toolbar is
	// attached; no titlebar UseToolbar preference is required.
	C.windowSetToolbarStyle(nsWindow, C.int(titleBar.ToolbarStyle))
	C.windowSetShowToolbarWhenFullscreen(nsWindow, C.bool(titleBar.ShowToolbarWhenFullscreen))
	C.windowSetHideToolbarSeparator(nsWindow, C.bool(titleBar.HideToolbarSeparator))
	return toolbar, nil
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

type macToolbarItemSnapshot struct {
	identifier         string
	kind               macToolbarItemKind
	label              string
	symbolName         string
	tooltip            string
	bordered           bool
	prominent          bool
	tintColor          *RGBA
	badgeCount         int
	disabled           bool
	hidden             bool
	items              []*MacToolbarItem
	selectionMode      MacToolbarGroupSelectionMode
	selectedIndex      int
	shareFormats       []MacShareRepresentation
	shareProvider      MacShareProvider
	shareSubject       string
	shareSuggestedName string
}

func snapshotMacToolbarItem(item *MacToolbarItem) macToolbarItemSnapshot {
	item.lock.RLock()
	defer item.lock.RUnlock()
	result := macToolbarItemSnapshot{
		identifier:         item.identifier,
		kind:               item.kind,
		label:              item.label,
		symbolName:         item.symbolName,
		tooltip:            item.tooltip,
		bordered:           item.bordered,
		prominent:          item.prominent,
		badgeCount:         item.badgeCount,
		disabled:           item.disabled,
		hidden:             item.hidden,
		items:              append([]*MacToolbarItem(nil), item.items...),
		selectionMode:      item.selectionMode,
		selectedIndex:      item.selectedIndex,
		shareFormats:       append([]MacShareRepresentation(nil), item.shareFormats...),
		shareProvider:      item.shareProvider,
		shareSubject:       item.shareSubject,
		shareSuggestedName: item.shareSuggestedName,
	}
	if item.tintColor != nil {
		copyOfColor := *item.tintColor
		result.tintColor = &copyOfColor
	}
	return result
}

func addToolbarItem(handle unsafe.Pointer, item *MacToolbarItem) []uint {
	snapshot := snapshotMacToolbarItem(item)
	idCStr := C.CString(snapshot.identifier)
	defer C.free(unsafe.Pointer(idCStr))
	labelCStr := C.CString(snapshot.label)
	defer C.free(unsafe.Pointer(labelCStr))
	tooltipCStr := C.CString(snapshot.tooltip)
	defer C.free(unsafe.Pointer(tooltipCStr))
	var itemIDs []uint

	switch snapshot.kind {
	case toolbarButton:
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)

		symbolCStr := C.CString(snapshot.symbolName)
		defer C.free(unsafe.Pointer(symbolCStr))

		hasTint, tintR, tintG, tintB, tintA := macToolbarTintComponents(snapshot.tintColor)
		C.toolbarAddButtonItem(handle, idCStr, C.uint(id), labelCStr, symbolCStr,
			tooltipCStr, C.bool(snapshot.bordered), C.bool(snapshot.prominent),
			C.bool(snapshot.disabled), C.bool(snapshot.hidden), C.bool(hasTint),
			tintR, tintG, tintB, tintA, C.int(snapshot.badgeCount))

	case toolbarGroup:
		memberPtrs := make([]unsafe.Pointer, len(snapshot.items))
		for index, member := range snapshot.items {
			memberSnapshot := snapshotMacToolbarItem(member)
			memberID := nextToolbarNativeID()
			itemIDs = append(itemIDs, memberID)
			addToToolbarItemMap(memberID, member)

			memberIDStr := C.CString(memberSnapshot.identifier)
			memberLabelStr := C.CString(memberSnapshot.label)
			memberSymbolStr := C.CString(memberSnapshot.symbolName)
			memberTooltipStr := C.CString(memberSnapshot.tooltip)
			memberPtrs[index] = C.toolbarBuildButtonItemStandalone(memberIDStr, C.uint(memberID),
				memberLabelStr, memberSymbolStr, memberTooltipStr,
				C.bool(memberSnapshot.bordered), C.bool(memberSnapshot.disabled), C.bool(memberSnapshot.hidden))
			C.free(unsafe.Pointer(memberIDStr))
			C.free(unsafe.Pointer(memberLabelStr))
			C.free(unsafe.Pointer(memberSymbolStr))
			C.free(unsafe.Pointer(memberTooltipStr))
		}
		C.toolbarAddGroupItem(handle, idCStr, labelCStr,
			(*unsafe.Pointer)(unsafe.Pointer(&memberPtrs[0])), C.int(len(snapshot.items)),
			C.int(snapshot.selectionMode), C.int(snapshot.selectedIndex))

	case toolbarSearchField:
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)
		C.toolbarAddSearchItem(handle, idCStr, C.uint(id), labelCStr, tooltipCStr,
			C.bool(snapshot.disabled), C.bool(snapshot.hidden))

	case toolbarTitle:
		C.toolbarAddTitleItem(handle, idCStr, labelCStr, C.bool(snapshot.hidden))

	case toolbarSeparator:
		C.toolbarAddSeparatorIdentifier(handle)

	case toolbarShare:
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)
		symbolCStr := C.CString(snapshot.symbolName)
		defer C.free(unsafe.Pointer(symbolCStr))
		providerID := registerToolbarShareProvider(snapshot.shareProvider, snapshot.shareFormats, snapshot.shareSuggestedName)
		providerJSON := macToolbarShareProviderJSON(providerID, snapshot.shareSubject,
			snapshot.shareSuggestedName, snapshot.shareFormats)
		providerCStr := C.CString(providerJSON)
		defer C.free(unsafe.Pointer(providerCStr))
		C.toolbarAddShareItem(handle, idCStr, C.uint(id), labelCStr, symbolCStr,
			tooltipCStr, C.bool(snapshot.disabled || len(snapshot.shareFormats) == 0),
			C.bool(snapshot.hidden), providerCStr)

	case toolbarFlexibleSpace:
		C.toolbarAddFlexibleSpaceIdentifier(handle)

	case toolbarSidebarToggle:
		// A standard AppKit identifier: the toolbar creates and owns the
		// item, so no callback registration or native ID is needed.
		C.toolbarAddSidebarToggleIdentifier(handle)

	case toolbarSidebarTrackingSeparator:
		// Omitted natively on macOS releases without the tracking-separator
		// API; the sidebar and toggle still work there.
		C.toolbarAddSidebarTrackingSeparatorIdentifier(handle)

	case toolbarInspectorToggle:
		// AppKit owns the standard item on macOS 14+. On older releases the
		// native fallback routes through this private item registration.
		id := nextToolbarNativeID()
		itemIDs = append(itemIDs, id)
		addToToolbarItemMap(id, item)
		C.toolbarAddInspectorToggleItem(handle, idCStr, C.uint(id))

	case toolbarInspectorTrackingSeparator:
		C.toolbarAddInspectorTrackingSeparatorIdentifier(handle)
	}
	return itemIDs
}

// hasSidebarSplitLayout reports whether this window has a pending or
// installed split view containing a sidebar pane.
func (w *macosWebviewWindow) hasSidebarSplitLayout() bool {
	w.parent.splitViewLock.RLock()
	split := w.parent.splitView
	w.parent.splitViewLock.RUnlock()
	return split != nil && split.hasSidebarPane()
}

func (w *macosWebviewWindow) hasInspectorSplitLayout() bool {
	w.parent.splitViewLock.RLock()
	split := w.parent.splitView
	w.parent.splitViewLock.RUnlock()
	return split != nil && split.hasInspectorPane()
}

func applyMacToolbarLatestState(toolbar *MacToolbar) {
	toolbar.stateLock.RLock()
	defer toolbar.stateLock.RUnlock()
	if toolbar.state == nil || toolbar.state.native == nil {
		return
	}
	macToolbarSetDisplayMode(toolbar.state.native, toolbar.displayMode)
	for _, item := range toolbar.itemSnapshot() {
		applyMacToolbarItemLatestState(toolbar.state.native, item)
	}
}

func macToolbarSetDisplayMode(native unsafe.Pointer, mode MacToolbarDisplayMode) {
	C.toolbarSetDisplayMode(native, C.int(mode))
}

func applyMacToolbarItemLatestState(native unsafe.Pointer, item *MacToolbarItem) {
	snapshot := snapshotMacToolbarItem(item)
	// Standard AppKit identifiers (spaces, the sidebar toggle, the tracking
	// separator) are created and owned by AppKit and have no entry in the
	// delegate's item map.
	if snapshot.kind == toolbarFlexibleSpace ||
		snapshot.kind == toolbarSeparator ||
		snapshot.kind == toolbarSidebarToggle ||
		snapshot.kind == toolbarSidebarTrackingSeparator ||
		snapshot.kind == toolbarInspectorTrackingSeparator {
		return
	}
	// The inspector toggle is standard on macOS 14+ and has no mutable
	// delegate item there. Setters are harmless no-ops in that configuration;
	// on older releases they update the Wails-owned fallback.

	macToolbarItemSetLabel(native, snapshot.identifier, snapshot.label)
	if snapshot.symbolName != "" {
		macToolbarItemSetSymbol(native, snapshot.identifier, snapshot.symbolName)
	}
	macToolbarItemSetTooltip(native, snapshot.identifier, snapshot.tooltip)
	macToolbarItemSetBordered(native, snapshot.identifier, snapshot.bordered)
	macToolbarItemSetProminent(native, snapshot.identifier, snapshot.prominent)
	macToolbarItemSetTintColor(native, snapshot.identifier, snapshot.tintColor)
	enabled := !snapshot.disabled && (snapshot.kind != toolbarShare || len(snapshot.shareFormats) > 0)
	if snapshot.kind == toolbarShare {
		macToolbarShareItemSetProvider(native, snapshot.identifier, snapshot.shareProvider, snapshot.shareSubject,
			snapshot.shareSuggestedName, snapshot.shareFormats)
	}
	macToolbarItemSetEnabled(native, snapshot.identifier, enabled)
	macToolbarItemSetHidden(native, snapshot.identifier, snapshot.hidden)
	macToolbarItemSetBadgeCount(native, snapshot.identifier, snapshot.badgeCount)

	if snapshot.kind == toolbarGroup {
		macToolbarGroupSetSelectionMode(native, snapshot.identifier, snapshot.selectionMode)
		macToolbarGroupSetSelectedIndex(native, snapshot.identifier, snapshot.selectedIndex)
		for _, member := range snapshot.items {
			applyMacToolbarItemLatestState(native, member)
		}
	}
}

func macToolbarTintComponents(color *RGBA) (bool, C.double, C.double, C.double, C.double) {
	if color == nil {
		return false, 0, 0, 0, 0
	}
	return true,
		C.double(float64(color.Red) / 255.0),
		C.double(float64(color.Green) / 255.0),
		C.double(float64(color.Blue) / 255.0),
		C.double(float64(color.Alpha) / 255.0)
}

func macToolbarShareProviderJSON(providerID uint, subject, suggestedName string, formats []MacShareRepresentation) string {
	type representationPayload struct {
		ContentType MacShareContentType `json:"contentType"`
	}
	type providerPayload struct {
		ProviderID      uint                    `json:"providerID,omitempty"`
		Subject         string                  `json:"subject,omitempty"`
		SuggestedName   string                  `json:"suggestedName,omitempty"`
		Representations []representationPayload `json:"representations,omitempty"`
	}
	payload := providerPayload{ProviderID: providerID, Subject: subject, SuggestedName: suggestedName}
	for _, format := range formats {
		if format.ContentType != "" {
			payload.Representations = append(payload.Representations, representationPayload(format))
		}
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func releaseMacToolbarResources(native unsafe.Pointer, itemIDs []uint) {
	for _, id := range itemIDs {
		removeFromToolbarItemMap(id)
	}
	if native != nil {
		C.toolbarRelease(native)
	}
}

func clearMacToolbarState(toolbar *MacToolbar, window macToolbarWindow, releaseOwnership bool) {
	if toolbar == nil {
		return
	}
	toolbar.stateLock.Lock()
	state := toolbar.state
	if state == nil || (window != nil && state.window != window) {
		toolbar.stateLock.Unlock()
		return
	}
	native := state.native
	itemIDs := state.itemIDs
	state.native = nil
	state.itemIDs = nil
	if releaseOwnership {
		state.window = nil
	}
	toolbar.stateLock.Unlock()
	releaseMacToolbarResources(native, itemIDs)
}

func macToolbarItemSetLabel(native unsafe.Pointer, id, label string) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	labelC := C.CString(label)
	defer C.free(unsafe.Pointer(labelC))
	C.toolbarItemSetLabel(native, idC, labelC)
}

func macToolbarItemSetSymbol(native unsafe.Pointer, id, symbol string) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	symbolC := C.CString(symbol)
	defer C.free(unsafe.Pointer(symbolC))
	C.toolbarItemSetSymbol(native, idC, symbolC)
}

func macToolbarItemSetTooltip(native unsafe.Pointer, id, tooltip string) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	tooltipC := C.CString(tooltip)
	defer C.free(unsafe.Pointer(tooltipC))
	C.toolbarItemSetTooltip(native, idC, tooltipC)
}

func macToolbarItemSetBordered(native unsafe.Pointer, id string, bordered bool) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarItemSetBordered(native, idC, C.bool(bordered))
}

func macToolbarItemSetProminent(native unsafe.Pointer, id string, prominent bool) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarItemSetProminent(native, idC, C.bool(prominent))
}

func macToolbarItemSetTintColor(native unsafe.Pointer, id string, color *RGBA) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	hasTint, tintR, tintG, tintB, tintA := macToolbarTintComponents(color)
	C.toolbarItemSetTintColor(native, idC, C.bool(hasTint), tintR, tintG, tintB, tintA)
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

func macToolbarGroupSetSelectionMode(native unsafe.Pointer, id string, mode MacToolbarGroupSelectionMode) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	C.toolbarGroupSetSelectionMode(native, idC, C.int(mode))
}

func macToolbarShareItemSetProvider(native unsafe.Pointer, id string, provider MacShareProvider,
	subject, suggestedName string, formats []MacShareRepresentation) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))
	providerID := registerToolbarShareProvider(provider, formats, suggestedName)
	providerC := C.CString(macToolbarShareProviderJSON(providerID, subject, suggestedName, formats))
	defer C.free(unsafe.Pointer(providerC))
	C.toolbarShareItemSetProvider(native, idC, providerC)
}
