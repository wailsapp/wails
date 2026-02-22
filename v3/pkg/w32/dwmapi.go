//go:build windows

package w32

import (
	"syscall"
	"unsafe"
)

var (
	moddwmapi = syscall.NewLazyDLL("dwmapi.dll")

	procDwmSetWindowAttribute        = moddwmapi.NewProc("DwmSetWindowAttribute")
	procDwmGetWindowAttribute        = moddwmapi.NewProc("DwmGetWindowAttribute")
	procDwmExtendFrameIntoClientArea = moddwmapi.NewProc("DwmExtendFrameIntoClientArea")
)

func DwmSetWindowAttribute(hwnd HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute unsafe.Pointer, cbAttribute uintptr) HRESULT {
	ret, _, _ := procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(dwAttribute),
		uintptr(pvAttribute),
		cbAttribute)
	return HRESULT(ret)
}

func DwmGetWindowAttribute(hwnd HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute unsafe.Pointer, cbAttribute uintptr) HRESULT {
	ret, _, _ := procDwmGetWindowAttribute.Call(
		hwnd,
		uintptr(dwAttribute),
		uintptr(pvAttribute),
		cbAttribute)
	return HRESULT(ret)
}

func dwmExtendFrameIntoClientArea(hwnd uintptr, margins *MARGINS) error {
	ret, _, _ := procDwmExtendFrameIntoClientArea.Call(
		hwnd,
		uintptr(unsafe.Pointer(margins)))

	if ret != 0 {
		return syscall.GetLastError()
	}

	return nil
}
