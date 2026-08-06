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
	"encoding/json"
	"unsafe"
)

func macInspectorRegisterControlIfInstalled(control *MacInspectorControl) {
	if control == nil || control.inspector == nil {
		return
	}
	control.inspector.lock.RLock()
	pane := control.inspector.pane
	control.inspector.lock.RUnlock()
	if pane != nil && pane.split.isInstalled() {
		registerMacInspectorControl(control)
	}
}

func macInspectorApplySnapshot(inspector *MacInspector) {
	if inspector == nil {
		return
	}
	snapshot := inspector.snapshot()
	inspector.lock.RLock()
	pane := inspector.pane
	inspector.lock.RUnlock()
	if pane == nil {
		return
	}
	if pane.split.isInstalled() {
		inspector.registerControls()
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		applyMacInspectorSnapshotToNative(handle, uint64(paneID), snapshot)
	})
}

func macInspectorApplyControl(control *MacInspectorControl) {
	if control == nil || control.inspector == nil {
		return
	}
	control.inspector.lock.RLock()
	pane := control.inspector.pane
	control.inspector.lock.RUnlock()
	if pane == nil {
		return
	}
	snapshot := snapshotMacInspectorControl(control)
	options, _ := json.Marshal(snapshot.options)
	labelC := C.CString(snapshot.label)
	valueC := C.CString(snapshot.value)
	optionsC := C.CString(string(options))
	tooltipC := C.CString(snapshot.tooltip)
	defer C.free(unsafe.Pointer(labelC))
	defer C.free(unsafe.Pointer(valueC))
	defer C.free(unsafe.Pointer(optionsC))
	defer C.free(unsafe.Pointer(tooltipC))
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewInspectorUpdateControl(handle, paneID, C.ulonglong(snapshot.internalID),
			C.int(snapshot.kind), labelC, valueC, C.bool(snapshot.checked), optionsC,
			C.int(snapshot.selected), tooltipC, C.bool(snapshot.disabled), C.bool(snapshot.hidden))
	})
}

func applyMacInspectorSnapshotToNative(handle unsafe.Pointer, paneID uint64, snapshot macInspectorSnapshot) {
	C.splitViewInspectorReset(handle, C.ulonglong(paneID))
	for _, section := range snapshot.sections {
		labelC := C.CString(section.label)
		C.splitViewInspectorAddSection(handle, C.ulonglong(paneID), C.ulonglong(section.internalID), labelC)
		C.free(unsafe.Pointer(labelC))
		for _, control := range section.controls {
			addMacInspectorControlToNative(handle, paneID, section.internalID, control)
		}
	}
	C.splitViewInspectorReload(handle, C.ulonglong(paneID))
}

func addMacInspectorControlToNative(handle unsafe.Pointer, paneID, sectionID uint64, control macInspectorControlSnapshot) {
	options, _ := json.Marshal(control.options)
	labelC := C.CString(control.label)
	valueC := C.CString(control.value)
	optionsC := C.CString(string(options))
	tooltipC := C.CString(control.tooltip)
	C.splitViewInspectorAddControl(handle, C.ulonglong(paneID), C.ulonglong(sectionID),
		C.ulonglong(control.internalID), C.int(control.kind), labelC, valueC,
		C.bool(control.checked), optionsC, C.int(control.selected), tooltipC,
		C.bool(control.disabled), C.bool(control.hidden))
	C.free(unsafe.Pointer(labelC))
	C.free(unsafe.Pointer(valueC))
	C.free(unsafe.Pointer(optionsC))
	C.free(unsafe.Pointer(tooltipC))
}
