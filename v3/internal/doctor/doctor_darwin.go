//go:build darwin

package doctor

import (
	"bytes"
	"github.com/samber/lo"
	"os/exec"
	"strings"
	"syscall"
)

func getSysctl(name string) string {
	value, err := syscall.Sysctl(name)
	if err != nil {
		return "unknown"
	}
	return value
}

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)
	ok := true

	// Determine if the app is running on Apple Silicon
	// Credit: https://www.yellowduck.be/posts/detecting-apple-silicon-via-go/
	appleSilicon := "unknown"
	r, err := syscall.Sysctl("sysctl.proc_translated")
	if err == nil {
		appleSilicon = lo.Ternary(r == "\x00\x00\x00" || r == "\x01\x00\x00", "true", "false")
	}
	result["Apple Silicon"] = appleSilicon
	result["CPU"] = getSysctl("machdep.cpu.brand_string")

	return result, ok
}

func checkPlatformDependencies(result map[string]string, ok *bool) {

	// Check for xcode command line tools
	output, err := exec.Command("xcode-select", "-v").Output()
	cliToolsVersion := "N/A. Install by running: `xcode-select --install`"
	if err != nil {
		*ok = false
	} else {
		cliToolsVersion = strings.TrimPrefix(string(output), "xcode-select version ")
		cliToolsVersion = strings.TrimSpace(cliToolsVersion)
		cliToolsVersion = strings.TrimSuffix(cliToolsVersion, ".")
	}
	result["Xcode cli tools"] = cliToolsVersion

	checkCommonDependencies(result, ok)

	// Check for nsis
	nsisVersion := []byte("Not Installed. Install with `brew install makensis`.")
	output, err = exec.Command("makensis", "-VERSION").Output()
	if err == nil && output != nil {
		nsisVersion = output
	}
	nsisVersion = bytes.TrimSpace(nsisVersion)

	result["*NSIS"] = string(nsisVersion)
}
