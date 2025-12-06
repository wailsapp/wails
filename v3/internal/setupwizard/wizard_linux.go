//go:build linux

package setupwizard

import (
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func (w *Wizard) checkAllDependencies() []DependencyStatus {
	var deps []DependencyStatus

	// Get OS info for package manager detection
	info, _ := operatingsystem.Info()

	// Find the package manager
	pm := packagemanager.Find(info.ID)
	if pm != nil {
		// Get platform dependencies from the doctor package
		platformDeps, _ := packagemanager.Dependencies(pm)
		for _, dep := range platformDeps {
			status := DependencyStatus{
				Name:     dep.Name,
				Required: !dep.Optional,
			}

			if dep.Installed {
				status.Installed = true
				status.Status = "installed"
				status.Version = dep.Version
			} else {
				status.Installed = false
				status.Status = "not_installed"
				if dep.InstallCommand != "" {
					status.Message = "Install with: " + dep.InstallCommand
				}
			}

			deps = append(deps, status)
		}
	}

	// Check npm (common dependency)
	deps = append(deps, checkNpm())

	// Check Docker (optional)
	deps = append(deps, checkDocker())

	return deps
}

func checkNpm() DependencyStatus {
	dep := DependencyStatus{
		Name:     "npm",
		Required: true,
	}

	version, err := execCommand("npm", "-v")
	if err != nil {
		dep.Status = "not_installed"
		dep.Installed = false
		dep.Message = "npm is required. Install Node.js from https://nodejs.org/"
		return dep
	}

	dep.Version = version

	// Check minimum version (7.0.0)
	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		major, _ := strconv.Atoi(parts[0])
		if major < 7 {
			dep.Status = "needs_update"
			dep.Installed = true
			dep.Message = "npm 7.0.0 or higher is required"
			return dep
		}
	}

	dep.Installed = true
	dep.Status = "installed"
	return dep
}

func checkDocker() DependencyStatus {
	dep := DependencyStatus{
		Name:     "docker",
		Required: false, // Optional for cross-compilation
	}

	version, err := execCommand("docker", "--version")
	if err != nil {
		dep.Status = "not_installed"
		dep.Installed = false
		dep.Message = "Optional - for cross-compilation"
		return dep
	}

	// Parse version from "Docker version 24.0.7, build afdd53b"
	parts := strings.Split(version, ",")
	if len(parts) > 0 {
		dep.Version = strings.TrimPrefix(strings.TrimSpace(parts[0]), "Docker version ")
	}

	// Check if daemon is running
	_, err = execCommand("docker", "info")
	if err != nil {
		dep.Installed = true
		dep.Status = "installed"
		dep.Message = "Daemon not running"
		return dep
	}

	// Check for wails-cross image
	imageCheck, _ := execCommand("docker", "image", "inspect", "wails-cross")
	if imageCheck == "" || strings.Contains(imageCheck, "Error") {
		dep.Installed = true
		dep.Status = "installed"
		dep.Message = "wails-cross image not built"
	} else {
		dep.Installed = true
		dep.Status = "installed"
		dep.Message = "Cross-compilation ready"
	}

	return dep
}
