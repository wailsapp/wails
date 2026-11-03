//go:build !(linux && gtk3)

package main

import "runtime"

// GTK4 (the default Linux build) highlights the first item as soon as F10
// opens the menu bar, and shows context menus as popovers that do not take
// the X keyboard focus without a window manager, so its context menu is
// driven with the mouse. On Windows F10 selects the menu bar's first menu
// without opening it, Down opens it on its first item, and popup menus take
// the keyboard. macOS does not navigate the menu bar and its context menus
// take the keyboard.
var menuNavigateKeys = map[string][]string{
	"linux":   {"return"},
	"windows": {"down", "return"},
}[runtime.GOOS]

const contextMenuByMouse = runtime.GOOS == "linux"
