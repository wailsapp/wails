//go:build windows

package operatingsystem

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func stripNulls(str string) string {
	// Split the string into substrings at each null character
	substrings := strings.Split(str, "\x00")

	// Join the substrings back into a single string
	strippedStr := strings.Join(substrings, "")

	return strippedStr
}

func mustStringToUTF16Ptr(input string) *uint16 {
	input = stripNulls(input)
	result, err := syscall.UTF16PtrFromString(input)
	if err != nil {
		panic(err)
	}
	return result
}

func getBranding() string {
	var modBranding = syscall.NewLazyDLL("winbrand.dll")
	var brandingFormatString = modBranding.NewProc("BrandingFormatString")

	windowsLong := mustStringToUTF16Ptr("%WINDOWS_LONG%\x00")
	ret, _, _ := brandingFormatString.Call(
		uintptr(unsafe.Pointer(windowsLong)),
	)
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ret)))
}

func platformInfo() (*OS, error) {
	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "Windows"
	result.Version = "Unknown"

	// Credit: https://stackoverflow.com/a/33288328
	// Ignore errors as it isn't a showstopper
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)

	productName, _, _ := key.GetStringValue("ProductName")
	currentBuild, _, _ := key.GetStringValue("CurrentBuildNumber")
	displayVersion, _, _ := key.GetStringValue("DisplayVersion")
	releaseId, _, _ := key.GetStringValue("ReleaseId")

	result.Name = productName
	result.Version = fmt.Sprintf("%s (Build: %s)", releaseId, currentBuild)
	result.ID = displayVersion
	result.Branding = getBranding()

	return &result, key.Close()
}
