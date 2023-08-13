//go:build windows

package w32

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"strconv"
	"syscall"
	"unsafe"
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
var (
	modBranding          = syscall.NewLazyDLL("winbrand.dll")
	brandingFormatString = modBranding.NewProc("BrandingFormatString")
)

var windowsVersion, _ = GetWindowsVersionInfo()

func IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return windowsVersion.Major >= major &&
		windowsVersion.Minor >= minor &&
		windowsVersion.Build >= buildNumber
}

type WindowsVersionInfo struct {
	Major          int
	Minor          int
	Build          int
	DisplayVersion string
}

func (w *WindowsVersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d (%s)", w.Major, w.Minor, w.Build, w.DisplayVersion)
}

func (w *WindowsVersionInfo) IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return w.Major >= major && w.Minor >= minor && w.Build >= buildNumber
}

func GetBranding() string {
	windowsLong := MustStringToUTF16Ptr("%WINDOWS_LONG%\x00")
	ret, _, _ := brandingFormatString.Call(
		uintptr(unsafe.Pointer(windowsLong)),
	)
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ret)))
}

func GetWindowsVersionInfo() (*WindowsVersionInfo, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}

	return &WindowsVersionInfo{
		Major:          regDWORDKeyAsInt(key, "CurrentMajorVersionNumber"),
		Minor:          regDWORDKeyAsInt(key, "CurrentMinorVersionNumber"),
		Build:          regStringKeyAsInt(key, "CurrentBuildNumber"),
		DisplayVersion: regKeyAsString(key, "DisplayVersion"),
	}, nil
}

func regDWORDKeyAsInt(key registry.Key, name string) int {
	result, _, err := key.GetIntegerValue(name)
	if err != nil {
		return -1
	}
	return int(result)
}

func regStringKeyAsInt(key registry.Key, name string) int {
	resultStr, _, err := key.GetStringValue(name)
	if err != nil {
		return -1
	}
	result, err := strconv.Atoi(resultStr)
	if err != nil {
		return -1
	}
	return result
}

func regKeyAsString(key registry.Key, name string) string {
	resultStr, _, err := key.GetStringValue(name)
	if err != nil {
		return ""
	}
	return resultStr
}
