//go:build windows

package w32

import "syscall"

var (
	moddwmapi = syscall.NewLazyDLL("dwmapi.dll")

	procDwmSetWindowAttribute = moddwmapi.NewProc("DwmSetWindowAttribute")
)

func DwmSetWindowAttribute(hwnd HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute LPCVOID, cbAttribute uint32) HRESULT {
	ret, _, _ := procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(dwAttribute),
		uintptr(pvAttribute),
		uintptr(cbAttribute))
	return HRESULT(ret)
}
