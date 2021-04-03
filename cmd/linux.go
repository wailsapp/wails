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
	// Arch linux distribution
	Arch
	// CentOS linux distribution
	CentOS
	// Fedora linux distribution
	Fedora
	// Gentoo distribution
	Gentoo
	// Zorin distribution
	Zorin
	// Parrot distribution
	Parrot
	// Linuxmint distribution
	Linuxmint
	// VoidLinux distribution
	VoidLinux
	// Elementary distribution
	Elementary
	// Kali distribution
	Kali
	// Neon distribution
	Neon
	// ArcoLinux distribution
	ArcoLinux
	// Manjaro distribution
	Manjaro
	// ManjaroARM distribution
	ManjaroARM
	// Deepin distribution
	Deepin
	// Raspbian distribution
	Raspbian
	// Tumbleweed (OpenSUSE) distribution
	Tumbleweed
	// Leap (OpenSUSE) distribution
	Leap
	// ArchLabs distribution
	ArchLabs
	// PopOS distribution
	PopOS
	// Solus distribution
	Solus
	// Ctlos Linux distribution
	Ctlos
	// EndeavourOS linux distribution
	EndeavourOS
	// Crux linux distribution
	Crux
	// RHEL distribution
	RHEL
)

// DistroInfo contains all the information relating to a linux distribution
type DistroInfo struct {
	Distribution LinuxDistribution
	Name         string
	ID           string
	Description  string
	Release      string
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
			osID = strings.ToLower(strings.Trim(splitLine[1], "\""))
		case "NAME":
			osNAME = strings.Trim(splitLine[1], "\"")
		case "VERSION_ID":
			version = strings.Trim(splitLine[1], "\"")
		}
	}

	// Check distro name against list of distros
	switch osID {
	case "fedora":
		result.Distribution = Fedora
	case "centos":
		result.Distribution = CentOS
	case "rhel":
		result.Distribution = RHEL
	case "arch":
		result.Distribution = Arch
	case "archlabs":
		result.Distribution = ArchLabs
	case "ctlos":
		result.Distribution = Ctlos
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
	case "void":
		result.Distribution = VoidLinux
	case "elementary":
		result.Distribution = Elementary
	case "kali":
		result.Distribution = Kali
	case "neon":
		result.Distribution = Neon
	case "arcolinux":
		result.Distribution = ArcoLinux
	case "manjaro":
		result.Distribution = Manjaro
	case "manjaro-arm":
		result.Distribution = ManjaroARM
	case "deepin":
		result.Distribution = Deepin
	case "raspbian":
		result.Distribution = Raspbian
	case "opensuse-tumbleweed":
		result.Distribution = Tumbleweed
	case "opensuse-leap":
		result.Distribution = Leap
	case "pop":
		result.Distribution = PopOS
	case "solus":
		result.Distribution = Solus
	case "endeavouros":
		result.Distribution = EndeavourOS
	case "crux":
		result.Distribution = Crux
	default:
		result.Distribution = Unknown
	}

	result.Name = osNAME
	result.ID = osID
	result.Release = version

	return result
}

// CheckPkgInstalled is all functions that use local programs to see if a package is installed
type CheckPkgInstalled func(string) (bool, error)

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

// EOpkgInstalled uses dpkg to see if a package is installed
func EOpkgInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	eopkg := program.FindProgram("eopkg")
	if eopkg == nil {
		return false, fmt.Errorf("cannot check dependencies: eopkg not found")
	}
	stdout, _, _, _ := eopkg.Run("info", packageName)
	return strings.HasPrefix(stdout, "Installed"), nil
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

// XbpsInstalled uses pacman to see if a package is installed.
func XbpsInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	xbpsQuery := program.FindProgram("xbps-query")
	if xbpsQuery == nil {
		return false, fmt.Errorf("cannot check dependencies: xbps-query not found")
	}
	_, _, exitCode, _ := xbpsQuery.Run("-S", packageName)
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

// PrtGetInstalled uses prt-get to see if a package is installed
func PrtGetInstalled(packageName string) (bool, error) {
	program := NewProgramHelper()
	prtget := program.FindProgram("prt-get")
	if prtget == nil {
		return false, fmt.Errorf("cannot check dependencies: prt-get not found")
	}
	_, _, exitCode, _ := prtget.Run("isinst", packageName)
	return exitCode == 0, nil
}

// RequestSupportForDistribution promts the user to submit a request to support their
// currently unsupported distribution
func RequestSupportForDistribution(distroInfo *DistroInfo) error {
	var logger = NewLogger()
	defaultError := fmt.Errorf("unable to check libraries on distribution '%s'", distroInfo.Name)

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
	result = Prompt("We have a guide for adding support for your distribution. Would you like to view it?", "yes")
	if strings.ToLower(result) == "yes" {
		browser.OpenURL("https://wails.app/guides/distro/")
	}
	return nil
}
