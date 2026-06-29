//go:build linux

package setupwizard

import (
	"os"
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
				// For gcc, the package manager reports the build-essential meta-package
				// version, not the actual gcc binary version. Use gcc --version instead.
				if dep.Name == "gcc" {
					if gccOut, err := execCommand("gcc", "--version"); err == nil {
						firstLine := strings.SplitN(gccOut, "\n", 2)[0]
						if idx := strings.LastIndex(firstLine, ")"); idx != -1 {
							if ver := strings.TrimSpace(firstLine[idx+1:]); ver != "" {
								status.Version = ver
							}
						}
					}
				}
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
				// System go is < 1.25. Check whether GOTOOLCHAIN provides a newer version
				// automatically. If GOTOOLCHAIN is not "off", Go will download and use
				// 1.25+ when a module requires it, so the user needs no manual action.
				if tcVer := findToolchainGo(); tcVer != "" {
					dep.Status = "installed"
					dep.Message = "Go 1.25 available via GOTOOLCHAIN (" + tcVer + ")"
				} else {
					dep.Status = "needs_update"
					dep.Message = "Go 1.25+ is required (found " + versionStr + ")"
					dep.HelpURL = "https://go.dev/dl/"
				}
			}
		}
	}

	return dep
}

// findToolchainGo returns a Go 1.25+ version string if one is available via
// GOTOOLCHAIN, or an empty string if the system go is the only option.
func findToolchainGo() string {
	gotoolchain, _ := execCommand("go", "env", "GOTOOLCHAIN")
	if gotoolchain == "off" {
		return ""
	}

	// GOTOOLCHAIN may name a specific version: "go1.25.0" or "go1.25.0+auto".
	tc := strings.TrimSuffix(gotoolchain, "+auto")
	tc = strings.TrimSuffix(tc, "+path")
	if strings.HasPrefix(tc, "go") {
		ver := strings.TrimPrefix(tc, "go")
		if isGoMinVersion(ver, 1, 25) {
			return ver
		}
	}

	// GOTOOLCHAIN=auto|path: scan the module cache for a downloaded toolchain.
	// Use GOMODCACHE rather than GOPATH so that a custom GOMODCACHE env var and
	// multi-entry GOPATH values are handled correctly.
	gomodcache, err := execCommand("go", "env", "GOMODCACHE")
	if err != nil || gomodcache == "" {
		return ""
	}
	entries, err := os.ReadDir(gomodcache + "/golang.org")
	if err != nil {
		return ""
	}
	for _, e := range entries {
		// Directory names follow the pattern:
		//   toolchain@v<modver>-go<goversion>.<goos>-<goarch>
		// e.g. "toolchain@v0.0.1-go1.25.0.linux-amd64"
		// Match any module version to stay forward-compatible with future
		// changes to the golang.org/toolchain module versioning.
		name := e.Name()
		if !strings.HasPrefix(name, "toolchain@v") {
			continue
		}
		goIdx := strings.Index(name, "-go")
		if goIdx == -1 {
			continue
		}
		// goVerAndPlatform is e.g. "1.25.0.linux-amd64"
		goVerAndPlatform := name[goIdx+3:]
		parts := strings.SplitN(goVerAndPlatform, ".", 3)
		if len(parts) < 2 {
			continue
		}
		minorStr := parts[1]
		for i, c := range minorStr {
			if c < '0' || c > '9' {
				minorStr = minorStr[:i]
				break
			}
		}
		if isGoMinVersion(parts[0]+"."+minorStr, 1, 25) {
			if len(parts) >= 3 {
				// Third segment is "0.linux-amd64"; extract leading digits only.
				patch := parts[2]
				for i, c := range patch {
					if c < '0' || c > '9' {
						patch = patch[:i]
						break
					}
				}
				if patch != "" {
					return parts[0] + "." + minorStr + "." + patch
				}
			}
			return parts[0] + "." + minorStr
		}
	}
	return ""
}

// isGoMinVersion reports whether the version string (e.g. "1.25") meets the
// given minimum major.minor requirement.
func isGoMinVersion(ver string, minMajor, minMinor int) bool {
	parts := strings.SplitN(ver, ".", 2)
	if len(parts) < 2 {
		return false
	}
	major, err1 := strconv.Atoi(parts[0])
	minorStr := parts[1]
	for i, c := range minorStr {
		if c < '0' || c > '9' {
			minorStr = minorStr[:i]
			break
		}
	}
	minor, err2 := strconv.Atoi(minorStr)
	if err1 != nil || err2 != nil {
		return false
	}
	return major > minMajor || (major == minMajor && minor >= minMinor)
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
