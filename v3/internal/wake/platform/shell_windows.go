//go:build windows

package platform

import "os/exec"

// ShellCommand returns an *exec.Cmd that evaluates script through the
// platform's default shell. On Windows we use `cmd /C` because `sh` isn't
// available on a stock Windows host — without this every wake-routed
// `vars: {X: {sh: ...}}`, status:, and precondition: shell-out would
// fail.
//
// Users who want a POSIX shell (Git Bash, MSYS, WSL) can keep
// WAILS_USE_WAKE unset and use the embedded Task runtime instead.
func ShellCommand(script string) *exec.Cmd {
	return exec.Command("cmd", "/C", script)
}
