package cmd

import "fmt"

// LinuxDistribution is of type int
type LinuxDistribution int

const (
	// Ubuntu distro
	Ubuntu LinuxDistribution = 0
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	distribution LinuxDistribution
	name         string
	release      string
}

func getLinuxDistroInfo() *DistroInfo {
	result := &DistroInfo{}
	program := NewProgramHelper()
	// Does lsb_release exist?

	lsbRelease := program.FindProgram("lsb_release")
	if lsbRelease != nil {
		stdout, _, err := lsbRelease.Run("-a")
		if err != nil {
			return nil
		}
		fmt.Println(stdout)
	}
	return result
}
