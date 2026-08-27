//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

//export processMacSplitPaneLoaded
func processMacSplitPaneLoaded(paneID C.ulonglong) {
	if event, ok := newMacSplitPaneLoadEvent(uint64(paneID)); ok {
		splitPaneLoadedEvents <- event
	}
}

//export processMacSplitPaneNavigationStarted
func processMacSplitPaneNavigationStarted(paneID C.ulonglong) {
	handleMacSplitPaneNavigationStarted(uint64(paneID))
}

//export processMacSplitPaneCollapsed
func processMacSplitPaneCollapsed(paneID C.ulonglong, collapsed C.bool) {
	splitPaneCollapseEvents <- splitPaneCollapseEvent{paneID: uint64(paneID), collapsed: bool(collapsed)}
}

//export processMacSidebarItemSelected
func processMacSidebarItemSelected(itemID C.ulonglong) {
	macSidebarItemSelected <- uint64(itemID)
}

//export processMacInspectorTextChanged
func processMacInspectorTextChanged(controlID C.ulonglong, value *C.char) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorTextField,
		value:     C.GoString(value),
	}
}

//export processMacInspectorToggleChanged
func processMacInspectorToggleChanged(controlID C.ulonglong, checked C.bool) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorCheckbox,
		checked:   bool(checked),
	}
}

//export processMacInspectorSelectionChanged
func processMacInspectorSelectionChanged(controlID C.ulonglong, selectedIndex C.int) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorPopup,
		selected:  int(selectedIndex),
	}
}

//export processMacTextEditorChanged
func processMacTextEditorChanged(editorID C.ulonglong) {
	macTextEditorChanged <- uint64(editorID)
}

// installSplitView installs a pending native split layout. It runs on the
// application thread during window creation, after the primary WebView exists
// and before any toolbar is attached: a sidebar tracking separator requires
// its split view to already be in the window.
func (w *macosWebviewWindow) installSplitView() {
	w.parent.splitViewLock.RLock()
	split := w.parent.splitView
	w.parent.splitViewLock.RUnlock()
	if split == nil {
		return
	}

	split.lock.RLock()
	autosaveName := split.autosaveName
	split.lock.RUnlock()
	panes := split.paneSnapshot()

	autosaveC := C.CString(autosaveName)
	handle := C.splitViewCreate(autosaveC)
	C.free(unsafe.Pointer(autosaveC))
	if handle == nil {
		w.parent.Error("SetSplitView: failed to create the native split view")
		return
	}

	for _, pane := range panes {
		snapshot := snapshotMacSplitPane(pane.MacSplitPane)
		pane.lock.RLock()
		paneLayout := pane.contentLayout
		scrollEdgeStyle := pane.scrollEdgeEffectStyle
		hasScrollEdgeStyle := pane.scrollEdgeEffectStyleSet
		pane.lock.RUnlock()
		snapshot.contentLayout = resolveMacContentLayout(w.parent.options.Mac, paneLayout)
		C.splitViewAddPane(handle,
			C.ulonglong(pane.internalID),
			C.int(snapshot.role),
			C.bool(pane.primary),
			C.double(snapshot.minimumThickness),
			C.double(snapshot.maximumThickness),
			C.double(snapshot.preferredThickness),
			C.bool(snapshot.preferredThicknessSet),
			C.double(snapshot.holdingPriority),
			C.bool(snapshot.holdingPrioritySet),
			C.bool(snapshot.collapsible),
			C.bool(snapshot.collapsibleSet),
			C.bool(snapshot.canCollapseFromResize),
			C.bool(snapshot.canCollapseFromResizeSet),
			C.bool(snapshot.collapsed),
			C.int(snapshot.contentLayout),
			C.int(scrollEdgeStyle),
			C.bool(hasScrollEdgeStyle))
		if pane.sidebar != nil {
			pane.sidebar.registerItems()
			applyMacSidebarSnapshotToNative(handle, pane.internalID, pane.sidebar.snapshot())
		}
		if pane.inspector != nil {
			pane.inspector.registerControls()
			applyMacInspectorSnapshotToNative(handle, pane.internalID, pane.inspector.snapshot())
		}

		// Registered before installation so no native callback can arrive for
		// an unknown pane.
		registerMacSplitPane(pane)
	}

	installed := C.splitViewInstall(handle, w.nsWindow,
		C.bool(w.parent.options.Mac.Backdrop == MacBackdropNormal))
	if !bool(installed) {
		for _, pane := range panes {
			unregisterMacSplitPane(pane.internalID)
			if pane.sidebar != nil {
				for _, item := range pane.sidebar.itemHandles() {
					unregisterMacSidebarItem(item.internalID)
				}
			}
			if pane.inspector != nil {
				for _, control := range pane.inspector.controlHandles() {
					unregisterMacInspectorControl(control.internalID)
				}
			}
		}
		C.splitViewRelease(handle)
		w.parent.Error("SetSplitView: failed to install the native split view")
		return
	}

	split.lock.Lock()
	split.native = unsafe.Pointer(handle)
	split.installed = true
	split.lock.Unlock()
	w.activeSplitView = split

	// A pane setter may have run between the snapshot above and the native
	// commit; reapply the latest state so no update is lost.
	for _, pane := range panes {
		applyMacSplitPaneLatestState(pane)
	}
}

// teardownSplitView detaches the Go callback registry, invalidates every pane
// handle, and releases the native split resources. It runs on the application
// thread before the native window is destroyed and is safe to call more than
// once.
func (w *macosWebviewWindow) teardownSplitView() {
	split := w.activeSplitView
	w.activeSplitView = nil
	w.parent.splitViewLock.Lock()
	if split == nil {
		split = w.parent.splitView
	}
	w.parent.splitView = nil
	w.parent.splitViewLock.Unlock()
	if split == nil {
		return
	}

	split.lock.Lock()
	native := split.native
	split.native = nil
	split.installed = false
	split.lock.Unlock()

	// Stop routing native callbacks before releasing native state, and mark
	// every handle dead so later calls are safe no-ops.
	for _, pane := range split.paneSnapshot() {
		unregisterMacSplitPane(pane.internalID)
		pane.markDead()
	}
	if native != nil {
		C.splitViewTeardown(native)
		C.splitViewRelease(native)
	}
}

// macSplitPaneWithNative runs a native update for an installed, live pane on
// the application thread. It re-checks installation under the split lock on
// that thread so a pane cannot be torn down between the check and the call.
func macSplitPaneWithNative(p *MacSplitPane, apply func(handle unsafe.Pointer, paneID C.ulonglong)) {
	if p == nil || p.split == nil {
		return
	}
	split := p.split
	split.lock.RLock()
	installed := split.installed && split.native != nil
	split.lock.RUnlock()
	if !installed {
		return
	}
	InvokeSync(func() {
		split.lock.RLock()
		defer split.lock.RUnlock()
		if split.native == nil || p.isDead() {
			return
		}
		apply(split.native, C.ulonglong(p.internalID))
	})
}

// macSplitPaneSnapshot is a lock-free copy of one pane's configuration.
type macSplitPaneSnapshot struct {
	role                     macSplitPaneRole
	minimumThickness         float64
	maximumThickness         float64
	preferredThickness       float64
	preferredThicknessSet    bool
	holdingPriority          float64
	holdingPrioritySet       bool
	collapsible              bool
	collapsibleSet           bool
	canCollapseFromResize    bool
	canCollapseFromResizeSet bool
	collapsed                bool
	contentLayout            MacContentLayout
}

func snapshotMacSplitPane(p *MacSplitPane) macSplitPaneSnapshot {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return macSplitPaneSnapshot{
		role:                     p.role,
		minimumThickness:         p.minimumThickness,
		maximumThickness:         p.maximumThickness,
		preferredThickness:       p.preferredThickness,
		preferredThicknessSet:    p.preferredThicknessSet,
		holdingPriority:          p.holdingPriority,
		holdingPrioritySet:       p.holdingPrioritySet,
		collapsible:              p.collapsible,
		collapsibleSet:           p.collapsibleSet,
		canCollapseFromResize:    p.canCollapseFromResize,
		canCollapseFromResizeSet: p.canCollapseFromResizeSet,
		collapsed:                p.collapsed,
	}
}

// applyMacSplitPaneLatestState replays a pane's explicitly configured state
// onto its freshly installed native item.
func applyMacSplitPaneLatestState(pane *MacSplitWebviewPane) {
	snapshot := snapshotMacSplitPane(pane.MacSplitPane)
	if snapshot.minimumThickness > 0 {
		macSplitPaneApplyMinimumThickness(pane.MacSplitPane, snapshot.minimumThickness)
	}
	if snapshot.maximumThickness > 0 {
		macSplitPaneApplyMaximumThickness(pane.MacSplitPane, snapshot.maximumThickness)
	}
	if snapshot.preferredThicknessSet && snapshot.preferredThickness > 0 && snapshot.preferredThickness <= 1 {
		macSplitPaneApplyPreferredFraction(pane.MacSplitPane, snapshot.preferredThickness)
	}
	if snapshot.holdingPrioritySet && snapshot.holdingPriority >= 1 && snapshot.holdingPriority <= 1000 {
		macSplitPaneApplyHoldingPriority(pane.MacSplitPane, snapshot.holdingPriority)
	}
	if snapshot.collapsibleSet {
		macSplitPaneApplyCollapsible(pane.MacSplitPane, snapshot.collapsible)
	}
	if snapshot.canCollapseFromResizeSet {
		macSplitPaneApplyCanCollapseFromResize(pane.MacSplitPane, snapshot.canCollapseFromResize)
	}
	pane.lock.RLock()
	layout := pane.contentLayout
	scrollEdgeStyle := pane.scrollEdgeEffectStyle
	hasScrollEdgeStyle := pane.scrollEdgeEffectStyleSet
	pane.lock.RUnlock()
	if owner := pane.split.ownerWindow(); owner != nil {
		macSplitPaneApplyContentLayout(pane.MacSplitPane, resolveMacContentLayout(owner.macSplitOptions(), layout))
	}
	if hasScrollEdgeStyle {
		macSplitPaneApplyScrollEdgeEffectStyle(pane.MacSplitPane, scrollEdgeStyle)
	}
}

func macSplitPaneApplyMinimumThickness(p *MacSplitPane, value float64) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetMinimumThickness(handle, paneID, C.double(value))
	})
}

func macSplitPaneApplyMaximumThickness(p *MacSplitPane, value float64) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetMaximumThickness(handle, paneID, C.double(value))
	})
}

func macSplitPaneApplyPreferredFraction(p *MacSplitPane, value float64) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetPreferredFraction(handle, paneID, C.double(value))
	})
}

func macSplitPaneApplyHoldingPriority(p *MacSplitPane, value float64) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetHoldingPriority(handle, paneID, C.double(value))
	})
}

func macSplitPaneApplyCollapsible(p *MacSplitPane, collapsible bool) {
	p.lock.RLock()
	resizeAllowed := p.canCollapseFromResize
	resizeConfigured := p.canCollapseFromResizeSet
	p.lock.RUnlock()
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetCollapsible(handle, paneID, C.bool(collapsible))
		if resizeConfigured {
			C.splitViewPaneSetCanCollapseFromWindowResize(handle, paneID, C.bool(resizeAllowed))
		}
	})
}

func macSplitPaneApplyCanCollapseFromResize(p *MacSplitPane, allowed bool) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetCanCollapseFromWindowResize(handle, paneID, C.bool(allowed))
	})
}

func macSplitPaneApplyContentLayout(p *MacSplitPane, layout MacContentLayout) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetContentLayout(handle, paneID, C.int(layout))
	})
}

func macSplitPaneApplyScrollEdgeEffectStyle(p *MacSplitPane, style MacScrollEdgeEffectStyle) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetScrollEdgeEffectStyle(handle, paneID, C.int(style))
	})
}

func macSplitPaneApplyCollapsed(p *MacSplitPane, collapsed bool) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneSetCollapsed(handle, paneID, C.bool(collapsed))
	})
}

func macSplitPaneApplyToggle(p *MacSplitPane) {
	macSplitPaneWithNative(p, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewPaneToggleCollapsed(handle, paneID)
	})
}
