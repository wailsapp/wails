//go:build windows

package doctor

import (
	"github.com/samber/lo"
	"github.com/wailsapp/go-webview2/webviewloader"
	"os/exec"
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

func checkPlatformDependencies(result map[string]string, ok *bool) {
	checkCommonDependencies(result, ok)
	// add nsis
	result["NSIS"] = getNSISVersion()
}
