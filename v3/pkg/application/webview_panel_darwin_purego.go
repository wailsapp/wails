//go:build darwin && purego && !ios && !server

package application

// CGO-free port of WebviewPanel (webview_panel_darwin.m): the NSPanel subclass
// used for MacWindowClassPanel windows. The panel shares the NSWindow method
// set (webviewWindowMethods) and overrides only what cgo's WebviewPanel
// overrides: canBecomeMainWindow (NO) and sendEvent: (key bindings first).

import (
	"sync"

	"github.com/ebitengine/purego/objc"
)

// AppKit constants (NSWindow.h, NSEvent.h of the macOS 26 SDK).
const (
	nsWindowStyleMaskUtilityWindow      = 1 << 4 // NSWindowStyleMaskUtilityWindow
	nsWindowStyleMaskNonactivatingPanel = 1 << 7 // NSWindowStyleMaskNonactivatingPanel
	nsEventTypeKeyDown                  = 10     // NSEventTypeKeyDown
)

// panelStyleMask returns the style mask of a panel window: the base mask from
// the window path plus the flags selected by the panel preferences, exactly as
// cgo createNativeWindow ORs them in. Bits outside the two panel flags are
// left untouched.
func panelStyleMask(base uint, prefs MacPanelPreferences) uint {
	mask := base
	if prefs.NonActivating {
		mask |= nsWindowStyleMaskNonactivatingPanel
	}
	if prefs.UtilityWindow {
		mask |= nsWindowStyleMaskUtilityWindow
	}
	return mask
}

// isNonActivatingPanel reports whether win is an NSPanel with the
// non-activating style. cgo show/focus use orderFrontRegardless for these so
// showing the panel does not activate the application.
func isNonActivatingPanel(win id) bool {
	if !get[bool](win, "isKindOfClass:", class("NSPanel")) {
		return false
	}
	return get[uint](win, "styleMask")&nsWindowStyleMaskNonactivatingPanel != 0
}

// configureNativePanel applies the NSPanel properties cgo createNativeWindow
// and WebviewPanel init set. alphaValue, clearColor, opaque and
// movableByWindowBackground are already applied to every window by createWindow.
func configureNativePanel(win id, prefs MacPanelPreferences) {
	win.send("setFloatingPanel:", prefs.FloatingPanel)
	win.send("setBecomesKeyOnlyIfNeeded:", prefs.BecomesKeyOnlyIfNeeded)
	// NSPanel defaults differ from NSWindow: keep Wails windows visible when
	// another app activates, and release on Close as cgo does.
	win.send("setHidesOnDeactivate:", false)
	win.send("setReleasedWhenClosed:", true)
}

// panelHandleKeyDown offers a key-down event to the Wails key bindings before
// AppKit delivers it (cgo WebviewPanel sendEvent:). It reports whether the
// event was consumed. A non-activating panel whose first responder is the
// WKWebView would otherwise never reach keyDown: or performKeyEquivalent:.
func panelHandleKeyDown(panel, event id) bool {
	const nsEventTypeKeyDown = 10
	if get[uint](event, "type") != nsEventTypeKeyDown {
		return false
	}
	v, ok := nsWindowToID.Load(uintptr(panel))
	if !ok {
		return false
	}
	accelerator := acceleratorFromNSEvent(event)
	return accelerator != "" && processKeyEquivalent(v.(uint), accelerator)
}

var (
	registerPanelClassOnce sync.Once
	webviewPanelClass      id
)

// panelOverrides are the selectors the NSPanel subclass implements itself.
// canBecomeMainWindow replaces the window's YES: a panel never becomes main.
var panelOverrides = []objc.MethodDef{
	{Cmd: sel_("canBecomeMainWindow"), Fn: func(self objc.ID, cmd objc.SEL) bool { return false }},
	{Cmd: sel_("sendEvent:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {
		if panelHandleKeyDown(id(self), id(event)) {
			return
		}
		sendSuperFrom(webviewPanelClass, self, "sendEvent:", event)
	}},
}

// panelMethodSet returns the window method set with the panel overrides
// replacing their counterparts, so no selector is registered twice.
func panelMethodSet() []objc.MethodDef {
	overridden := make(map[objc.SEL]bool, len(panelOverrides))
	for _, m := range panelOverrides {
		overridden[m.Cmd] = true
	}
	var methods []objc.MethodDef
	for _, m := range webviewWindowMethods(&webviewPanelClass) {
		if !overridden[m.Cmd] {
			methods = append(methods, m)
		}
	}
	return append(methods, panelOverrides...)
}

// registerWebviewPanelClass registers WailsWebviewNSPanel, the NSPanel subclass
// used for panel windows. It is registered once and cached.
func registerWebviewPanelClass() id {
	registerPanelClassOnce.Do(func() {
		webviewPanelClass = registerDelegateClass("WailsWebviewNSPanel", "NSPanel", nil, panelMethodSet())
	})
	return webviewPanelClass
}
