package cmd

import (
	"fmt"
	"runtime"
)

func newPrerequisite(name, help string) *Prerequisite {
	return &Prerequisite{Name: name, Help: help}
}

// Prerequisites is a list of things required to use Wails
type Prerequisites []*Prerequisite

// Add given prereq object to list
func (p *Prerequisites) Add(prereq *Prerequisite) {
	*p = append(*p, prereq)
}

// GetRequiredPrograms returns a list of programs required for the platform
func GetRequiredPrograms() (*Prerequisites, error) {
	switch runtime.GOOS {
	case "darwin":
		return getRequiredProgramsOSX(), nil
	case "linux":
		return getRequiredProgramsLinux(), nil
	case "windows":
		return getRequiredProgramsWindows(), nil
	default:
		return nil, fmt.Errorf("platform '%s' not supported at this time", runtime.GOOS)
	}
}

func getRequiredProgramsOSX() *Prerequisites {
	result := &Prerequisites{}
	result.Add(newPrerequisite("clang", "Please install with `xcode-select --install` and try again"))
	result.Add(newPrerequisite("npm", "Please install from https://nodejs.org/en/download/ and try again"))
	return result
}

func getRequiredProgramsLinux() *Prerequisites {
	result := &Prerequisites{}
	distroInfo := GetLinuxDistroInfo()
	if distroInfo.Distribution != Unknown {
		var linuxDB = NewLinuxDB()
		distro := linuxDB.GetDistro(distroInfo.ID)
		release := distro.GetRelease(distroInfo.Release)
		for _, program := range release.Programs {
			result.Add(program)
		}
	}
	return result
}

// TODO: Test this on Windows
func getRequiredProgramsWindows() *Prerequisites {
	result := &Prerequisites{}
	result.Add(newPrerequisite("gcc", "Please install gcc from here and try again: http://tdm-gcc.tdragon.net/download. You will need to add the bin directory to your path, EG: C:\\TDM-GCC-64\\bin\\"))
	result.Add(newPrerequisite("npm", "Please install node/npm from here and try again: https://nodejs.org/en/download/"))
	return result
}

// GetRequiredLibraries returns a list of libraries (packages) required for the platform
func GetRequiredLibraries() (*Prerequisites, error) {
	switch runtime.GOOS {
	case "darwin":
		return getRequiredLibrariesOSX()
	case "linux":
		return getRequiredLibrariesLinux()
	case "windows":
		return getRequiredLibrariesWindows()
	default:
		return nil, fmt.Errorf("platform '%s' not supported at this time", runtime.GOOS)
	}
}

func getRequiredLibrariesOSX() (*Prerequisites, error) {
	result := &Prerequisites{}
	return result, nil
}

func getRequiredLibrariesLinux() (*Prerequisites, error) {
	result := &Prerequisites{}
	// The Linux Distribution DB
	distroInfo := GetLinuxDistroInfo()
	if distroInfo.Distribution != Unknown {
		var linuxDB = NewLinuxDB()
		distro := linuxDB.GetDistro(distroInfo.ID)
		release := distro.GetRelease(distroInfo.Release)
		for _, library := range release.Libraries {
			result.Add(library)
		}
	}
	return result, nil
}

func getRequiredLibrariesWindows() (*Prerequisites, error) {
	result := &Prerequisites{}
	return result, nil
}
