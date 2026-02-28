//go:build windows

package w32

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modshcore = syscall.NewLazyDLL("shcore.dll")

	procGetDpiForMonitor       = modshcore.NewProc("GetDpiForMonitor")
	procGetProcessDpiAwareness = modshcore.NewProc("GetProcessDpiAwareness")
	procSetProcessDpiAwareness = modshcore.NewProc("SetProcessDpiAwareness")
)

func HasGetProcessDpiAwarenessFunc() bool {
	err := procGetProcessDpiAwareness.Find()
	return err == nil
}

// GetProcessDpiAwareness retrieves the DPI awareness of the current process.
// Returns one of: PROCESS_DPI_UNAWARE, PROCESS_SYSTEM_DPI_AWARE, or PROCESS_PER_MONITOR_DPI_AWARE.
func GetProcessDpiAwareness() (uint, error) {
	var awareness uint
	status, _, err := procGetProcessDpiAwareness.Call(0, uintptr(unsafe.Pointer(&awareness)))
	if status != S_OK {
		return 0, fmt.Errorf("GetProcessDpiAwareness failed: %v", err)
	}
	return awareness, nil
}

func HasSetProcessDpiAwarenessFunc() bool {
	err := procSetProcessDpiAwareness.Find()
	return err == nil
}

func SetProcessDpiAwareness(val uint) error {
	status, r, err := procSetProcessDpiAwareness.Call(uintptr(val))
	if status != S_OK {
		return fmt.Errorf("procSetProcessDpiAwareness failed %d: %v %v", status, r, err)
	}
	return nil
}

func HasGetDPIForMonitorFunc() bool {
	err := procGetDpiForMonitor.Find()
	return err == nil
}

func GetDPIForMonitor(hmonitor HMONITOR, dpiType MONITOR_DPI_TYPE, dpiX *UINT, dpiY *UINT) uintptr {
	ret, _, _ := procGetDpiForMonitor.Call(
		hmonitor,
		uintptr(dpiType),
		uintptr(unsafe.Pointer(dpiX)),
		uintptr(unsafe.Pointer(dpiY)))

	return ret
}

func GetNotificationFlyoutBounds() (*RECT, error) {
	var rect RECT
	res, _, err := procSystemParametersInfo.Call(SPI_GETNOTIFYWINDOWRECT, 0, uintptr(unsafe.Pointer(&rect)), 0)
	if res == 0 {
		_ = err
		return nil, err
	}
	return &rect, nil
}
