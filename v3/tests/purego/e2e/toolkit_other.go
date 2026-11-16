//go:build !(linux && gtk3)

package main

// GTK4 (the default Linux build) highlights the first item as soon as F10
// opens the menu bar, and shows context menus as popovers that do not take
// the X keyboard focus without a window manager, so its context menu is
// driven with the mouse. macOS does not use these.
const (
	menuOpensOnFirstItem = true
	contextMenuByMouse   = true
)
