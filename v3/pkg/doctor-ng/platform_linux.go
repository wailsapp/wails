//go:build linux

package doctorng

import (
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/doctor-ng/packagemanager"
)

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	extras["XDG_SESSION_TYPE"] = getEnvOrDefault("XDG_SESSION_TYPE", "unset")
	extras["Desktop Environment"] = getEnvOrDefault("XDG_CURRENT_DESKTOP", "unset")
	extras["NVIDIA Driver"] = getNvidiaDriverInfo()

	return extras
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getNvidiaDriverInfo() string {
	version, err := os.ReadFile("/sys/module/nvidia/version")
	if err != nil {
		return "N/A"
	}

	versionStr := strings.TrimSpace(string(version))

	srcVersion, err := os.ReadFile("/sys/module/nvidia/srcversion")
	if err != nil {
		return versionStr
	}

	return versionStr + " (" + strings.TrimSpace(string(srcVersion)) + ")"
}

func (d *Doctor) collectDependencies() error {
	info, _ := operatingsystem.Info()
	pm := packagemanager.Find(info.ID)
	if pm == nil {
		return nil
	}

	deps, err := packagemanager.Dependencies(pm)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		status := StatusMissing
		if dep.Installed {
			status = StatusOK
		}

		d.report.Dependencies = append(d.report.Dependencies, &Dependency{
			Name:           dep.Name,
			PackageName:    dep.PackageName,
			Version:        dep.Version,
			Status:         status,
			Required:       !dep.Optional,
			InstallCommand: dep.InstallCommand,
			Category:       categorizeLinuxDep(dep.Name),
		})
	}

	return nil
}

func categorizeLinuxDep(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "gtk"):
		return "gtk"
	case strings.Contains(lower, "webkit"):
		return "webkit"
	case name == "gcc" || name == "pkg-config":
		return "build-tools"
	case name == "npm":
		return "frontend"
	case name == "docker":
		return "optional"
	default:
		return "other"
	}
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkLinuxSpecific()
}

func (d *Doctor) checkGoInstallation() {
	if d.report.Build.GoVersion == "" {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Go Installation",
			Message:  "Go installation not found",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/",
			Fix: &Fix{
				Description: "Install Go from https://go.dev/dl/",
				ManualSteps: []string{
					"Download Go from https://go.dev/dl/",
					"Extract and add to PATH",
				},
			},
		})
	}
}

func (d *Doctor) checkLinuxSpecific() {
	missingRequired := d.report.Dependencies.RequiredMissing()
	if len(missingRequired) > 0 {
		var commands []string
		for _, dep := range missingRequired {
			if dep.InstallCommand != "" {
				commands = append(commands, dep.InstallCommand)
			}
		}

		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Missing Dependencies",
			Message:  "Required system packages are not installed",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/#linux",
			Fix: &Fix{
				Description:  "Install missing packages",
				Command:      strings.Join(commands, " && "),
				RequiresSudo: true,
			},
		})
	}
}
