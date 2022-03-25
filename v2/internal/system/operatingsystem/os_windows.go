package operatingsystem

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

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

	return &result, key.Close()
}

func UseDarkMode() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	AppsUseLightTheme, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return false
	}

	return AppsUseLightTheme == 0
}
