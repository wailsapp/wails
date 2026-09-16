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
	"time"
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

func macInspectorRegisterSectionIfInstalled(section *MacInspectorSection) {
	if section == nil || section.inspector == nil {
		return
	}
	section.inspector.lock.RLock()
	pane := section.inspector.pane
	section.inspector.lock.RUnlock()
	if pane != nil && pane.split.isInstalled() {
		registerMacInspectorSection(section)
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
	spec, release := newMacInspectorControlSpec(snapshot)
	defer release()
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewInspectorUpdateControl(handle, paneID, C.ulonglong(snapshot.internalID), spec)
	})
}

func applyMacInspectorSnapshotToNative(handle unsafe.Pointer, paneID uint64, snapshot macInspectorSnapshot) {
	C.splitViewInspectorReset(handle, C.ulonglong(paneID))
	for _, section := range snapshot.sections {
		labelC := C.CString(section.label)
		C.splitViewInspectorAddSection(handle, C.ulonglong(paneID), C.ulonglong(section.internalID), labelC,
			C.bool(section.collapsible), C.bool(section.collapsed))
		C.free(unsafe.Pointer(labelC))
		for _, control := range section.controls {
			addMacInspectorControlToNative(handle, paneID, section.internalID, control)
		}
	}
	C.splitViewInspectorReload(handle, C.ulonglong(paneID))
}

func addMacInspectorControlToNative(handle unsafe.Pointer, paneID, sectionID uint64, control macInspectorControlSnapshot) {
	spec, release := newMacInspectorControlSpec(control)
	defer release()
	C.splitViewInspectorAddControl(handle, C.ulonglong(paneID), C.ulonglong(sectionID),
		C.ulonglong(control.internalID), spec)
}

// newMacInspectorControlSpec marshals a control snapshot for AppKit. The
// returned release function frees the C strings it allocates.
func newMacInspectorControlSpec(control macInspectorControlSnapshot) (C.WailsInspectorControlSpec, func()) {
	options, _ := json.Marshal(control.options)
	labelC := C.CString(control.label)
	valueC := C.CString(control.value)
	optionsC := C.CString(string(options))
	tooltipC := C.CString(control.tooltip)
	var dateSeconds float64
	if !control.date.IsZero() {
		dateSeconds = float64(control.date.UnixNano()) / float64(time.Second)
	}
	spec := C.WailsInspectorControlSpec{
		kind:          C.int(control.kind),
		label:         labelC,
		value:         valueC,
		checked:       C.bool(control.checked),
		optionsJSON:   optionsC,
		selectedIndex: C.int(control.selected),
		tooltip:       tooltipC,
		disabled:      C.bool(control.disabled),
		hidden:        C.bool(control.hidden),
		number:        C.double(control.number),
		minimum:       C.double(control.minimum),
		maximum:       C.double(control.maximum),
		step:          C.double(control.step),
		red:           C.int(control.colour.Red),
		green:         C.int(control.colour.Green),
		blue:          C.int(control.colour.Blue),
		alpha:         C.int(control.colour.Alpha),
		dateSeconds:   C.double(dateSeconds),
	}
	return spec, func() {
		C.free(unsafe.Pointer(labelC))
		C.free(unsafe.Pointer(valueC))
		C.free(unsafe.Pointer(optionsC))
		C.free(unsafe.Pointer(tooltipC))
	}
}

//export processMacInspectorNumberChanged
func processMacInspectorNumberChanged(controlID C.ulonglong, kind C.int, value C.double) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorControlKind(kind),
		number:    float64(value),
	}
}

//export processMacInspectorSegmentChanged
func processMacInspectorSegmentChanged(controlID C.ulonglong, selectedIndex C.int) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorSegmented,
		selected:  int(selectedIndex),
	}
}

//export processMacInspectorColorChanged
func processMacInspectorColorChanged(controlID C.ulonglong, red, green, blue, alpha C.int) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorColorWell,
		colour:    RGBA{Red: uint8(red), Green: uint8(green), Blue: uint8(blue), Alpha: uint8(alpha)},
	}
}

//export processMacInspectorDateChanged
func processMacInspectorDateChanged(controlID C.ulonglong, unixSeconds C.double) {
	seconds := float64(unixSeconds)
	whole := int64(seconds)
	nanos := int64((seconds - float64(whole)) * float64(time.Second))
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorDatePicker,
		date:      time.Unix(whole, nanos),
	}
}

//export processMacInspectorButtonClicked
func processMacInspectorButtonClicked(controlID C.ulonglong) {
	macInspectorControlEvents <- macInspectorControlEvent{
		controlID: uint64(controlID),
		kind:      MacInspectorButton,
	}
}

//export processMacInspectorSectionCollapsed
func processMacInspectorSectionCollapsed(sectionID C.ulonglong, collapsed C.bool) {
	macInspectorSectionEvents <- macInspectorSectionEvent{sectionID: uint64(sectionID), collapsed: bool(collapsed)}
}
