//go:build darwin && !ios

package mac

import "path/filepath"

func resourcePathForExecutable(executable string) (string, bool) {
	executable = filepath.Clean(executable)
	macOSPath := filepath.Dir(executable)
	contentsPath := filepath.Dir(macOSPath)
	bundlePath := filepath.Dir(contentsPath)

	if filepath.Base(macOSPath) != "MacOS" ||
		filepath.Base(contentsPath) != "Contents" ||
		filepath.Ext(bundlePath) != ".app" {
		return "", false
	}

	return filepath.Join(contentsPath, "Resources"), true
}
