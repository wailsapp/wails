package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

// Info holds information about the current operating system,
// package manager and required dependancies
type Info struct {
	OS           *operatingsystem.OS
	PM           packagemanager.PackageManager
	Dependencies packagemanager.DependencyList
}

// GetInfo scans the system for operating system details,
// the system package manager and the status of required
// dependancies.
func GetInfo() (*Info, error) {
	var result Info
	err := result.discover()
	if err != nil {
		return nil, err
	}
	return &result, nil
}
