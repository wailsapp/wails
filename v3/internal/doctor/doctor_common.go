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
}
