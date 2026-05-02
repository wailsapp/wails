//go:build windows

package doctor

import (
	"os/exec"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/go-webview2/webviewloader"
)

func getInfo() (map[string]string, bool) {
	ok := true
	result := make(map[string]string)
	result["Go WebView2Loader"] = lo.Ternary(webviewloader.UsingGoWebview2Loader, "true", "false")
	webviewVersion, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		ok = false
		webviewVersion = "Error:" + err.Error()
	}
	result["WebView2 Version"] = webviewVersion
	return result, ok
}

func getNSISVersion() string {
	// Execute nsis
	output, err := exec.Command("makensis", "-VERSION").Output()
	if err != nil {
		return "Not Installed"
	}
	return string(output)
}

func getMakeAppxVersion() string {
	// Check if MakeAppx.exe is available (part of Windows SDK)
	_, err := exec.LookPath("MakeAppx.exe")
	if err != nil {
		return "Not Installed"
	}
	return "Installed"
}

func getMSIXPackagingToolVersion() string {
	// Check if MSIX Packaging Tool is installed
	// Use PowerShell to check if the app is installed from Microsoft Store
	cmd := exec.Command("powershell", "-Command", "Get-AppxPackage -Name Microsoft.MSIXPackagingTool")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 || !strings.Contains(string(output), "Microsoft.MSIXPackagingTool") {
		return "Not Installed"
	}

	if strings.Contains(string(output), "Version") {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Version") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Installed (Version Unknown)"
}

func getSignToolVersion() string {
	// Check if signtool.exe is available (part of Windows SDK)
	_, err := exec.LookPath("signtool.exe")
	if err != nil {
		return "Not Installed"
	}
	return "Installed"
}

func checkPlatformDependencies(result map[string]string, ok *bool) {
	checkCommonDependencies(result, ok)
	// add nsis
	result["NSIS"] = getNSISVersion()

	// Add MSIX tooling checks
	result["MakeAppx.exe (Windows SDK)"] = getMakeAppxVersion()
	result["MSIX Packaging Tool"] = getMSIXPackagingToolVersion()
	result["SignTool.exe (Windows SDK)"] = getSignToolVersion()
}
