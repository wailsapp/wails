package cmd

import (
	"fmt"
	"strings"
)

// LinuxDistribution is of type int
type LinuxDistribution int

const (
	// Unknown is the catch-all distro
	Unknown LinuxDistribution = 0
	// Ubuntu distribution
	Ubuntu  LinuxDistribution = 1
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	Distribution  LinuxDistribution
	Description   string
	Release       string
	Codename      string
	DistributorID string
}

// GetLinuxDistroInfo returns information about the running linux distribution
func GetLinuxDistroInfo() *DistroInfo {
	result := &DistroInfo{Distribution: Unknown}
	program := NewProgramHelper()
	// Does lsb_release exist?

	lsbRelease := program.FindProgram("lsb_release")
	if lsbRelease != nil {
		stdout, _, _, err := lsbRelease.Run("-a")
		if err != nil {
			return result
		}

		for _, line := range strings.Split(stdout, "\n") {
			if strings.Contains(line, ":") {
				// Iterate lines a
				details := strings.Split(line, ":")
				key := strings.TrimSpace(details[0])
				value := strings.TrimSpace(details[1])
				switch key {
				case "Distributor ID":
					result.DistributorID = value
					switch value {
					case "Ubuntu":
						result.Distribution = Ubuntu
					}
				case "Description":
					result.Description = value
				case "Release":
					result.Release = value
				case "Codename":
					result.Codename = value

				}
			}
		}

	}
	return result
}

// DpkgInstalled uses dpkg to see if a package is installed
func DpkgInstalled(packageName string) (bool, error) {
	result := false
	program := NewProgramHelper()
	dpkg := program.FindProgram("dpkg")
	if dpkg == nil {
		return false, fmt.Errorf("cannot check dependencies: dpkg not found")
	}
	_, _, exitCode, _ := dpkg.Run("-L", packageName)
	result = exitCode == 0
	return result, nil
}
