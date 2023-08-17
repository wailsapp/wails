//go:build darwin

package operatingsystem

import (
	"os/exec"
	"strings"
)

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

	// 		cmd := CreateCommand(directory, command, args...)
	// 		var stdo, stde bytes.Buffer
	// 		cmd.Stdout = &stdo
	// 		cmd.Stderr = &stde
	// 		err := cmd.Run()
	// 		return stdo.String(), stde.String(), err
	// 	}
	// 	sysctl := shell.NewCommand("sysctl")
	// 	kern.ostype: Darwin
	// kern.osrelease: 20.1.0
	// kern.osrevision: 199506

	return &result, nil
}
