//go:build linux

package setupwizard

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func (w *Wizard) checkAllDependencies() []DependencyStatus {
	var deps []DependencyStatus

	deps = append(deps, checkGo())

	info, _ := operatingsystem.Info()

	pm := packagemanager.Find(info.ID)
	if pm != nil {
		platformDeps, _ := packagemanager.Dependencies(pm)
		for _, dep := range platformDeps {
			// Skip npm from package manager - we'll check it via PATH instead
			if dep.Name == "npm" {
				continue
			}
			status := DependencyStatus{
				Name:     dep.Name,
				Required: !dep.Optional,
			}

			if dep.Installed {
				status.Installed = true
				status.Status = "installed"
				status.Version = dep.Version
			} else {
				// Also check if the binary is in PATH (might be installed via other means)
				if _, err := exec.LookPath(dep.Name); err == nil {
					status.Installed = true
					status.Status = "installed"
				} else {
					status.Installed = false
					status.Status = "not_installed"
					status.InstallCommand = dep.InstallCommand
				}
			}

			deps = append(deps, status)
		}
	}

	// Always check npm via PATH (might be installed via nvm, fnm, etc.)
	deps = append(deps, checkNpm())

	deps = append(deps, checkDocker())

	return deps
}

func checkGo() DependencyStatus {
	dep := DependencyStatus{
		Name:     "Go",
		Required: true,
	}

	version, err := execCommand("go", "version")
	if err != nil {
		dep.Status = "not_installed"
		dep.Installed = false
		dep.Message = "Go 1.25+ is required"
		dep.HelpURL = "https://go.dev/dl/"
		dep.InstallCommand = "Download from https://go.dev/dl/"
		return dep
	}

	dep.Installed = true
	dep.Status = "installed"

	parts := strings.Split(version, " ")
	if len(parts) >= 3 {
		versionStr := strings.TrimPrefix(parts[2], "go")
		dep.Version = versionStr

		versionParts := strings.Split(versionStr, ".")
		if len(versionParts) >= 2 {
			major, majorErr := strconv.Atoi(versionParts[0])
			// Handle versions like "25beta1" by extracting leading digits
			minorStr := versionParts[1]
			for i, c := range minorStr {
				if c < '0' || c > '9' {
					minorStr = minorStr[:i]
					break
				}
			}
			minor, minorErr := strconv.Atoi(minorStr)
			if majorErr != nil || minorErr != nil {
				// Couldn't parse version; assume it's acceptable
				return dep
			}
			if major < 1 || (major == 1 && minor < 25) {
				dep.Status = "needs_update"
				dep.Message = "Go 1.25+ is required (found " + versionStr + ")"
				dep.HelpURL = "https://go.dev/dl/"
			}
		}
	}

	return dep
}

func checkNpm() DependencyStatus {
	dep := DependencyStatus{
		Name:     "npm",
		Required: false, // Optional - not strictly required for Go-only projects
	}

	version, err := execCommand("npm", "-v")
	if err != nil {
		dep.Status = "not_installed"
		dep.Installed = false
		dep.Message = "Required for frontend development"
		dep.HelpURL = "https://nodejs.org/"
		dep.InstallCommand = "Install Node.js from https://nodejs.org/"
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
			dep.Message = "npm 7.0.0 or higher recommended"
			dep.HelpURL = "https://nodejs.org/"
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
		dep.Message = "Enables cross-platform builds"
		dep.HelpURL = "https://docs.docker.com/get-docker/"
		dep.InstallCommand = "Install Docker from https://docs.docker.com/get-docker/"
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
		dep.Message = "Start Docker to enable cross-compilation"
		return dep
	}

	// Check for wails-cross image
	// docker image inspect returns "[]" (empty JSON array) on stdout when image doesn't exist
	imageCheck, _ := execCommand("docker", "image", "inspect", crossImageName)
	if imageCheck == "" || imageCheck == "[]" || strings.Contains(imageCheck, "Error") {
		dep.Installed = true
		dep.Status = "installed"
		dep.ImageBuilt = false
		dep.Message = "Run 'wails3 task setup:docker' to build cross-compilation image"
	} else {
		dep.Installed = true
		dep.Status = "installed"
		dep.ImageBuilt = true
		dep.Message = "Cross-compilation ready"
	}

	return dep
}
