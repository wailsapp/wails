package operatingsystem

import (
	"os/exec"
	"strings"
)

func platformInfo() (*OS, error) {
	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "FreeBSD"
	result.Version = "Unknown"

	output, err := exec.Command("sysctl", "-a", "kern.osrevision").Output()
	id := ""
	if err == nil {
		id = strings.TrimSpace(strings.Split(string(output), ":")[1])
	}
	result.ID = id

	output, err = exec.Command("freebsd-version").Output()
	version := ""
	if err == nil {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	result.Version = version

	return &result, nil
}
