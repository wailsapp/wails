//go:build windows

package w32

import (
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procSetWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")
var windowsVersion, _ = operatingsystem.GetWindowsVersionInfo()

const (
	WDA_NONE               = 0x00000000
	WDA_MONITOR            = 0x00000001
	WDA_EXCLUDEFROMCAPTURE = 0x00000011 // windows 10 2004+
)

func isWindowsVersionAtLeast(major, minor, build int) bool {
	if windowsVersion.Major > major {
		return true
	}
	if windowsVersion.Major < major {
		return false
	}
	if windowsVersion.Minor > minor {
		return true
	}
	if windowsVersion.Minor < minor {
		return false
	}
	return windowsVersion.Build >= build
}

func SetWindowDisplayAffinity(hwnd uintptr, affinity uint32) bool {
	if affinity == WDA_EXCLUDEFROMCAPTURE && !isWindowsVersionAtLeast(10, 0, 19041) {
		// for older windows versions, use WDA_MONITOR
		affinity = WDA_MONITOR
	}
	ret, _, _ := procSetWindowDisplayAffinity.Call(
		hwnd,
		uintptr(affinity),
	)
	return ret != 0
}
