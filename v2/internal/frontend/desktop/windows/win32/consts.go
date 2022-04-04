//go:build windows

package win32

import (
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

type HRESULT int32
type HANDLE uintptr

var (
	moduser32                = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo = moduser32.NewProc("SystemParametersInfoW")
	procGetWindowLong        = moduser32.NewProc("GetWindowLongW")
)
var (
	moddwmapi                        = syscall.NewLazyDLL("dwmapi.dll")
	procDwmSetWindowAttribute        = moddwmapi.NewProc("DwmSetWindowAttribute")
	procDwmExtendFrameIntoClientArea = moddwmapi.NewProc("DwmExtendFrameIntoClientArea")
)
var windowsVersion, _ = operatingsystem.GetWindowsVersionInfo()

func IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return windowsVersion.Major >= major &&
		windowsVersion.Minor >= minor &&
		windowsVersion.Build >= buildNumber
}
