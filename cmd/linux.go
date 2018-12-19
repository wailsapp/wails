package cmd

import (
	"strings"
)

// LinuxDistribution is of type int
type LinuxDistribution int

const (
	// Ubuntu distro
	Ubuntu LinuxDistribution = 0
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	Distribution  LinuxDistribution
	Description   string
	Release       string
	Codename      string
	DistributorId string
}

func GetLinuxDistroInfo() *DistroInfo {
	result := &DistroInfo{}
	program := NewProgramHelper()
	// Does lsb_release exist?

	lsbRelease := program.FindProgram("lsb_release")
	if lsbRelease != nil {
		stdout, _, err, _ := lsbRelease.Run("-a")
		if err != nil {
			return nil
		}

		for _, line := range strings.Split(stdout, "\n") {
			if strings.Contains(line, ":") {
				// Iterate lines a
				details := strings.Split(line, ":")
				key := strings.TrimSpace(details[0])
				value := strings.TrimSpace(details[1])
				switch key {
				case "Distributor ID":
					result.DistributorId = value
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
