//go:build linux && gtk3

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// GTK3 opens a menu with nothing highlighted: the first Down selects the first
// item. Its context menus take the keyboard.
func menuNavigateKeys(*application.App) []string { return []string{"down", "return"} }

const contextMenuByMouse = false
