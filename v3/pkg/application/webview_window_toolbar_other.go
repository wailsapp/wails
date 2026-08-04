//go:build !darwin || ios || server

package application

import "unsafe"

func macToolbarItemSetLabel(unsafe.Pointer, string, string)       {}
func macToolbarItemSetSymbol(unsafe.Pointer, string, string)      {}
func macToolbarItemSetTooltip(unsafe.Pointer, string, string)     {}
func macToolbarItemSetBordered(unsafe.Pointer, string, bool)      {}
func macToolbarItemSetProminent(unsafe.Pointer, string, bool)     {}
func macToolbarItemSetTintColor(unsafe.Pointer, string, *RGBA)    {}
func macToolbarItemSetEnabled(unsafe.Pointer, string, bool)       {}
func macToolbarItemSetHidden(unsafe.Pointer, string, bool)        {}
func macToolbarItemSetBadgeCount(unsafe.Pointer, string, int)     {}
func macToolbarGroupSetSelectedIndex(unsafe.Pointer, string, int) {}
func macToolbarGroupSetSelectionMode(unsafe.Pointer, string, MacToolbarGroupSelectionMode) {
}
