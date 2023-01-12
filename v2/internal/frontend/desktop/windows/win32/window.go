//go:build windows

package win32

import (
	"fmt"
	"log"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
)

const (
	WS_MAXIMIZE = 0x01000000
	WS_MINIMIZE = 0x20000000

	GWL_STYLE = -16

	MONITOR_DEFAULTTOPRIMARY = 0x00000001
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

const (
	GCLP_HBRBACKGROUND int32 = -10
)

// Power
const (
	// WM_POWERBROADCAST - Notifies applications that a power-management event has occurred.
	WM_POWERBROADCAST = 536

	// PBT_APMPOWERSTATUSCHANGE - Power status has changed.
	PBT_APMPOWERSTATUSCHANGE = 10

	// PBT_APMRESUMEAUTOMATIC -Operation is resuming automatically from a low-power state. This message is sent every time the system resumes.
	PBT_APMRESUMEAUTOMATIC = 18

	// PBT_APMRESUMESUSPEND - Operation is resuming from a low-power state. This message is sent after PBT_APMRESUMEAUTOMATIC if the resume is triggered by user input, such as pressing a key.
	PBT_APMRESUMESUSPEND = 7

	// PBT_APMSUSPEND - System is suspending operation.
	PBT_APMSUSPEND = 4

	// PBT_POWERSETTINGCHANGE - A power setting change event has been received.
	PBT_POWERSETTINGCHANGE = 32787
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb773244.aspx
type MARGINS struct {
	CxLeftWidth, CxRightWidth, CyTopHeight, CyBottomHeight int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162897.aspx
type RECT struct {
	Left, Top, Right, Bottom int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145065.aspx
type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

func ExtendFrameIntoClientArea(hwnd uintptr, extend bool) {
	// -1: Adds the default frame styling (aero shadow and e.g. rounded corners on Windows 11)
	//     Also shows the caption buttons if transparent ant translucent but they don't work.
	//  0: Adds the default frame styling but no aero shadow, does not show the caption buttons.
	//  1: Adds the default frame styling (aero shadow and e.g. rounded corners on Windows 11) but no caption buttons
	//     are shown if transparent ant translucent.
	var margins MARGINS
	if extend {
		margins = MARGINS{1, 1, 1, 1} // Only extend 1 pixel to have the default frame styling but no caption buttons
	}
	if err := dwmExtendFrameIntoClientArea(hwnd, &margins); err != nil {
		log.Fatal(fmt.Errorf("DwmExtendFrameIntoClientArea failed: %s", err))
	}
}

func IsVisible(hwnd uintptr) bool {
	ret, _, _ := procIsWindowVisible.Call(hwnd)
	return ret != 0
}

func IsWindowFullScreen(hwnd uintptr) bool {
	wRect := GetWindowRect(hwnd)
	m := MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY)
	var mi MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	if !GetMonitorInfo(m, &mi) {
		return false
	}
	return wRect.Left == mi.RcMonitor.Left &&
		wRect.Top == mi.RcMonitor.Top &&
		wRect.Right == mi.RcMonitor.Right &&
		wRect.Bottom == mi.RcMonitor.Bottom
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

func SetBackgroundColour(hwnd uintptr, r, g, b uint8) {
	col := winc.RGB(r, g, b)
	hbrush, _, _ := procCreateSolidBrush.Call(uintptr(col))
	setClassLongPtr(hwnd, GCLP_HBRBACKGROUND, hbrush)
}

func IsWindowNormal(hwnd uintptr) bool {
	return !IsWindowMaximised(hwnd) && !IsWindowMinimised(hwnd) && !IsWindowFullScreen(hwnd)
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

func setClassLongPtr(hwnd uintptr, param int32, val uintptr) bool {
	proc := procSetClassLongPtr
	if strconv.IntSize == 32 {
		/*
			https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setclasslongptrw
			Note: 	To write code that is compatible with both 32-bit and 64-bit Windows, use SetClassLongPtr.
					When compiling for 32-bit Windows, SetClassLongPtr is defined as a call to the SetClassLong function

			=> We have to do this dynamically when directly calling the DLL procedures
		*/
		proc = procSetClassLong
	}

	ret, _, _ := proc.Call(
		hwnd,
		uintptr(param),
		val,
	)
	return ret != 0
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

func GetWindowRect(hwnd uintptr) *RECT {
	var rect RECT
	procGetWindowRect.Call(
		hwnd,
		uintptr(unsafe.Pointer(&rect)))

	return &rect
}

func MonitorFromWindow(hwnd uintptr, dwFlags uint32) HMONITOR {
	ret, _, _ := procMonitorFromWindow.Call(
		hwnd,
		uintptr(dwFlags),
	)
	return HMONITOR(ret)
}

func GetMonitorInfo(hMonitor HMONITOR, lmpi *MONITORINFO) bool {
	ret, _, _ := procGetMonitorInfo.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lmpi)),
	)
	return ret != 0
}
