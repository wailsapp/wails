package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultMacOSDeploymentTarget = "10.13"
	defaultSafariStagedFramework = "/Library/Apple/System/Library/StagedFrameworks/Safari"
)

var safariStagedFrameworks = []string{
	"WebKit",
	"JavaScriptCore",
	"WebCore",
	"WebKitLegacy",
}

type macOSVersion struct {
	major int
	minor int
}

type macOSBuildConfig struct {
	deploymentTarget  string
	stagedFrameworks  string
	externalLinkerArg string
}

// newMacOSBuildConfig turns the standard deployment-target environment variable
// into the additional flags needed by Apple's Safari-updated WebKit on Monterey.
// The staged WebKit shipped with Safari 17.6 declares macOS 12.3 as its own
// minimum. Older targets therefore retain Wails' normal system-WebKit path;
// macOS 13 and newer already expose their current WebKit publicly.
func newMacOSBuildConfig(deploymentTarget, stagedFrameworks string) (macOSBuildConfig, error) {
	if deploymentTarget == "" {
		deploymentTarget = defaultMacOSDeploymentTarget
	}

	version, err := parseMacOSVersion(deploymentTarget)
	if err != nil {
		return macOSBuildConfig{}, err
	}

	config := macOSBuildConfig{deploymentTarget: deploymentTarget}
	if version.major != 12 || version.minor < 3 {
		return config, nil
	}

	if stagedFrameworks == "" {
		stagedFrameworks = defaultSafariStagedFramework
	}
	if err := validateSafariStagedFrameworks(stagedFrameworks); err != nil {
		return macOSBuildConfig{}, err
	}

	config.stagedFrameworks = stagedFrameworks
	config.externalLinkerArg = safariExternalLinkerArg(stagedFrameworks)
	return config, nil
}

func parseMacOSVersion(value string) (macOSVersion, error) {
	parts := strings.Split(value, ".")
	if len(parts) > 2 || value == "" {
		return macOSVersion{}, fmt.Errorf("MACOSX_DEPLOYMENT_TARGET must be a numeric macOS version: %s", value)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return macOSVersion{}, fmt.Errorf("MACOSX_DEPLOYMENT_TARGET must be a numeric macOS version: %s", value)
	}

	minor := 0
	if len(parts) == 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil || minor < 0 {
			return macOSVersion{}, fmt.Errorf("MACOSX_DEPLOYMENT_TARGET must be a numeric macOS version: %s", value)
		}
	}

	return macOSVersion{major: major, minor: minor}, nil
}

func validateSafariStagedFrameworks(root string) error {
	for _, framework := range safariStagedFrameworks {
		binary := filepath.Join(root, framework+".framework", "Versions", "A", framework)
		if _, err := os.Stat(binary); err != nil {
			return fmt.Errorf("Safari staged WebKit is incomplete: missing %s", filepath.Join(root, framework+".framework"))
		}
	}
	return nil
}

func safariExternalLinkerArg(stagedFrameworks string) string {
	flags := []string{
		"-F" + stagedFrameworks,
		"-framework UniformTypeIdentifiers",
		"-Wl,-dyld_env,DYLD_VERSIONED_FRAMEWORK_PATH=" + stagedFrameworks,
		"-Wl,-dyld_env,DYLD_VERSIONED_LIBRARY_PATH=" + stagedFrameworks,
	}
	return "-extldflags " + strconv.Quote(strings.Join(flags, " "))
}
