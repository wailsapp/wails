package doctor

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func checkCommonDependencies(result map[string]string, ok *bool) {
	// Check for npm
	npmVersion := []byte("Not Installed. Requires npm >= 7.0.0")
	npmVersion, err := exec.Command("npm", "-v").Output()
	if err != nil {
		*ok = false
	} else {
		npmVersion = bytes.TrimSpace(npmVersion)
		// Check that it's at least version 7 by converting first byte to int and checking if it's >= 7
		// Parse the semver string
		semver := strings.Split(string(npmVersion), ".")
		if len(semver) > 0 {
			major, _ := strconv.Atoi(semver[0])
			if major < 7 {
				*ok = false
				npmVersion = append(npmVersion, []byte(". Installed, but requires npm >= 7.0.0")...)
			} else {
				*ok = true
			}
		}
	}
	result["npm"] = string(npmVersion)

	// Check for Docker (optional - used for macOS cross-compilation from Linux)
	checkDocker(result)
}

func checkDocker(result map[string]string) {
	dockerVersion, err := exec.Command("docker", "--version").Output()
	if err != nil {
		result["docker"] = "*Not installed (optional - for cross-compilation)"
		return
	}

	// Check if Docker daemon is running
	_, err = exec.Command("docker", "info").Output()
	if err != nil {
		version := strings.TrimSpace(string(dockerVersion))
		result["docker"] = "*" + version + " (daemon not running)"
		return
	}

	version := strings.TrimSpace(string(dockerVersion))

	// Check for the unified cross-compilation image
	imageCheck, _ := exec.Command("docker", "image", "inspect", "wails-cross").Output()
	if len(imageCheck) == 0 {
		result["docker"] = "*" + version + " (wails-cross image not built - run: wails3 task setup:docker)"
	} else {
		result["docker"] = "*" + version + " (cross-compilation ready)"
	}
}
