package application

import "github.com/wailsapp/wails/v3/internal/operatingsystem"

// EnvironmentInfo represents information about the current environment.
//
// Fields:
// - OS: the operating system that the program is running on.
// - Arch: the architecture of the operating system.
// - Debug: indicates whether debug mode is enabled.
// - OSInfo: information about the operating system.
type EnvironmentInfo struct {
	OS           string
	Arch         string
	Debug        bool
	OSInfo       *operatingsystem.OS
	PlatformInfo map[string]any
}
