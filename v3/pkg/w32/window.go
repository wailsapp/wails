//go:build windows

package w32

import (
	"fmt"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const (
	SC_CLOSE    = 0xF060
	SC_MOVE     = 0xF010
	SC_MAXIMIZE = 0xF030
	SC_MINIMIZE = 0xF020
	SC_SIZE     = 0xF000
	SC_RESTORE  = 0xF120
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	getSystemMenu  = user32.NewProc("GetSystemMenu")
	enableMenuItem = user32.NewProc("EnableMenuItem")
	findWindow     = user32.NewProc("FindWindowW")
	sendMessage    = user32.NewProc("SendMessageW")
)

const (
	WMCOPYDATA_SINGLE_INSTANCE_DATA = 1542
)

type COPYDATASTRUCT struct {
	DwData uintptr
	CbData uint32
	LpData uintptr
}

var Fatal func(error)

const (
	GCLP_HBRBACKGROUND int32 = -10
	GCLP_HICON         int32 = -14
)

type WINDOWPOS struct {
	HwndInsertAfter HWND
	X               int32
	Y               int32
	Cx              int32
	Cy              int32
	Flags           uint32
}

func ExtendFrameIntoClientArea(hwnd uintptr, extend bool) error {
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
		return fmt.Errorf("DwmExtendFrameIntoClientArea failed: %s", err)
	}
	return nil
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

func ShowWindowMaximised(hwnd uintptr) {
	showWindow(hwnd, SW_MAXIMIZE)
}
func ShowWindowMinimised(hwnd uintptr) {
	showWindow(hwnd, SW_MINIMIZE)
}

func SetApplicationIcon(hwnd uintptr, icon HICON) {
	setClassLongPtr(hwnd, GCLP_HICON, icon)
}

func SetBackgroundColour(hwnd uintptr, r, g, b uint8) {
	col := uint32(r) | uint32(g)<<8 | uint32(b)<<16
	hbrush, _, _ := procCreateSolidBrush.Call(uintptr(col))
	setClassLongPtr(hwnd, GCLP_HBRBACKGROUND, hbrush)
}

func IsWindowNormal(hwnd uintptr) bool {
	return !IsWindowMaximised(hwnd) && !IsWindowMinimised(hwnd) && !IsWindowFullScreen(hwnd)
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

func stripNulls(str string) string {
	// Split the string into substrings at each null character
	substrings := strings.Split(str, "\x00")

	// Join the substrings back into a single string
	strippedStr := strings.Join(substrings, "")

	return strippedStr
}

func MustStringToUTF16Ptr(input string) *uint16 {
	input = stripNulls(input)
	result, err := syscall.UTF16PtrFromString(input)
	if err != nil {
		Fatal(err)
	}
	return result
}

func MustStringToUTF16uintptr(input string) uintptr {
	input = stripNulls(input)
	ret := lo.Must(syscall.UTF16PtrFromString(input))
	return uintptr(unsafe.Pointer(ret))
}

func MustStringToUTF16(input string) []uint16 {
	input = stripNulls(input)
	return lo.Must(syscall.UTF16FromString(input))
}

func CenterWindow(hwnd HWND) {
	windowInfo := getWindowInfo(hwnd)
	frameless := windowInfo.IsPopup()

	info := GetMonitorInfoForWindow(hwnd)
	workRect := info.RcWork
	screenMiddleW := workRect.Left + (workRect.Right-workRect.Left)/2
	screenMiddleH := workRect.Top + (workRect.Bottom-workRect.Top)/2
	var winRect *RECT
	if !frameless {
		winRect = GetWindowRect(hwnd)
	} else {
		winRect = GetClientRect(hwnd)
	}
	winWidth := winRect.Right - winRect.Left
	winHeight := winRect.Bottom - winRect.Top
	windowX := screenMiddleW - (winWidth / 2)
	windowY := screenMiddleH - (winHeight / 2)
	SetWindowPos(hwnd, HWND_TOP, int(windowX), int(windowY), int(winWidth), int(winHeight), SWP_NOSIZE)
}

func getWindowInfo(hwnd HWND) *WINDOWINFO {
	var info WINDOWINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	GetWindowInfo(hwnd, &info)
	return &info
}

func GetMonitorInfoForWindow(hwnd HWND) *MONITORINFO {
	currentMonitor := MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST)
	var info MONITORINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	GetMonitorInfo(currentMonitor, &info)
	return &info
}

type WindowProc func(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr

var windowClasses = make(map[string]HINSTANCE)
var windowClassesLock sync.Mutex

func getWindowClass(name string) (HINSTANCE, bool) {
	windowClassesLock.Lock()
	defer windowClassesLock.Unlock()
	result, exists := windowClasses[name]
	return result, exists
}

func setWindowClass(name string, instance HINSTANCE) {
	windowClassesLock.Lock()
	defer windowClassesLock.Unlock()
	windowClasses[name] = instance
}

func RegisterWindow(name string, proc WindowProc) (HINSTANCE, error) {
	classInstance, exists := getWindowClass(name)
	if exists {
		return classInstance, nil
	}
	applicationInstance := GetModuleHandle("")
	if applicationInstance == 0 {
		return 0, fmt.Errorf("get module handle failed")
	}

	var wc WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.WndProc = syscall.NewCallback(proc)
	wc.Instance = applicationInstance
	wc.Icon = LoadIconWithResourceID(0, uint16(IDI_APPLICATION))
	wc.Cursor = LoadCursorWithResourceID(0, uint16(IDC_ARROW))
	wc.Background = COLOR_BTNFACE + 1
	wc.ClassName = MustStringToUTF16Ptr(name)

	atom := RegisterClassEx(&wc)
	if atom == 0 {
		panic(syscall.GetLastError())
	}

	setWindowClass(name, applicationInstance)

	return applicationInstance, nil
}

func FlashWindow(hwnd HWND, enabled bool) {
	var flashInfo FLASHWINFO
	flashInfo.CbSize = uint32(unsafe.Sizeof(flashInfo))
	flashInfo.Hwnd = hwnd
	if enabled {
		flashInfo.DwFlags = FLASHW_ALL | FLASHW_TIMERNOFG
	} else {
		flashInfo.DwFlags = FLASHW_STOP
	}
	_, _, _ = procFlashWindowEx.Call(uintptr(unsafe.Pointer(&flashInfo)))
}

func EnumChildWindows(hwnd HWND, callback func(hwnd HWND, lparam LPARAM) LRESULT) LRESULT {
	r, _, _ := procEnumChildWindows.Call(hwnd, syscall.NewCallback(callback), 0)
	return r
}

func DisableCloseButton(hwnd HWND) error {
	hSysMenu, _, err := getSystemMenu.Call(hwnd, 0)
	if hSysMenu == 0 {
		return err
	}

	r1, _, err := enableMenuItem.Call(hSysMenu, SC_CLOSE, MF_BYCOMMAND|MF_DISABLED|MF_GRAYED)
	if r1 == 0 {
		return err
	}

	return nil
}

func EnableCloseButton(hwnd HWND) error {
	hSysMenu, _, err := getSystemMenu.Call(hwnd, 0)
	if hSysMenu == 0 {
		return err
	}

	r1, _, err := enableMenuItem.Call(hSysMenu, SC_CLOSE, MF_BYCOMMAND|MF_ENABLED)
	if r1 == 0 {
		return err
	}

	return nil
}

func FindWindowW(className, windowName *uint16) HWND {
	ret, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
	)
	return HWND(ret)
}

func SendMessageToWindow(hwnd HWND, msg string) {
	// Convert data to UTF16 string
	dataUTF16 := MustStringToUTF16(msg)

	// Prepare COPYDATASTRUCT
	cds := COPYDATASTRUCT{
		DwData: WMCOPYDATA_SINGLE_INSTANCE_DATA,
		CbData: uint32((len(dataUTF16) * 2) + 1), // +1 for null terminator
		LpData: uintptr(unsafe.Pointer(&dataUTF16[0])),
	}

	// Send message to first instance
	_, _, _ = procSendMessage.Call(
		hwnd,
		WM_COPYDATA,
		0,
		uintptr(unsafe.Pointer(&cds)),
	)
}
