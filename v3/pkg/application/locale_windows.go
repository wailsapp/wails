//go:build windows && !server

package application

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName = kernel32.NewProc("GetUserDefaultLocaleName")
)

// SystemLocale returns the system's configured locale as a BCP-47 language tag
// (e.g. "nb-NO", "en-US").
func SystemLocale() string {
	buf := make([]uint16, 85) // LOCALE_NAME_MAX_LENGTH
	r, _, _ := procGetUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if r == 0 {
		return "en"
	}
	return syscall.UTF16ToString(buf)
}
