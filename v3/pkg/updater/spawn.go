package updater

import (
	"os"
	"os/exec"
)

// selfExecutable returns the path of the running executable. It exists as a
// thin wrapper so tests can override it for spawn assertions without poking
// at os.Executable directly.
func selfExecutable() (string, error) {
	return os.Executable()
}

// newDetachedCommand builds an exec.Cmd for the helper invocation. Stdio is
// disconnected from the parent so the helper survives the parent's exit on
// every platform.
func newDetachedCommand(path string) *exec.Cmd {
	cmd := exec.Command(path)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	// Platform-specific session/process-group detachment is handled in
	// spawn_unix.go and spawn_windows.go.
	applyDetachAttrs(cmd)
	return cmd
}
