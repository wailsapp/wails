//go:build darwin

package doctorng

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"

	"github.com/samber/lo"
)

type macPackageManager int

const (
	macPMNone macPackageManager = iota
	macPMBrew
	macPMMacPorts
	macPMNix
)

var detectedMacPM macPackageManager
var macPMDetected bool

func detectMacPackageManager() macPackageManager {
	if macPMDetected {
		return detectedMacPM
	}
	macPMDetected = true

	if _, err := exec.LookPath("brew"); err == nil {
		detectedMacPM = macPMBrew
		return detectedMacPM
	}
	if _, err := exec.LookPath("port"); err == nil {
		detectedMacPM = macPMMacPorts
		return detectedMacPM
	}
	if _, err := exec.LookPath("nix-env"); err == nil {
		detectedMacPM = macPMNix
		return detectedMacPM
	}

	detectedMacPM = macPMNone
	return detectedMacPM
}

func macInstallCmd(brew, macports, nix, manual string) string {
	switch detectMacPackageManager() {
	case macPMBrew:
		return brew
	case macPMMacPorts:
		if macports != "" {
			return "sudo " + macports
		}
		return manual
	case macPMNix:
		if nix != "" {
			return nix
		}
		return manual
	default:
		return manual
	}
}

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	appleSilicon := "unknown"
	r, err := syscall.Sysctl("sysctl.proc_translated")
	if err == nil {
		appleSilicon = lo.Ternary(r == "\x00\x00\x00" || r == "\x01\x00\x00", "true", "false")
	}
	extras["Apple Silicon"] = appleSilicon

	pm := "none"
	switch detectMacPackageManager() {
	case macPMBrew:
		pm = "homebrew"
	case macPMMacPorts:
		pm = "macports"
	case macPMNix:
		pm = "nix"
	}
	extras["Package Manager"] = pm

	return extras
}

func (d *Doctor) collectDependencies() error {
	output, err := exec.Command("xcode-select", "-v").Output()
	xcodeStatus := StatusMissing
	xcodeVersion := ""
	if err == nil {
		xcodeStatus = StatusOK
		xcodeVersion = strings.TrimPrefix(string(output), "xcode-select version ")
		xcodeVersion = strings.TrimSpace(xcodeVersion)
		xcodeVersion = strings.TrimSuffix(xcodeVersion, ".")
	}

	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:           "Xcode CLI Tools",
		Version:        xcodeVersion,
		Status:         xcodeStatus,
		Required:       true,
		InstallCommand: "xcode-select --install",
		Category:       "build-tools",
	})

	d.checkCommonDependencies()

	nsisVersion := ""
	nsisStatus := StatusMissing
	output, err = exec.Command("makensis", "-VERSION").Output()
	if err == nil && output != nil {
		nsisStatus = StatusOK
		nsisVersion = strings.TrimSpace(string(output))
	}

	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "NSIS",
		Version: nsisVersion,
		Status:  nsisStatus,
		InstallCommand: macInstallCmd(
			"brew install makensis",
			"port install nsis",
			"nix-env -iA nixpkgs.nsis",
			"Download from https://nsis.sourceforge.io/",
		),
		Required:    false,
		Category:    "optional",
		Description: "For Windows installer generation",
	})

	return nil
}

func (d *Doctor) checkCommonDependencies() {
	npmVersion := ""
	npmStatus := StatusMissing
	output, err := exec.Command("npm", "--version").Output()
	if err == nil {
		npmStatus = StatusOK
		npmVersion = strings.TrimSpace(string(output))
	}

	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "npm",
		Version: npmVersion,
		Status:  npmStatus,
		InstallCommand: macInstallCmd(
			"brew install node",
			"port install nodejs18",
			"nix-env -iA nixpkgs.nodejs",
			"Download from https://nodejs.org/",
		),
		Required: true,
		Category: "frontend",
	})

	dockerVersion := ""
	dockerStatus := StatusMissing
	output, err = exec.Command("docker", "--version").Output()
	if err == nil {
		dockerStatus = StatusOK
		dockerVersion = strings.TrimSpace(string(output))
		output = bytes.Replace(output, []byte("Docker version "), []byte(""), 1)
		dockerVersion = strings.TrimSpace(string(output))
	}

	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "docker",
		Version: dockerVersion,
		Status:  dockerStatus,
		InstallCommand: macInstallCmd(
			"brew install --cask docker",
			"",
			"",
			"Download from https://docker.com/",
		),
		Required:    false,
		Category:    "optional",
		Description: "For cross-compilation",
	})
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkMacSpecific()
	d.checkPackageManager()
}

func (d *Doctor) checkGoInstallation() {
	if d.report.Build.GoVersion == "" {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Go Installation",
			Message:  "Go installation not found",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/",
			Fix: &Fix{
				Description: "Install Go",
				Command: macInstallCmd(
					"brew install go",
					"port install go",
					"nix-env -iA nixpkgs.go",
					"Download from https://go.dev/dl/",
				),
			},
		})
	}
}

func (d *Doctor) checkMacSpecific() {
	matches, _ := exec.Command("sh", "-c", "ls *.syso 2>/dev/null").Output()
	if len(matches) > 0 {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     ".syso files found",
			Message:  "Found .syso file(s) which may cause issues on macOS",
			Severity: SeverityWarning,
			HelpURL:  "/troubleshooting/mac-syso",
			Fix: &Fix{
				Description: "Remove .syso files before building on macOS",
				Command:     "rm *.syso",
			},
		})
	}
}

func (d *Doctor) checkPackageManager() {
	if detectMacPackageManager() == macPMNone {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Package Manager",
			Message:  "No package manager found (homebrew, macports, or nix)",
			Severity: SeverityWarning,
			HelpURL:  "/getting-started/installation/#macos",
			Fix: &Fix{
				Description: "Install Homebrew for easier dependency management",
				Command:     `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
			},
		})
	}
}
