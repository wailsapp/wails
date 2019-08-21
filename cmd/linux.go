package cmd

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/browser"
)

// LinuxDistribution is of type int
type LinuxDistribution int

const (
	// Unknown is the catch-all distro
	Unknown LinuxDistribution = iota
	// Debian distribution
	Debian
	// Ubuntu distribution
	Ubuntu
	// Parrot distribution
	Parrot
	// Zorin distribution
	Zorin
	// Linuxmint distribution
	Linuxmint
	// Elementary distribution
	Elementary
	// Centos linux distribution
	Centos
	// Fedora linux distribution
	Fedora
	// Arch linux distribution
	Arch
	// Gentoo distribution
	Gentoo
	// Opensuse distribution
	Opensuse
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	Distribution LinuxDistribution
	// ie. NAME="Ubuntu"
	Name string
	// ie. ID=ubuntu
	ID          string
	Description string
	// ie. VERSION_ID="18.04"
	Release string
}

// GetLinuxDistroInfo returns information about the running linux distribution
func GetLinuxDistroInfo() *DistroInfo {
	result := &DistroInfo{
		Distribution: Unknown,
		ID:           "unknown",
		Name:         "Unknown",
	}
	_, err := os.Stat("/etc/os-release")
	if !os.IsNotExist(err) {
		osRelease, _ := ioutil.ReadFile("/etc/os-release")
		result = parseOsRelease(string(osRelease))
	}
	return result
}

// parseOsRelease parses the given os-release data and returns
// a DistroInfo struct with the details
func parseOsRelease(osRelease string) *DistroInfo {
	result := &DistroInfo{Distribution: Unknown}

	// Default value
	osID := "unknown"
	osNAME := "Unknown"
	version := ""

	// Split into lines
	lines := strings.Split(osRelease, "\n")
	// Iterate lines
	for _, line := range lines {
		// Split each line by the equals char
		splitLine := strings.SplitN(line, "=", 2)
		// Check we have
		if len(splitLine) != 2 {
			continue
		}

		switch splitLine[0] {
		case "ID":
			osID = strings.Trim(splitLine[1], "\"")
		case "NAME":
			osNAME = strings.Trim(splitLine[1], "\"")
		// for debian, ubuntu, based distros
		case "VERSION_ID":
			version = strings.Trim(splitLine[1], "\"")
		// for arch
		case "BUILD_ID":
			version = strings.Trim(splitLine[1], "\"")
		}
	}
	// Check distro name against list of distros
	switch osID {
	case "fedora":
		result.Distribution = Fedora
	case "centos":
		result.Distribution = Centos
	case "arch":
		result.Distribution = Arch
	case "debian":
		result.Distribution = Debian
	case "ubuntu":
		result.Distribution = Ubuntu
	case "gentoo":
		result.Distribution = Gentoo
	case "zorin":
		result.Distribution = Zorin
	case "parrot":
		result.Distribution = Parrot
	case "linuxmint":
		result.Distribution = Linuxmint
	case "\"elementary\"":
		result.Distribution = Elementary
	case "\"opensuse-tumbleweed\"":
		result.Distribution = Opensuse
	default:
		result.Distribution = Unknown
	}

	result.Release = version
	result.ID = osID
	result.Name = osNAME
	return result
}

// EqueryInstalled uses equery to see if a package is installed
func EqueryInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	equery := program.FindProgram("equery")
	if equery == nil {
		return false, fmt.Errorf("cannont check dependencies: equery not found")
	}
	_, _, exitCode, _ := equery.Run("l", packageName)
	return exitCode == 0, nil
}

// DpkgInstalled uses dpkg to see if a package is installed
func DpkgInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	dpkg := program.FindProgram("dpkg")
	if dpkg == nil {
		return false, fmt.Errorf("cannot check dependencies: dpkg not found")
	}
	_, _, exitCode, _ := dpkg.Run("-L", packageName)
	return exitCode == 0, nil
}

// PacmanInstalled uses pacman to see if a package is installed.
func PacmanInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	pacman := program.FindProgram("pacman")
	if pacman == nil {
		return false, fmt.Errorf("cannot check dependencies: pacman not found")
	}
	_, _, exitCode, _ := pacman.Run("-Qs", packageName)
	return exitCode == 0, nil
}

// RpmInstalled uses rpm to see if a package is installed
func RpmInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	rpm := program.FindProgram("rpm")
	if rpm == nil {
		return false, fmt.Errorf("cannot check dependencies: rpm not found")
	}
	_, _, exitCode, _ := rpm.Run("--query", packageName)
	return exitCode == 0, nil
}

// RequestSupportForDistribution promts the user to submit a request to support their
// currently unsupported distribution
func RequestSupportForDistribution(distroInfo *DistroInfo, libraryName string) error {
	var logger = NewLogger()
	defaultError := fmt.Errorf("unable to check libraries on distribution '%s'. Please ensure that the '%s' equivalent is installed", distroInfo.Name, libraryName)

	logger.Yellow("Distribution '%s' is not currently supported, but we would love to!", distroInfo.Name)
	q := fmt.Sprintf("Would you like to submit a request to support distribution '%s'?", distroInfo.Name)
	result := Prompt(q, "yes")
	if strings.ToLower(result) != "yes" {
		return defaultError
	}

	title := fmt.Sprintf("Support Distribution '%s'", distroInfo.Name)

	var str strings.Builder

	gomodule, exists := os.LookupEnv("GO111MODULE")
	if !exists {
		gomodule = "(Not Set)"
	}

	str.WriteString("\n| Name   | Value |\n| ----- | ----- |\n")
	str.WriteString(fmt.Sprintf("| Wails Version | %s |\n", Version))
	str.WriteString(fmt.Sprintf("| Go Version    | %s |\n", runtime.Version()))
	str.WriteString(fmt.Sprintf("| Platform      | %s |\n", runtime.GOOS))
	str.WriteString(fmt.Sprintf("| Arch          | %s |\n", runtime.GOARCH))
	str.WriteString(fmt.Sprintf("| GO111MODULE   | %s |\n", gomodule))
	str.WriteString(fmt.Sprintf("| Distribution ID   | %s |\n", distroInfo.ID))
	str.WriteString(fmt.Sprintf("| Distribution Name   | %s |\n", distroInfo.Name))
	str.WriteString(fmt.Sprintf("| Distribution Version   | %s |\n", distroInfo.Release))

	body := fmt.Sprintf("**Description**\nDistribution '%s' is currently unsupported.\n\n**Further Information**\n\n%s\n\n*Please add any extra information here, EG: libraries that are needed to make the distribution work, or commands to install them*", distroInfo.ID, str.String())
	fullURL := "https://github.com/wailsapp/wails/issues/new?"
	params := "title=" + title + "&body=" + body

	fmt.Println("Opening browser to file request.")
	browser.OpenURL(fullURL + url.PathEscape(params))
	return nil
}
