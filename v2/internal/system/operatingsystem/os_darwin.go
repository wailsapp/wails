package operatingsystem

import (
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

func getSysctlValue(key string) (string, error) {
	stdout, _, err := shell.RunCommand(".", "sysctl", key)
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(stdout, key+": ")
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
