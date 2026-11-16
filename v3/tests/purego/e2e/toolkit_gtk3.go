//go:build linux && gtk3

package main

// GTK3 opens a menu with nothing highlighted: the first Down selects the first
// item. Its context menus take the keyboard.
const (
	menuOpensOnFirstItem = false
	contextMenuByMouse   = false
)
