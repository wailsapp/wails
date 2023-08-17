//go:build windows

package operatingsystem

import (
	"strconv"

	"golang.org/x/sys/windows/registry"
)

type WindowsVersionInfo struct {
	Major          int
	Minor          int
	Build          int
	DisplayVersion string
}

func (w *WindowsVersionInfo) IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return w.Major >= major && w.Minor >= minor && w.Build >= buildNumber
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
