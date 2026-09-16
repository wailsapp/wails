//go:build windows

package platform

import "os/exec"

// ShellCommand returns an *exec.Cmd that evaluates script through the
// platform's default shell. On Windows we use `cmd /C` because `sh` isn't
// available on a stock Windows host — without this every wake-routed
// `vars: {X: {sh: ...}}`, status:, and precondition: shell-out would
// fail.
//
// CLI Taskfile commands use the embedded Task runtime and its shell semantics.
func ShellCommand(script string) *exec.Cmd {
	return exec.Command("cmd", "/C", script)
}
