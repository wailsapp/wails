//go:build linux

package doctor

import (
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/usernssandbox"
)

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)

	// Check session type (X11/Wayland)
	result["XDG_SESSION_TYPE"] = getSessionType()

	// Check desktop environment
	result["Desktop Environment"] = getDesktopEnvironment()

	// Check for NVIDIA driver
	result["NVIDIA Driver"] = getNvidiaDriverInfo()

	return result, true
}

func getSessionType() string {
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType == "" {
		return "unset"
	}
	return sessionType
}

func getDesktopEnvironment() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de == "" {
		return "unset"
	}
	return de
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

func checkPlatformDependencies(result map[string]string, ok *bool) {
	info, _ := operatingsystem.Info()

	pm := packagemanager.Find(info.ID)
	deps, _ := packagemanager.Dependencies(pm)
	for _, dep := range deps {
		var status string

		switch true {
		case !dep.Installed:
			if dep.Optional {
				status = "[Optional] "
			} else {
				*ok = false
			}
			status += "not installed."
			if dep.InstallCommand != "" {
				status += " Install with: " + dep.InstallCommand
			}
		case dep.Version != "":
			status = dep.Version
		}

		result[dep.Name] = status
	}

	checkWebKitSandbox(result, ok)

	checkCommonDependencies(result, ok)
}

// checkWebKitSandbox verifies bwrap can create a user namespace, which
// WebKitGTK's sandbox needs (https://github.com/wailsapp/wails/issues/6172).
func checkWebKitSandbox(result map[string]string, ok *bool) {
	const key = "WebKit sandbox (bwrap)"
	res := usernssandbox.Check(usernssandbox.Env{}, true)
	switch {
	case res.SandboxDisabled:
		result[key] = "[Warning] " + usernssandbox.DisableEnvVar + " is set: WebKit runs without its sandbox"
	case res.Status == usernssandbox.Blocked:
		*ok = false
		result[key] = "user namespaces blocked: " + res.Message() + "\n" + usernssandbox.Advice()
	case res.Status == usernssandbox.Unknown:
		result[key] = "[Warning] " + res.Message()
	}
}
