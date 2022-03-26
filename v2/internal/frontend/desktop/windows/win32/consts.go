package win32

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"syscall"
)

type HRESULT int32
type HANDLE uintptr

var moduser32 = syscall.NewLazyDLL("user32.dll")
var procSystemParametersInfo = moduser32.NewProc("SystemParametersInfoW")

var moddwmapi = syscall.NewLazyDLL("dwmapi.dll")
var procDwmSetWindowAttribute = moddwmapi.NewProc("DwmSetWindowAttribute")

var windowsVersion, _ = operatingsystem.GetWindowsVersionInfo()

func IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return windowsVersion.Major >= major &&
		windowsVersion.Minor >= minor &&
		windowsVersion.Build >= buildNumber
}
