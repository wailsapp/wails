//go:build windows

package win32

import (
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

type HRESULT int32
type HANDLE uintptr
type HMONITOR HANDLE
type HWND HANDLE
type HICON HANDLE

type NOTIFYICONDATA struct {
	CbSize           uint32
	HWnd             HWND
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            HICON
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
}

const (
	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004
)

var (
	moduser32                      = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo       = moduser32.NewProc("SystemParametersInfoW")
	procGetWindowLong              = moduser32.NewProc("GetWindowLongW")
	procSetClassLong               = moduser32.NewProc("SetClassLongW")
	procSetClassLongPtr            = moduser32.NewProc("SetClassLongPtrW")
	procShowWindow                 = moduser32.NewProc("ShowWindow")
	procIsWindowVisible            = moduser32.NewProc("IsWindowVisible")
	procGetWindowRect              = moduser32.NewProc("GetWindowRect")
	procGetMonitorInfo             = moduser32.NewProc("GetMonitorInfoW")
	procMonitorFromWindow          = moduser32.NewProc("MonitorFromWindow")
	procIsClipboardFormatAvailable = moduser32.NewProc("IsClipboardFormatAvailable")
	procOpenClipboard              = moduser32.NewProc("OpenClipboard")
	procCloseClipboard             = moduser32.NewProc("CloseClipboard")
	procEmptyClipboard             = moduser32.NewProc("EmptyClipboard")
	procGetClipboardData           = moduser32.NewProc("GetClipboardData")
	procSetClipboardData           = moduser32.NewProc("SetClipboardData")
	procCreateIconFromResourceEx   = moduser32.NewProc("CreateIconFromResourceEx")
)
var (
	moddwmapi                        = syscall.NewLazyDLL("dwmapi.dll")
	procDwmSetWindowAttribute        = moddwmapi.NewProc("DwmSetWindowAttribute")
	procDwmExtendFrameIntoClientArea = moddwmapi.NewProc("DwmExtendFrameIntoClientArea")
)
var (
	modshell32          = syscall.NewLazyDLL("shell32.dll")
	procShellNotifyIcon = modshell32.NewProc("Shell_NotifyIconW")
)
var (
	modwingdi            = syscall.NewLazyDLL("gdi32.dll")
	procCreateSolidBrush = modwingdi.NewProc("CreateSolidBrush")
)
var (
	kernel32           = syscall.NewLazyDLL("kernel32")
	kernelGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	kernelGlobalFree   = kernel32.NewProc("GlobalFree")
	kernelGlobalLock   = kernel32.NewProc("GlobalLock")
	kernelGlobalUnlock = kernel32.NewProc("GlobalUnlock")
	kernelLstrcpy      = kernel32.NewProc("lstrcpyW")
)

var windowsVersion, _ = operatingsystem.GetWindowsVersionInfo()

func IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return windowsVersion.Major >= major &&
		windowsVersion.Minor >= minor &&
		windowsVersion.Build >= buildNumber
}
