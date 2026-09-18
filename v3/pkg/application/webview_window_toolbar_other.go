//go:build !darwin || ios || server

package application

import "unsafe"

func macToolbarSetDisplayMode(unsafe.Pointer, MacToolbarDisplayMode) {}
func macToolbarItemSetLabel(unsafe.Pointer, string, string)          {}
func macToolbarItemSetSymbol(unsafe.Pointer, string, string)         {}
func macToolbarItemSetTooltip(unsafe.Pointer, string, string)        {}
func macToolbarItemSetBordered(unsafe.Pointer, string, bool)         {}
func macToolbarItemSetProminent(unsafe.Pointer, string, bool)        {}
func macToolbarItemSetTintColor(unsafe.Pointer, string, *RGBA)       {}
func macToolbarItemSetEnabled(unsafe.Pointer, string, bool)          {}
func macToolbarItemSetHidden(unsafe.Pointer, string, bool)           {}
func macToolbarItemSetBadgeCount(unsafe.Pointer, string, int)        {}
func macToolbarGroupSetSelectedIndex(unsafe.Pointer, string, int)    {}
func macToolbarGroupSetSelectionMode(unsafe.Pointer, string, MacToolbarGroupSelectionMode) {
}
func macToolbarShareItemSetProvider(unsafe.Pointer, string, MacShareProvider, string, string, []MacShareRepresentation) {
}
func macToolbarSetCustomizable(unsafe.Pointer, bool)                   {}
func macToolbarRunCustomizationPalette(unsafe.Pointer)                 {}
func macToolbarSyncNative(*MacToolbar, *MacToolbarItem)                {}
func macToolbarMenuItemSetShowsIndicator(unsafe.Pointer, string, bool) {}
func macToolbarItemSetVisibilityPriority(unsafe.Pointer, string, MacToolbarVisibilityPriority) {
}
func macToolbarItemSetNavigational(unsafe.Pointer, string, bool)         {}
func macToolbarSearchItemSetRecents(unsafe.Pointer, string, string, int) {}
func macToolbarSearchItemSetMenu(unsafe.Pointer, string, *Menu, bool)    {}
func macToolbarSearchItemSetIncremental(unsafe.Pointer, string, bool)    {}
func macToolbarSearchItemSetPlaceholder(unsafe.Pointer, string, string)  {}
