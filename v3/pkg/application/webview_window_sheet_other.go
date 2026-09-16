//go:build !darwin || ios || server

package application

import "unsafe"

// Off macOS window sheets are documented no-ops; PresentSheet returns
// ErrMacSheetUnsupported and the queries report no sheet.

func macSheetsSupported() bool { return false }

func macSheetBegin(unsafe.Pointer, unsafe.Pointer, bool) error { return ErrMacSheetUnsupported }
func macSheetEnd(unsafe.Pointer, int)                          {}
func macSheetParent(unsafe.Pointer) unsafe.Pointer             { return nil }
func macSheetAttached(unsafe.Pointer) unsafe.Pointer           { return nil }
