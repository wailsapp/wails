//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "native_window_darwin.h"
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type macosNativeWindow struct {
	parent        *NativeWindow
	nsWindow      unsafe.Pointer
	activeToolbar *MacToolbar
	activeSplit   *MacSplitView
}

var nativeWindowClosed = make(chan uint, 16)

//export processNativeWindowClosed
func processNativeWindowClosed(windowID C.uint) {
	nativeWindowClosed <- uint(windowID)
}

func newNativeWindowImpl(parent *NativeWindow) nativeWindowImpl {
	return &macosNativeWindow{parent: parent}
}

func (w *macosNativeWindow) run() error {
	globalApplication.dispatchOnMainThread(func() {
		options := w.parent.options
		w.nsWindow = C.nativeWindowCreate(C.uint(w.parent.id), C.int(options.Width), C.int(options.Height), C.bool(options.HideOnClose))
		if w.nsWindow == nil {
			w.parent.Error("failed to create NSWindow")
			return
		}
		w.setTitle(options.Title)
		C.nativeWindowSetResizable(w.nsWindow, C.bool(!options.DisableResize))
		if options.MinWidth != 0 || options.MinHeight != 0 {
			C.nativeWindowSetMinSize(w.nsWindow, C.int(options.MinWidth), C.int(options.MinHeight))
		}
		if options.MaxWidth != 0 || options.MaxHeight != 0 {
			C.nativeWindowSetMaxSize(w.nsWindow, C.int(options.MaxWidth), C.int(options.MaxHeight))
		}
		C.nativeWindowSetAlwaysOnTop(w.nsWindow, C.bool(options.AlwaysOnTop))
		titlebar := options.Mac.TitleBar
		C.nativeWindowConfigureTitlebar(w.nsWindow,
			C.bool(titlebar.AppearsTransparent),
			C.bool(titlebar.FullSizeContent),
			C.bool(titlebar.HideTitle),
			C.bool(titlebar.HideToolbarSeparator),
			C.int(titlebar.ToolbarStyle))
		if err := w.installSplitView(); err != nil {
			w.parent.Error("%s", err)
			w.close()
			return
		}
		w.parent.lock.RLock()
		toolbar := w.parent.toolbar
		w.parent.lock.RUnlock()
		if toolbar != nil {
			if err := w.setToolbar(toolbar); err != nil {
				w.parent.Error("SetToolbar: %s", err)
				w.close()
				return
			}
		}
		if options.InitialPosition == WindowCentered {
			C.nativeWindowCenter(w.nsWindow)
		} else {
			C.nativeWindowSetPosition(w.nsWindow, C.int(options.X), C.int(options.Y))
		}
		if !options.Hidden {
			w.show()
			// A focus request made while constructing the hierarchy cannot take
			// effect before the NSWindow exists. Make the native primary editor
			// first responder after the window is key so typing works immediately.
			for _, pane := range w.activeSplit.paneSnapshot() {
				if pane.editor != nil {
					macTextEditorFocus(pane.editor)
					break
				}
			}
		}
	})
	return nil
}

func (w *macosNativeWindow) installSplitView() error {
	w.parent.lock.RLock()
	split := w.parent.split
	w.parent.lock.RUnlock()
	if split == nil {
		return fmt.Errorf("NativeWindow requires a MacSplitView content layout")
	}

	split.lock.RLock()
	autosaveName := split.autosaveName
	split.lock.RUnlock()
	autosaveC := C.CString(autosaveName)
	handle := C.splitViewCreate(autosaveC)
	C.free(unsafe.Pointer(autosaveC))
	if handle == nil {
		return fmt.Errorf("failed to create native split view")
	}

	panes := split.paneSnapshot()
	configuredTextVersions := make(map[uint64]uint64)
	committed := false
	defer func() {
		if !committed {
			C.splitViewRelease(handle)
		}
	}()
	for _, pane := range panes {
		snapshot := snapshotMacSplitPane(pane.MacSplitPane)
		C.splitViewAddPane(handle,
			C.ulonglong(pane.internalID), C.int(snapshot.role), C.bool(pane.primary),
			C.double(snapshot.minimumThickness), C.double(snapshot.maximumThickness),
			C.double(snapshot.preferredThickness), C.bool(snapshot.preferredThicknessSet),
			C.double(snapshot.holdingPriority), C.bool(snapshot.holdingPrioritySet),
			C.bool(snapshot.collapsible), C.bool(snapshot.collapsibleSet),
			C.bool(snapshot.canCollapseFromResize), C.bool(snapshot.canCollapseFromResizeSet),
			C.bool(snapshot.collapsed), C.int(MacContentLayoutEdgeToEdge))

		if pane.sidebar != nil {
			pane.sidebar.registerItems()
			applyMacSidebarSnapshotToNative(unsafe.Pointer(handle), pane.internalID, pane.sidebar.snapshot())
		}
		if pane.inspector != nil {
			pane.inspector.registerControls()
			applyMacInspectorSnapshotToNative(unsafe.Pointer(handle), pane.internalID, pane.inspector.snapshot())
		}
		if pane.editor != nil {
			editorID, text, editable, version := pane.editor.snapshot()
			textC := C.CString(text)
			C.splitViewConfigureTextEditor(handle, C.ulonglong(pane.internalID), C.ulonglong(editorID), textC, C.bool(editable))
			C.free(unsafe.Pointer(textC))
			configuredTextVersions[pane.internalID] = version
			registerMacTextEditor(pane.editor)
		}
		registerMacSplitPane(pane)
	}
	if !bool(C.splitViewInstallNative(handle, w.nsWindow,
		C.bool(w.parent.options.Mac.Backdrop == MacBackdropNormal))) {
		for _, pane := range panes {
			unregisterMacSplitPane(pane.internalID)
			if pane.editor != nil {
				unregisterMacTextEditor(pane.editor.internalID)
			}
		}
		return fmt.Errorf("failed to install native split view; NativeWindow requires exactly one AddTextEditor primary pane")
	}
	split.lock.Lock()
	split.native = unsafe.Pointer(handle)
	split.installed = true
	split.lock.Unlock()
	w.activeSplit = split
	committed = true
	for _, pane := range panes {
		applyMacSplitPaneLatestState(pane)
		if pane.editor != nil {
			_, text, editable, version := pane.editor.snapshot()
			// Only replay text when SetText raced with native installation.
			// Replacing an NSTextView's large backing store with identical text
			// needlessly raises both peak and retained allocator capacity.
			if configuredTextVersions[pane.internalID] != version {
				textC := C.CString(text)
				C.splitViewTextEditorSetText(handle, C.ulonglong(pane.internalID), textC)
				C.free(unsafe.Pointer(textC))
			}
			C.splitViewTextEditorSetEditable(handle, C.ulonglong(pane.internalID), C.bool(editable))
			pane.editor.clearCachedText(version)
		}
	}
	return nil
}

func (w *macosNativeWindow) setToolbar(toolbar *MacToolbar) error {
	hasSidebar := w.activeSplit != nil && w.activeSplit.hasSidebarPane()
	hasInspector := w.activeSplit != nil && w.activeSplit.hasInspectorPane()
	attached, err := attachMacToolbar(w.parent, w.nsWindow, w.activeToolbar, toolbar,
		hasSidebar, hasInspector, w.parent.options.Mac.TitleBar)
	if err == nil {
		w.activeToolbar = attached
	}
	return err
}

func (w *macosNativeWindow) teardownSplitView() {
	split := w.activeSplit
	w.activeSplit = nil
	if split == nil {
		return
	}
	split.lock.Lock()
	native := split.native
	split.native = nil
	split.installed = false
	split.lock.Unlock()
	for _, pane := range split.paneSnapshot() {
		unregisterMacSplitPane(pane.internalID)
		pane.markDead()
	}
	if native != nil {
		C.splitViewTeardown(native)
		C.splitViewRelease(native)
	}
}

func (w *macosNativeWindow) show()           { C.nativeWindowShow(w.nsWindow) }
func (w *macosNativeWindow) hide()           { C.nativeWindowHide(w.nsWindow) }
func (w *macosNativeWindow) focus()          { C.nativeWindowFocus(w.nsWindow) }
func (w *macosNativeWindow) isVisible() bool { return bool(C.nativeWindowIsVisible(w.nsWindow)) }
func (w *macosNativeWindow) setTitle(title string) {
	titleC := C.CString(title)
	C.nativeWindowSetTitle(w.nsWindow, titleC)
	C.free(unsafe.Pointer(titleC))
}
func (w *macosNativeWindow) nativeWindow() unsafe.Pointer { return w.nsWindow }
func (w *macosNativeWindow) close() {
	if w.nsWindow == nil {
		return
	}
	if w.activeToolbar != nil {
		_, _ = attachMacToolbar(w.parent, w.nsWindow, w.activeToolbar, nil, false, false, w.parent.options.Mac.TitleBar)
		w.activeToolbar = nil
	}
	w.teardownSplitView()
	C.nativeWindowDestroy(w.nsWindow)
	w.nsWindow = nil
}
