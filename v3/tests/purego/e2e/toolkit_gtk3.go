//go:build linux && gtk3

package main

// GTK3 opens a menu with nothing highlighted: the first Down selects the first
// item. Its context menus take the keyboard.
var menuNavigateKeys = []string{"down", "return"}

const contextMenuByMouse = false
