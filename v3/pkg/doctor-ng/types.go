// Package doctorng provides system diagnostics and dependency checking for Wails.
// It exposes a public API suitable for both CLI and GUI consumption.
package doctorng

import (
	"encoding/json"
	"time"
)

// Status represents the health status of a check
type Status int

const (
	StatusUnknown Status = iota
	StatusOK
	StatusWarning
	StatusError
	StatusMissing
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusWarning:
		return "warning"
	case StatusError:
		return "error"
	case StatusMissing:
		return "missing"
	default:
		return "unknown"
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// SystemInfo contains information about the host system
type SystemInfo struct {
	// Operating system details
	OS OSInfo `json:"os"`

	// Hardware information
	Hardware HardwareInfo `json:"hardware"`

	// Environment variables relevant to Wails
	Environment map[string]string `json:"environment"`

	// Platform-specific extras (e.g., XDG_SESSION_TYPE on Linux)
	PlatformExtras map[string]string `json:"platform_extras,omitempty"`
}

// OSInfo contains operating system details
type OSInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ID       string `json:"id"`
	Branding string `json:"branding,omitempty"`
	Platform string `json:"platform"` // darwin, linux, windows
	Arch     string `json:"arch"`     // amd64, arm64, etc.
}

// HardwareInfo contains hardware details
type HardwareInfo struct {
	CPUs   []CPUInfo `json:"cpus"`
	GPUs   []GPUInfo `json:"gpus"`
	Memory string    `json:"memory"`
}

// CPUInfo contains CPU details
type CPUInfo struct {
	Model string `json:"model"`
	Cores int    `json:"cores,omitempty"`
}

// GPUInfo contains GPU details
type GPUInfo struct {
	Name   string `json:"name"`
	Vendor string `json:"vendor,omitempty"`
	Driver string `json:"driver,omitempty"`
}

// BuildInfo contains build environment information
type BuildInfo struct {
	WailsVersion string            `json:"wails_version"`
	GoVersion    string            `json:"go_version"`
	BuildMode    string            `json:"build_mode,omitempty"`
	Compiler     string            `json:"compiler,omitempty"`
	CGOEnabled   bool              `json:"cgo_enabled"`
	Settings     map[string]string `json:"settings,omitempty"`
}

// Dependency represents a system dependency
type Dependency struct {
	// Name is the display name for this dependency
	Name string `json:"name"`

	// PackageName is the actual package name in the package manager
	PackageName string `json:"package_name,omitempty"`

	// Version is the installed version (empty if not installed)
	Version string `json:"version,omitempty"`

	// Status indicates the installation status
	Status Status `json:"status"`

	// Required indicates if this dependency is required (vs optional)
	Required bool `json:"required"`

	// InstallCommand is the command to install this dependency
	InstallCommand string `json:"install_command,omitempty"`

	// Description provides context about what this dependency is for
	Description string `json:"description,omitempty"`

	// Category groups related dependencies (e.g., "gtk", "build-tools")
	Category string `json:"category,omitempty"`
}

// DependencyList is a collection of dependencies with helper methods
type DependencyList []*Dependency

// RequiredMissing returns all required dependencies that are missing
func (d DependencyList) RequiredMissing() DependencyList {
	var result DependencyList
	for _, dep := range d {
		if dep.Required && dep.Status != StatusOK {
			result = append(result, dep)
		}
	}
	return result
}

// OptionalMissing returns all optional dependencies that are missing
func (d DependencyList) OptionalMissing() DependencyList {
	var result DependencyList
	for _, dep := range d {
		if !dep.Required && dep.Status != StatusOK {
			result = append(result, dep)
		}
	}
	return result
}

// AllInstalled returns true if all required dependencies are installed
func (d DependencyList) AllInstalled() bool {
	for _, dep := range d {
		if dep.Required && dep.Status != StatusOK {
			return false
		}
	}
	return true
}

// ByCategory groups dependencies by their category
func (d DependencyList) ByCategory() map[string]DependencyList {
	result := make(map[string]DependencyList)
	for _, dep := range d {
		cat := dep.Category
		if cat == "" {
			cat = "other"
		}
		result[cat] = append(result[cat], dep)
	}
	return result
}

// InstallCommands returns install commands for all missing dependencies
func (d DependencyList) InstallCommands(requiredOnly bool) []string {
	var commands []string
	for _, dep := range d {
		if dep.Status != StatusOK && dep.InstallCommand != "" {
			if requiredOnly && !dep.Required {
				continue
			}
			commands = append(commands, dep.InstallCommand)
		}
	}
	return commands
}

// DiagnosticSeverity indicates the severity of a diagnostic issue
type DiagnosticSeverity int

const (
	SeverityInfo DiagnosticSeverity = iota
	SeverityWarning
	SeverityError
)

func (s DiagnosticSeverity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}

// DiagnosticResult represents the result of a diagnostic check
type DiagnosticResult struct {
	// Name is a short identifier for this diagnostic
	Name string `json:"name"`

	// Message describes the issue or status
	Message string `json:"message"`

	// Severity indicates how serious the issue is
	Severity DiagnosticSeverity `json:"severity"`

	// HelpURL points to documentation about this issue
	HelpURL string `json:"help_url,omitempty"`

	// Fix contains instructions or a command to fix the issue
	Fix *Fix `json:"fix,omitempty"`
}

// Fix describes how to fix an issue
type Fix struct {
	// Description explains what the fix does
	Description string `json:"description"`

	// Command is a shell command that can be run to fix the issue
	// May be empty if manual intervention is required
	Command string `json:"command,omitempty"`

	// RequiresSudo indicates if the fix requires elevated privileges
	RequiresSudo bool `json:"requires_sudo,omitempty"`

	// ManualSteps are human-readable instructions if no command is available
	ManualSteps []string `json:"manual_steps,omitempty"`
}

// Report is the complete doctor report
type Report struct {
	// Timestamp when the report was generated
	Timestamp time.Time `json:"timestamp"`

	// System information
	System SystemInfo `json:"system"`

	// Build environment
	Build BuildInfo `json:"build"`

	// Dependencies and their status
	Dependencies DependencyList `json:"dependencies"`

	// Diagnostic results (issues found)
	Diagnostics []DiagnosticResult `json:"diagnostics"`

	// Overall status
	Ready bool `json:"ready"`

	// Summary message
	Summary string `json:"summary"`
}

// NewReport creates a new empty report with the current timestamp
func NewReport() *Report {
	return &Report{
		Timestamp:    time.Now(),
		Dependencies: make(DependencyList, 0),
		Diagnostics:  make([]DiagnosticResult, 0),
	}
}
