//go:build !darwin || ios || server

package application

import (
	"fmt"
	"os/exec"
)

// platformOpenWith starts app (an executable on PATH or an executable path)
// with path as its argument. There is no portable "open with" beyond that.
func platformOpenWith(path string, app string) error {
	executable, err := exec.LookPath(app)
	if err != nil {
		return fmt.Errorf("%w: %q is not an executable", ErrOpenWithUnsupported, app)
	}
	cmd := exec.Command(executable, path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("browser: OpenWith %q: %w", app, err)
	}
	go cmd.Wait() //nolint:errcheck
	return nil
}

func platformApplicationsForFile(path string) []AppInfo {
	return []AppInfo{}
}

func platformActivateApplication(bundleID string) error {
	return fmt.Errorf("%w: %s", ErrApplicationNotRunning, bundleID)
}
