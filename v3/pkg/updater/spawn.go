package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// selfExecutable returns the path of the running executable. Held in a
// package-level var so tests can override without poking at os.Executable.
var selfExecutable = func() (string, error) {
	return os.Executable()
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

const helperReadyTimeout = 5 * time.Second

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
