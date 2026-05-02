//go:build windows

package doctorng

import (
	"os/exec"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/go-webview2/webviewloader"
)

type windowsPackageManager int

const (
	pmNone windowsPackageManager = iota
	pmWinget
	pmChoco
	pmScoop
)

var detectedPM windowsPackageManager
var pmDetected bool

func detectWindowsPackageManager() windowsPackageManager {
	if pmDetected {
		return detectedPM
	}
	pmDetected = true

	if _, err := exec.LookPath("winget"); err == nil {
		detectedPM = pmWinget
		return detectedPM
	}
	if _, err := exec.LookPath("scoop"); err == nil {
		detectedPM = pmScoop
		return detectedPM
	}
	if _, err := exec.LookPath("choco"); err == nil {
		detectedPM = pmChoco
		return detectedPM
	}

	detectedPM = pmNone
	return detectedPM
}

func windowsInstallCmd(winget, scoop, choco, manual string) string {
	switch detectWindowsPackageManager() {
	case pmWinget:
		return winget
	case pmScoop:
		return scoop
	case pmChoco:
		return choco + " (requires admin)"
	default:
		return manual
	}
}

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	extras["Go WebView2Loader"] = lo.Ternary(webviewloader.UsingGoWebview2Loader, "true", "false")

	webviewVersion, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		extras["WebView2 Version"] = "Error: " + err.Error()
	} else {
		extras["WebView2 Version"] = webviewVersion
	}

	pm := "none"
	switch detectWindowsPackageManager() {
	case pmWinget:
		pm = "winget"
	case pmScoop:
		pm = "scoop"
	case pmChoco:
		pm = "choco"
	}
	extras["Package Manager"] = pm

	return extras
}

func (d *Doctor) collectDependencies() error {
	d.checkCommonDependencies()

	nsisVersion, nsisStatus := d.checkCommand("makensis", "-VERSION")
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "NSIS",
		Version: nsisVersion,
		Status:  nsisStatus,
		InstallCommand: windowsInstallCmd(
			"winget install NSIS.NSIS",
			"scoop install nsis",
			"choco install nsis",
			"Download from https://nsis.sourceforge.io/",
		),
		Required:    false,
		Category:    "optional",
		Description: "For Windows installer generation",
	})

	makeAppxStatus := StatusMissing
	if _, err := exec.LookPath("MakeAppx.exe"); err == nil {
		makeAppxStatus = StatusOK
	}
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:        "MakeAppx.exe",
		Status:      makeAppxStatus,
		Required:    false,
		Category:    "optional",
		Description: "Part of Windows SDK, for MSIX packaging",
	})

	signToolStatus := StatusMissing
	if _, err := exec.LookPath("signtool.exe"); err == nil {
		signToolStatus = StatusOK
	}
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:        "SignTool.exe",
		Status:      signToolStatus,
		Required:    false,
		Category:    "optional",
		Description: "Part of Windows SDK, for code signing",
	})

	return nil
}

func (d *Doctor) checkCommand(cmd string, args ...string) (string, Status) {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", StatusMissing
	}
	return strings.TrimSpace(string(output)), StatusOK
}

func (d *Doctor) checkCommonDependencies() {
	npmVersion, npmStatus := d.checkCommand("npm", "--version")
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "npm",
		Version: npmVersion,
		Status:  npmStatus,
		InstallCommand: windowsInstallCmd(
			"winget install OpenJS.NodeJS.LTS",
			"scoop install nodejs-lts",
			"choco install nodejs-lts",
			"Download from https://nodejs.org/",
		),
		Required: true,
		Category: "frontend",
	})

	dockerVersion, dockerStatus := d.checkCommand("docker", "--version")
	if dockerStatus == StatusOK {
		dockerVersion = strings.TrimPrefix(dockerVersion, "Docker version ")
		if idx := strings.Index(dockerVersion, ","); idx > 0 {
			dockerVersion = dockerVersion[:idx]
		}
	}
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:    "docker",
		Version: dockerVersion,
		Status:  dockerStatus,
		InstallCommand: windowsInstallCmd(
			"winget install Docker.DockerDesktop",
			"Download from https://docker.com/",
			"choco install docker-desktop",
			"Download from https://docker.com/",
		),
		Required:    false,
		Category:    "optional",
		Description: "For cross-compilation",
	})
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkWebView2()
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
				Command: windowsInstallCmd(
					"winget install GoLang.Go",
					"scoop install go",
					"choco install golang",
					"Download from https://go.dev/dl/",
				),
			},
		})
	}
}

func (d *Doctor) checkWebView2() {
	_, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "WebView2 Runtime",
			Message:  "WebView2 runtime not found or unavailable",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/#windows",
			Fix: &Fix{
				Description: "Install Microsoft Edge WebView2 Runtime",
				ManualSteps: []string{
					"Download from https://developer.microsoft.com/en-us/microsoft-edge/webview2/",
					"Or it may be bundled with recent Windows updates",
				},
			},
		})
	}
}

func (d *Doctor) checkPackageManager() {
	if detectWindowsPackageManager() == pmNone {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Package Manager",
			Message:  "No package manager found (winget, scoop, or choco)",
			Severity: SeverityWarning,
			HelpURL:  "/getting-started/installation/#windows",
			Fix: &Fix{
				Description: "Install a package manager for easier dependency management",
				ManualSteps: []string{
					"winget: Built into Windows 11, or install from Microsoft Store",
					"scoop: Run in PowerShell: irm get.scoop.sh | iex",
					"choco: See https://chocolatey.org/install",
				},
			},
		})
	}
}
