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
	procSetProcessDpiAwareness = modshcore.NewProc("SetProcessDpiAwareness")
)

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
