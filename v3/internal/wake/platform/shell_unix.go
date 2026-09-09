//go:build !windows

package platform

import "os/exec"

// ShellCommand returns an *exec.Cmd that evaluates script through the
// platform's default shell. On Unix-like systems we use `sh -c …`,
// matching the upstream Taskfile reference implementation.
func ShellCommand(script string) *exec.Cmd {
	return exec.Command("sh", "-c", script)
}
