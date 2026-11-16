//go:build darwin && !ios && !server

package application

import "unsafe"

type macosMenuItem struct {
	menuItem *MenuItem

	nsMenuItem unsafe.Pointer

	// macosMenuItemNative is the backend's own state: empty for cgo.
	macosMenuItemNative
}

func (m *macosMenuItem) destroy() {
	if m.nsMenuItem == nil {
		return
	}
	m.destroyNSMenuItem()
	m.nsMenuItem = nil
}
