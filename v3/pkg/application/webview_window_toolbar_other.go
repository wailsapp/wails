//go:build !darwin || ios || server

package application

import "unsafe"

func macToolbarItemSetLabel(unsafe.Pointer, string, string)       {}
func macToolbarItemSetEnabled(unsafe.Pointer, string, bool)       {}
func macToolbarItemSetHidden(unsafe.Pointer, string, bool)        {}
func macToolbarItemSetBadgeCount(unsafe.Pointer, string, int)     {}
func macToolbarGroupSetSelectedIndex(unsafe.Pointer, string, int) {}
