//go:build windows

package operatingsystem

import (
	"fmt"
	"strconv"
	"strings"

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

	// Windows 10 and 11 both declare the product name as "Windows 10"
	// We determine the difference using the build number:
	// https://en.wikipedia.org/wiki/List_of_Microsoft_Windows_versions#Personal_computer_versions
	if buildNum, err := strconv.Atoi(currentBuild); err == nil {
		if buildNum >= 22000 {
			productName = strings.Replace(productName, " 10", " 11", -1)
		}
	}

	result.Name = productName
	result.Version = fmt.Sprintf("%s (Build: %s)", releaseId, currentBuild)
	result.ID = displayVersion

	return &result, key.Close()
}
