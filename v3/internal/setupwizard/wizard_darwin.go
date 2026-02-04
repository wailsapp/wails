//go:build darwin

package setupwizard

import (
	"os/exec"
	"strconv"
	"strings"
)

func (w *Wizard) checkAllDependencies() []DependencyStatus {
	var deps []DependencyStatus

	// Check Go (required)
	deps = append(deps, checkGo())

	// Check Xcode Command Line Tools
	deps = append(deps, checkXcode())

	// Check npm (common dependency)
	deps = append(deps, checkNpm())

	// Check Docker (optional)
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

func checkXcode() DependencyStatus {
	dep := DependencyStatus{
		Name:     "Xcode Command Line Tools",
		Required: true,
	}

	path, err := execCommand("xcode-select", "-p")
	if err != nil {
		dep.Status = "not_installed"
		dep.Installed = false
		dep.Message = "Run: xcode-select --install"
		return dep
	}

	dep.Installed = true
	dep.Status = "installed"

	// Try to get version
	cmd := exec.Command("pkgutil", "--pkg-info=com.apple.pkg.CLTools_Executables")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "version:") {
				dep.Version = strings.TrimSpace(strings.TrimPrefix(line, "version:"))
				break
			}
		}
	}

	_ = path // suppress unused warning
	return dep
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
		dep.Message = "Required for frontend development"
		dep.HelpURL = "https://nodejs.org/"
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
	// docker image inspect returns "[]" (empty JSON array) on stdout when image doesn't exist
	imageCheck, _ := execCommand("docker", "image", "inspect", crossImageName)
	if imageCheck == "" || imageCheck == "[]" || strings.Contains(imageCheck, "Error") {
		dep.Installed = true
		dep.Status = "installed"
		dep.ImageBuilt = false
		dep.Message = "wails-cross image not built"
	} else {
		dep.Installed = true
		dep.Status = "installed"
		dep.ImageBuilt = true
		dep.Message = "Cross-compilation ready"
	}

	return dep
}
