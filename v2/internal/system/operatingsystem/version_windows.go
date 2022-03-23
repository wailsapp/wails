package operatingsystem

import (
	"golang.org/x/sys/windows/registry"
	"strconv"
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
		Major:          regKeyAsInt(key, "CurrentMajorVersionNumber"),
		Minor:          regKeyAsInt(key, "CurrentMinorVersionNumber"),
		Build:          regKeyAsInt(key, "CurrentBuildNumber"),
		DisplayVersion: regKeyAsString(key, "DisplayVersion"),
	}, nil
}

func regKeyAsInt(key registry.Key, name string) int {
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
