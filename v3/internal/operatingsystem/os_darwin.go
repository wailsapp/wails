//go:build darwin

package operatingsystem

import (
	"os/exec"
	"strings"
)

var macOSNames = map[string]string{
	"10.10": "Yosemite",
	"10.11": "El Capitan",
	"10.12": "Sierra",
	"10.13": "High Sierra",
	"10.14": "Mojave",
	"10.15": "Catalina",
	"11":    "Big Sur",
	"12":    "Monterey",
	"13":    "Ventura",
	"14":    "Sonoma",
	"15":    "Sequoia",
	// Add newer versions as they are released...
}

func getOSName(version string) string {
	trimmedVersion := version
	if !strings.HasPrefix(version, "10.") {
		trimmedVersion = strings.SplitN(version, ".", 2)[0]
	}
	name, ok := macOSNames[trimmedVersion]
	if ok {
		return name
	}
	return "MacOS " + version
}

func getSysctlValue(key string) (string, error) {
	// Run "sysctl" command
	command := exec.Command("sysctl", key)
	// Capture stdout
	var stdout strings.Builder
	command.Stdout = &stdout
	// Run command
	err := command.Run()
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(stdout.String(), key+": ")
	return strings.TrimSpace(version), nil
}

func platformInfo() (*OS, error) {
	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "MacOS"
	result.Version = "Unknown"

	version, err := getSysctlValue("kern.osproductversion")
	if err != nil {
		return nil, err
	}
	result.Version = version
	ID, err := getSysctlValue("kern.osversion")
	if err != nil {
		return nil, err
	}
	result.ID = ID
	result.Branding = getOSName(result.Version)

	return &result, nil
}
