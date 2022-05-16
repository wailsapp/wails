package w32

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	ole32               = windows.NewLazySystemDLL("ole32")
	Ole32CoInitializeEx = ole32.NewProc("CoInitializeEx")

	kernel32                   = windows.NewLazySystemDLL("kernel32")
	Kernel32GetCurrentThreadID = kernel32.NewProc("GetCurrentThreadId")

	shlwapi                  = windows.NewLazySystemDLL("shlwapi")
	shlwapiSHCreateMemStream = shlwapi.NewProc("SHCreateMemStream")

	user32                   = windows.NewLazySystemDLL("user32")
	User32LoadImageW         = user32.NewProc("LoadImageW")
	User32GetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	User32RegisterClassExW   = user32.NewProc("RegisterClassExW")
	User32CreateWindowExW    = user32.NewProc("CreateWindowExW")
	User32DestroyWindow      = user32.NewProc("DestroyWindow")
	User32ShowWindow         = user32.NewProc("ShowWindow")
	User32UpdateWindow       = user32.NewProc("UpdateWindow")
	User32SetFocus           = user32.NewProc("SetFocus")
	User32GetMessageW        = user32.NewProc("GetMessageW")
	User32TranslateMessage   = user32.NewProc("TranslateMessage")
	User32DispatchMessageW   = user32.NewProc("DispatchMessageW")
	User32DefWindowProcW     = user32.NewProc("DefWindowProcW")
	User32GetClientRect      = user32.NewProc("GetClientRect")
	User32PostQuitMessage    = user32.NewProc("PostQuitMessage")
	User32SetWindowTextW     = user32.NewProc("SetWindowTextW")
	User32PostThreadMessageW = user32.NewProc("PostThreadMessageW")
	User32GetWindowLongPtrW  = user32.NewProc("GetWindowLongPtrW")
	User32SetWindowLongPtrW  = user32.NewProc("SetWindowLongPtrW")
	User32AdjustWindowRect   = user32.NewProc("AdjustWindowRect")
	User32SetWindowPos       = user32.NewProc("SetWindowPos")
)

const (
	SystemMetricsCxIcon = 11
	SystemMetricsCyIcon = 12
)

const (
	SWShow = 5
)

const (
	SWPNoZOrder     = 0x0004
	SWPNoActivate   = 0x0010
	SWPNoMove       = 0x0002
	SWPFrameChanged = 0x0020
)

const (
	WMDestroy       = 0x0002
	WMMove          = 0x0003
	WMSize          = 0x0005
	WMClose         = 0x0010
	WMQuit          = 0x0012
	WMGetMinMaxInfo = 0x0024
	WMNCLButtonDown = 0x00A1
	WMMoving        = 0x0216
	WMApp           = 0x8000
)

const (
	GWLStyle = -16
)

const (
	WSOverlapped       = 0x00000000
	WSMaximizeBox      = 0x00020000
	WSThickFrame       = 0x00040000
	WSCaption          = 0x00C00000
	WSSysMenu          = 0x00080000
	WSMinimizeBox      = 0x00020000
	WSOverlappedWindow = (WSOverlapped | WSCaption | WSSysMenu | WSThickFrame | WSMinimizeBox | WSMaximizeBox)
)

type WndClassExW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type MinMaxInfo struct {
	PtReserved     Point
	PtMaxSize      Point
	PtMaxPosition  Point
	PtMinTrackSize Point
	PtMaxTrackSize Point
}

type Point struct {
	X, Y int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

func Utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	s := (*[(1 << 30) - 1]uint16)(unsafe.Pointer(p))[:n:n]
	return string(utf16.Decode(s))
}

func SHCreateMemStream(data []byte) (uintptr, error) {
	ret, _, err := shlwapiSHCreateMemStream.Call(
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
	)
	if ret == 0 {
		return 0, err
	}

	return ret, nil
}
