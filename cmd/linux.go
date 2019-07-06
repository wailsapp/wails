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
	// Ubuntu distribution
	Ubuntu
	// Arch linux distribution
	Arch
	// RedHat linux distribution
	RedHat
	// Debian distribution
	Debian
	// Gentoo distribution
	Gentoo
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	Distribution  LinuxDistribution
	Description   string
	Release       string
	Codename      string
	DistributorID string
	DiscoveredBy  string
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
		result.DiscoveredBy = "lsb"
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
					case "Arch", "ManjaroLinux":
						result.Distribution = Arch
					case "Debian":
						result.Distribution = Debian
					case "Gentoo":
						result.Distribution = Gentoo
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
		// check if /etc/os-release exists
	} else if _, err := os.Stat("/etc/os-release"); !os.IsNotExist(err) {
		// Default value
		osName := "Unknown"
		version := ""
		// read /etc/os-release
		osRelease, _ := ioutil.ReadFile("/etc/os-release")
		// Split into lines
		lines := strings.Split(string(osRelease), "\n")
		// Iterate lines
		for _, line := range lines {
			// Split each line by the equals char
			splitLine := strings.SplitN(line, "=", 2)
			// Check we have
			if len(splitLine) != 2 {
				continue
			}
			switch splitLine[0] {
			case "NAME":
				osName = strings.Trim(splitLine[1], "\"")
			case "VERSION_ID":
				version = strings.Trim(splitLine[1], "\"")
			}

		}
		// Check distro name against list of distros
		result.Release = version
		result.DiscoveredBy = "os-release"
		switch osName {
		case "Fedora":
			result.Distribution = RedHat
		case "CentOS":
			result.Distribution = RedHat
		case "Arch Linux":
			result.Distribution = Arch
		case "Debian GNU/Linux":
			result.Distribution = Debian
		case "Gentoo/Linux":
			result.Distribution = Gentoo
		default:
			result.Distribution = Unknown
			result.DistributorID = osName
		}
	}
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
	defaultError := fmt.Errorf("unable to check libraries on distribution '%s'. Please ensure that the '%s' equivalent is installed", distroInfo.DistributorID, libraryName)

	logger.Yellow("Distribution '%s' is not currently supported, but we would love to!", distroInfo.DistributorID)
	q := fmt.Sprintf("Would you like to submit a request to support distribution '%s'?", distroInfo.DistributorID)
	result := Prompt(q, "yes")
	if strings.ToLower(result) != "yes" {
		return defaultError
	}

	title := fmt.Sprintf("Support Distribution '%s'", distroInfo.DistributorID)

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
	str.WriteString(fmt.Sprintf("| Distribution ID   | %s |\n", distroInfo.DistributorID))
	str.WriteString(fmt.Sprintf("| Distribution Version   | %s |\n", distroInfo.Release))
	str.WriteString(fmt.Sprintf("| Discovered by   | %s |\n", distroInfo.DiscoveredBy))

	body := fmt.Sprintf("**Description**\nDistribution '%s' is currently unsupported.\n\n**Further Information**\n\n%s\n\n*Please add any extra information here, EG: libraries that are needed to make the distribution work, or commands to install them*", distroInfo.DistributorID, str.String())
	fullURL := "https://github.com/wailsapp/wails/issues/new?"
	params := "title=" + title + "&body=" + body

	fmt.Println("Opening browser to file request.")
	browser.OpenURL(fullURL + url.PathEscape(params))
	return nil

}
