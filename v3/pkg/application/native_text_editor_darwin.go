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

func macTextEditorApplyText(editor *MacTextEditor, text string) bool {
	if editor == nil {
		return false
	}
	editor.lock.RLock()
	pane := editor.pane
	editor.lock.RUnlock()
	if pane == nil {
		return false
	}
	applied := false
	textC := C.CString(text)
	defer C.free(unsafe.Pointer(textC))
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewTextEditorSetText(handle, paneID, textC)
		applied = true
	})
	return applied
}

func macTextEditorReadText(editor *MacTextEditor) (string, bool) {
	if editor == nil {
		return "", false
	}
	editor.lock.RLock()
	pane := editor.pane
	editor.lock.RUnlock()
	if pane == nil || pane.split == nil || !pane.split.isInstalled() {
		return "", false
	}
	var result string
	var read bool
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		value := C.splitViewTextEditorCopyText(handle, paneID)
		if value == nil {
			return
		}
		result = C.GoString(value)
		C.free(unsafe.Pointer(value))
		read = true
	})
	return result, read
}

func macTextEditorApplyEditable(editor *MacTextEditor, editable bool) {
	if editor == nil {
		return
	}
	editor.lock.RLock()
	pane := editor.pane
	editor.lock.RUnlock()
	if pane == nil {
		return
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewTextEditorSetEditable(handle, paneID, C.bool(editable))
	})
}

func macTextEditorFocus(editor *MacTextEditor) {
	if editor == nil {
		return
	}
	editor.lock.RLock()
	pane := editor.pane
	editor.lock.RUnlock()
	if pane == nil {
		return
	}
	macSplitPaneWithNative(pane, func(handle unsafe.Pointer, paneID C.ulonglong) {
		C.splitViewTextEditorFocus(handle, paneID)
	})
}
