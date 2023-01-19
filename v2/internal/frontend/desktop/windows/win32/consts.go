//go:build windows

package win32

import (
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

type HRESULT int32
type HANDLE uintptr
type HMONITOR HANDLE

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
)
var (
	moddwmapi                        = syscall.NewLazyDLL("dwmapi.dll")
	procDwmSetWindowAttribute        = moddwmapi.NewProc("DwmSetWindowAttribute")
	procDwmExtendFrameIntoClientArea = moddwmapi.NewProc("DwmExtendFrameIntoClientArea")
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
