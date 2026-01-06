//go:build darwin

package doctorng

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"

	"github.com/samber/lo"
)

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	appleSilicon := "unknown"
	r, err := syscall.Sysctl("sysctl.proc_translated")
	if err == nil {
		appleSilicon = lo.Ternary(r == "\x00\x00\x00" || r == "\x01\x00\x00", "true", "false")
	}
	extras["Apple Silicon"] = appleSilicon

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
		Name:           "NSIS",
		Version:        nsisVersion,
		Status:         nsisStatus,
		Required:       false,
		InstallCommand: "brew install makensis",
		Category:       "optional",
		Description:    "For Windows installer generation",
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
		Name:           "npm",
		Version:        npmVersion,
		Status:         npmStatus,
		Required:       true,
		InstallCommand: "brew install node",
		Category:       "frontend",
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
		Name:           "docker",
		Version:        dockerVersion,
		Status:         dockerStatus,
		Required:       false,
		InstallCommand: "brew install --cask docker",
		Category:       "optional",
		Description:    "For cross-compilation",
	})
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkMacSpecific()
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
				Command:     "brew install go",
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
