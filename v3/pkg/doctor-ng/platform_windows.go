//go:build windows

package doctorng

import (
	"os/exec"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/go-webview2/webviewloader"
)

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	extras["Go WebView2Loader"] = lo.Ternary(webviewloader.UsingGoWebview2Loader, "true", "false")

	webviewVersion, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		extras["WebView2 Version"] = "Error: " + err.Error()
	} else {
		extras["WebView2 Version"] = webviewVersion
	}

	return extras
}

func (d *Doctor) collectDependencies() error {
	d.checkCommonDependencies()

	nsisVersion, nsisStatus := d.checkCommand("makensis", "-VERSION")
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:           "NSIS",
		Version:        nsisVersion,
		Status:         nsisStatus,
		Required:       false,
		InstallCommand: "choco install nsis",
		Category:       "optional",
		Description:    "For Windows installer generation",
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
		Name:           "npm",
		Version:        npmVersion,
		Status:         npmStatus,
		Required:       true,
		InstallCommand: "choco install nodejs",
		Category:       "frontend",
	})

	dockerVersion, dockerStatus := d.checkCommand("docker", "--version")
	if dockerStatus == StatusOK {
		dockerVersion = strings.TrimPrefix(dockerVersion, "Docker version ")
		if idx := strings.Index(dockerVersion, ","); idx > 0 {
			dockerVersion = dockerVersion[:idx]
		}
	}
	d.report.Dependencies = append(d.report.Dependencies, &Dependency{
		Name:           "docker",
		Version:        dockerVersion,
		Status:         dockerStatus,
		Required:       false,
		InstallCommand: "choco install docker-desktop",
		Category:       "optional",
		Description:    "For cross-compilation",
	})
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkWebView2()
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
				Command:     "choco install golang",
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
