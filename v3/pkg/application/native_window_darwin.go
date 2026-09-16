//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "native_window_darwin.h"
#include "mac_window_chrome_darwin.h"
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// NativeWindowOptions.Mac contract on macOS.
//
// A NativeWindow applies the MacWindow options below with the same AppKit
// semantics WebviewWindow uses; the appliers live in mac_window_chrome_darwin.m
// and native_window_darwin.m. NativeWindowOptions has no Frameless field, so
// TitleBar.Hide is the native frameless switch and follows the
// WebviewWindowOptions.Frameless rules on macOS.
//
//	Field                         Native window behaviour
//	Backdrop                      Normal: opaque windowBackgroundColor surface.
//	                              Transparent: clear window, transparent editor.
//	                              Translucent: NSVisualEffectView (behind-window)
//	                              beneath every pane, transparent editor.
//	                              LiquidGlass: NSGlassEffectView on macOS 26+
//	                              beneath every pane (translucent fallback
//	                              elsewhere), transparent editor. Sidebar,
//	                              inspector and content-list panes keep their
//	                              own material in every mode.
//	LiquidGlass                   Style, Material, CornerRadius, TintColor apply;
//	                              GroupID/GroupSpacing need -tags private_mac_apis
//	                              exactly as for WebviewWindow.
//	DisableShadow                 NSWindow.hasShadow.
//	TitleBar.Hide                 Frameless: with CornerType Rounded and
//	                              CornerRadius 0 the AppKit frame is kept with a
//	                              transparent titlebar, hidden title and hidden
//	                              window buttons; otherwise the window is
//	                              borderless. The other TitleBar fields are
//	                              ignored, as for a frameless WebviewWindow.
//	CornerType / CornerRadius     Only with TitleBar.Hide: Square gives a
//	                              borderless square window; a radius gives a
//	                              borderless window whose content is masked to
//	                              that radius.
//	TitleBar.AppearsTransparent   NSWindow.titlebarAppearsTransparent.
//	TitleBar.HideTitle            NSWindow.titleVisibility.
//	TitleBar.FullSizeContent      NSWindowStyleMaskFullSizeContentView (also the
//	                              ContentLayout Automatic input).
//	TitleBar.UseToolbar           Ignored: a NativeWindow only shows a toolbar
//	                              attached with SetToolbar.
//	TitleBar.ToolbarStyle         NSWindow.toolbarStyle.
//	TitleBar.HideToolbarSeparator NSToolbar.showsBaselineSeparator.
//	TitleBar.ShowToolbarWhenFullscreen  Fullscreen presentation options.
//	ContentLayout                 EdgeToEdge keeps the split layout full-size so
//	                              the editor scrolls beneath the toolbar with
//	                              AppKit's automatic insets. BelowToolbar removes
//	                              full-size content, which places the whole
//	                              layout, sidebar included, beneath the toolbar
//	                              (the split view is the window's content view).
//	                              Automatic follows TitleBar.FullSizeContent.
//	Appearance                    NSWindow.appearance by name.
//	WindowLevel                   NSWindow.level; precedence WindowLevel, then
//	                              AlwaysOnTop or a floating panel, then normal.
//	CollectionBehavior            NSWindow.collectionBehavior; zero selects
//	                              FullScreenPrimary.
//	TabbingMode                   NSWindow.tabbingMode; Default resolves to
//	                              Disallowed.
//	DisableEscapeExitsFullscreen  cancelOperation: is swallowed in fullscreen.
//	WindowClass / PanelPreferences  Panel creates an NSPanel subclass with the
//	                              FloatingPanel, BecomesKeyOnlyIfNeeded,
//	                              NonActivating and UtilityWindow preferences.
//	InvisibleTitleBarHeight       Not applicable: it sizes the WebView drag
//	                              region, and native views handle their own drag.
//	EventMapping                  Not applicable: NativeWindow has no window
//	                              event API yet, so nothing is emitted to remap.
//	EnableFraudulentWebsiteWarnings, WebviewPreferences  Not applicable: no WKWebView.
//	SplitView, Toolbar            Ignored: NativeWindowOptions.SplitView and
//	                              NativeWindowOptions.Toolbar are used instead.
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
		macOptions := options.Mac
		frame := resolveNativeMacFrame(macOptions)
		w.nsWindow = C.nativeWindowCreate(w.windowConfig(frame))
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
		// The order below matches macosWebviewWindow.run so both window kinds
		// resolve the same MacWindow options identically.
		C.windowChromeSetShadow(w.nsWindow, C.bool(!macOptions.DisableShadow))
		C.windowChromeSetLevel(w.nsWindow, C.int(nativeMacWindowLevelCode(effectiveNativeMacWindowLevel(options))))
		C.windowChromeSetCollectionBehavior(w.nsWindow, C.int(macOptions.CollectionBehavior))
		C.windowChromeSetTabbingMode(w.nsWindow, C.int(nativeMacTabbingModeValue(macOptions.TabbingMode)))
		titlebar := macOptions.TitleBar
		if !frame.frameless {
			C.nativeWindowConfigureTitlebar(w.nsWindow,
				C.bool(titlebar.AppearsTransparent),
				C.bool(titlebar.FullSizeContent),
				C.bool(titlebar.HideTitle),
				C.bool(titlebar.HideToolbarSeparator),
				C.int(titlebar.ToolbarStyle))
		}
		if err := w.installSplitView(); err != nil {
			w.parent.Error("%s", err)
			w.close()
			return
		}
		w.applyContentChrome(frame)
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
		if macOptions.Appearance != "" {
			appearance := C.CString(string(macOptions.Appearance))
			C.windowChromeSetAppearanceByName(w.nsWindow, appearance)
			C.free(unsafe.Pointer(appearance))
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

// windowConfig gathers the options AppKit needs at construction time.
func (w *macosNativeWindow) windowConfig(frame nativeMacFrame) C.WailsNativeWindowConfig {
	options := w.parent.options
	panel := options.Mac.PanelPreferences
	return C.WailsNativeWindowConfig{
		windowID:                     C.uint(w.parent.id),
		width:                        C.int(options.Width),
		height:                       C.int(options.Height),
		hideOnClose:                  C.bool(options.HideOnClose),
		frameless:                    C.bool(frame.frameless),
		borderless:                   C.bool(frame.borderless),
		cornerRadius:                 C.double(frame.cornerRadius),
		isPanel:                      C.bool(options.Mac.WindowClass == MacWindowClassPanel),
		floatingPanel:                C.bool(panel.FloatingPanel),
		becomesKeyOnlyIfNeeded:       C.bool(panel.BecomesKeyOnlyIfNeeded),
		nonActivating:                C.bool(panel.NonActivating),
		utilityWindow:                C.bool(panel.UtilityWindow),
		disableEscapeExitsFullscreen: C.bool(options.Mac.DisableEscapeExitsFullscreen),
	}
}

// applyContentChrome applies ContentLayout, the custom corner mask and the
// Backdrop once the split layout owns the window's content view. The backdrop
// is created here rather than before installation because the split installer
// replaces the content view, which would discard an earlier backdrop.
func (w *macosNativeWindow) applyContentChrome(frame nativeMacFrame) {
	macOptions := w.parent.options.Mac
	glass := macOptions.LiquidGlass
	backdrop := C.WailsNativeBackdropConfig{
		backdrop:          C.int(macOptions.Backdrop),
		glassStyle:        C.int(glass.Style),
		glassMaterial:     C.int(glass.Material),
		glassCornerRadius: C.double(nativeMacLiquidGlassCornerRadius(glass)),
		groupSpacing:      C.double(glass.GroupSpacing),
	}
	if glass.TintColor != nil {
		backdrop.tintR = C.int(glass.TintColor.Red)
		backdrop.tintG = C.int(glass.TintColor.Green)
		backdrop.tintB = C.int(glass.TintColor.Blue)
		backdrop.tintA = C.int(glass.TintColor.Alpha)
	}
	if glass.GroupID != "" {
		backdrop.groupID = C.CString(glass.GroupID)
		defer C.free(unsafe.Pointer(backdrop.groupID))
	}
	if macOptions.Backdrop == MacBackdropLiquidGlass && !bool(C.windowChromeLiquidGlassSupported()) {
		globalApplication.debug("Liquid Glass not supported on this macOS version, falling back to translucent", "window", w.parent.id)
	}
	belowToolbar := resolveNativeMacContentLayout(macOptions) == MacContentLayoutBelowToolbar
	C.nativeWindowApplyContentChrome(w.nsWindow, C.bool(belowToolbar), C.double(frame.cornerRadius), backdrop)
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
		if pane.contentList != nil {
			pane.contentList.registerRows()
			applyMacContentListSnapshotToNative(unsafe.Pointer(handle), pane.internalID, pane.contentList.snapshot())
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
