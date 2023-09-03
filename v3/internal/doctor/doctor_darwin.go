//go:build darwin

package doctor

import (
	"github.com/samber/lo"
	"os/exec"
	"strings"
	"syscall"
)

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

	// Check for xcode command line tools
	output, err := exec.Command("xcode-select", "-v").Output()
	cliToolsVersion := "N/A. Install by running: `xcode-select --install`"
	if err != nil {
		ok = false
	} else {
		cliToolsVersion = strings.TrimPrefix(string(output), "xcode-select version ")
		cliToolsVersion = strings.TrimSpace(cliToolsVersion)
		cliToolsVersion = strings.TrimSuffix(cliToolsVersion, ".")
	}
	result["Xcode cli tools"] = cliToolsVersion

	return result, ok
}
