//go:build windows

package win32

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	WS_MAXIMIZE = 0x01000000
	WS_MINIMIZE = 0x20000000

	GWL_STYLE = -16
)

const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb773244.aspx
type MARGINS struct {
	CxLeftWidth, CxRightWidth, CyTopHeight, CyBottomHeight int32
}

func ExtendFrameIntoClientArea(hwnd uintptr) {
	// -1: Adds the default frame styling (aero shadow and e.g. rounded corners on Windows 11)
	//     Also shows the caption buttons if transparent ant translucent but they don't work.
	//  0: Adds the default frame styling but no aero shadow, does not show the caption buttons.
	//  1: Adds the default frame styling (aero shadow and e.g. rounded corners on Windows 11) but no caption buttons
	//     are shown if transparent ant translucent.
	margins := &MARGINS{1, 1, 1, 1} // Only extend 1 pixel to have the default frame styling but no caption buttons
	if err := dwmExtendFrameIntoClientArea(hwnd, margins); err != nil {
		log.Fatal(fmt.Errorf("DwmExtendFrameIntoClientArea failed: %s", err))
	}
}

func IsWindowMaximised(hwnd uintptr) bool {
	style := uint32(getWindowLong(hwnd, GWL_STYLE))
	return style&WS_MAXIMIZE != 0
}
func IsWindowMinimised(hwnd uintptr) bool {
	style := uint32(getWindowLong(hwnd, GWL_STYLE))
	return style&WS_MINIMIZE != 0
}

func RestoreWindow(hwnd uintptr) {
	showWindow(hwnd, SW_RESTORE)
}

func ShowWindow(hwnd uintptr) {
	showWindow(hwnd, SW_SHOW)
}

func ShowWindowMaximised(hwnd uintptr) {
	showWindow(hwnd, SW_MAXIMIZE)
}
func ShowWindowMinimised(hwnd uintptr) {
	showWindow(hwnd, SW_MINIMIZE)
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

func getWindowLong(hwnd uintptr, index int) int32 {
	ret, _, _ := procGetWindowLong.Call(
		hwnd,
		uintptr(index))

	return int32(ret)
}

func showWindow(hwnd uintptr, cmdshow int) bool {
	ret, _, _ := procShowWindow.Call(
		hwnd,
		uintptr(cmdshow))

	return ret != 0
}
