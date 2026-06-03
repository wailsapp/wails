//go:build !windows

package parse

import "os/exec"

// shellCommand returns an *exec.Cmd that will evaluate script through the
// platform's default POSIX shell. On Unix-like systems we use `sh -c …`,
// matching the Taskfile reference implementation and what wake's executor
// uses elsewhere.
func shellCommand(script string) *exec.Cmd {
	return exec.Command("sh", "-c", script)
}
