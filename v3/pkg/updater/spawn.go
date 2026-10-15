package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// selfExecutable returns the path of the running executable, or of the
// AppImage it was launched from. Held in a package-level var so tests can
// override without poking at os.Executable.
var selfExecutable = func() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return appImageOr(exe), nil
}

// appImageOr returns $APPIMAGE when exe runs from that AppImage's mount.
// The mount is a read-only squashfs that vanishes on exit, so the AppImage
// file is what must be replaced and re-executed. APPIMAGE is inherited by
// child processes, hence the check that exe really lives under $APPDIR.
func appImageOr(exe string) string {
	img, mount := os.Getenv("APPIMAGE"), os.Getenv("APPDIR")
	if img == "" || mount == "" || !strings.HasPrefix(exe, filepath.Clean(mount)+string(filepath.Separator)) {
		return exe
	}
	if info, err := os.Stat(img); err != nil || !info.Mode().IsRegular() {
		return exe
	}
	return img
}

// resolveTarget returns the path the helper replaces: the running
// executable, or its enclosing .app bundle on macOS.
func resolveTarget() (string, error) {
	self, err := selfExecutable()
	if err != nil {
		return "", err
	}
	return bundleTarget(self), nil
}

// newDetachedCommand builds an exec.Cmd for the helper invocation. Stdio is
// disconnected from the parent so the helper survives the parent's exit on
// every platform. Held in a package-level var so tests can substitute a
// command that's safe to actually Start() (e.g. /usr/bin/true) without
// re-execing the test binary.
var newDetachedCommand = func(path string) *exec.Cmd {
	cmd := exec.Command(path)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	// Platform-specific session/process-group detachment is handled in
	// spawn_unix.go and spawn_windows.go.
	applyDetachAttrs(cmd)
	return cmd
}

const helperReadyTimeout = 30 * time.Second

var waitForHelperReady = func(readyPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		payload, err := os.ReadFile(readyPath)
		if err == nil {
			status := strings.TrimSpace(string(payload))
			switch {
			case status == "ready":
				return nil
			case strings.HasPrefix(status, "error:"):
				return fmt.Errorf("%w: %s", ErrHelperNotReady, strings.TrimSpace(strings.TrimPrefix(status, "error:")))
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("%w: read readiness signal: %v", ErrHelperNotReady, err)
		}
		if time.Now().After(deadline) {
			return ErrHelperNotReady
		}
		time.Sleep(25 * time.Millisecond)
	}
}
