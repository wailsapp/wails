//go:build windows

package setupwizard

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// refreshPath updates the process PATH environment variable from the Windows registry.
// This is needed because when software is installed, the PATH is updated in the registry
// but running processes still have the old PATH until they restart.
func refreshPath() {
	var paths []string

	// Get system PATH from HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment
	if key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE); err == nil {
		if systemPath, _, err := key.GetStringValue("Path"); err == nil {
			paths = append(paths, strings.Split(systemPath, ";")...)
		}
		key.Close()
	}

	// Get user PATH from HKEY_CURRENT_USER\Environment
	if key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE); err == nil {
		if userPath, _, err := key.GetStringValue("Path"); err == nil {
			paths = append(paths, strings.Split(userPath, ";")...)
		}
		key.Close()
	}

	// Build new PATH, removing empty entries
	var cleanPaths []string
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" {
			cleanPaths = append(cleanPaths, p)
		}
	}

	if len(cleanPaths) > 0 {
		os.Setenv("PATH", strings.Join(cleanPaths, ";"))
	}
}

// execCommandRefreshed refreshes PATH and then executes a command.
// This ensures newly installed software is found.
func execCommandRefreshed(name string, args ...string) (string, error) {
	refreshPath()
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func (w *Wizard) checkAllDependencies() []DependencyStatus {
	var deps []DependencyStatus

	// Check WebView2 Runtime
	deps = append(deps, checkWebView2())

	// Check npm (common dependency)
	deps = append(deps, checkNpm())

	// Check Docker (optional)
	deps = append(deps, checkDocker())

	return deps
}

func checkWebView2() DependencyStatus {
	dep := DependencyStatus{
		Name:     "WebView2 Runtime",
		Required: true,
	}

	// Check common installation paths
	paths := []string{
		filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Microsoft", "EdgeWebView", "Application"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "EdgeWebView", "Application"),
		filepath.Join(os.Getenv("PROGRAMFILES"), "Microsoft", "EdgeWebView", "Application"),
	}

	for _, path := range paths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			dep.Installed = true
			dep.Status = "installed"

			// Try to get version from directory name
			entries, _ := os.ReadDir(path)
			for _, entry := range entries {
				if entry.IsDir() {
					name := entry.Name()
					// Version directories look like "120.0.2210.91"
					if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
						dep.Version = name
						break
					}
				}
			}
			return dep
		}
	}

	dep.Status = "not_installed"
	dep.Installed = false
	dep.Message = "Required for rendering the application UI"
	dep.HelpURL = "https://developer.microsoft.com/en-us/microsoft-edge/webview2/"
	return dep
}

func checkNpm() DependencyStatus {
	dep := DependencyStatus{
		Name:     "npm",
		Required: true,
	}

	version, err := execCommandRefreshed("npm", "-v")
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

	version, err := execCommandRefreshed("docker", "--version")
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
	imageCheck, _ := execCommand("docker", "image", "inspect", "wails-cross")
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
