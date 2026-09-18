//go:build !darwin || ios || server

package application

import "unsafe"

// Off macOS popovers are documented no-ops: every show method returns
// ErrMacPopoverUnsupported and IsShown is false.

func macPopoversSupported() bool { return false }

func macPopoverEnsureNative(*MacPopover) (unsafe.Pointer, error) {
	return nil, ErrMacPopoverUnsupported
}

func macPopoverShowInWindow(unsafe.Pointer, unsafe.Pointer, Rect, MacRectEdge) error {
	return ErrMacPopoverUnsupported
}

func macPopoverShowFromToolbarItem(unsafe.Pointer, unsafe.Pointer, string, MacRectEdge) error {
	return ErrMacPopoverUnsupported
}

func macPopoverShowFromStatusItem(unsafe.Pointer, unsafe.Pointer, MacRectEdge) error {
	return ErrMacPopoverUnsupported
}

func macPopoverClose(unsafe.Pointer)                            {}
func macPopoverIsShown(unsafe.Pointer) bool                     { return false }
func macPopoverSetContentSize(unsafe.Pointer, float64, float64) {}
func macPopoverSetBehavior(unsafe.Pointer, MacPopoverBehavior)  {}
func macPopoverRelease(*MacPopover)                             {}
func macSystemTrayStatusItem(*SystemTray) unsafe.Pointer        { return nil }
