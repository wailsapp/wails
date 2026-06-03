//go:build windows

package parse

import "os/exec"

// shellCommand returns an *exec.Cmd that will evaluate script through the
// platform's default shell. On Windows we use `cmd /C` because `sh` isn't
// installed on a stock Windows machine — without this, every Taskfile with
// `vars: {X: {sh: ...}}` would break under wake on Windows (the original
// implementation hard-coded `sh -c`, which is fine on Unix but not here).
//
// Users with Git Bash / WSL / MSYS who explicitly want a POSIX shell can
// keep using `task` instead — wake's fallback path picks the embedded
// Task runtime for any feature wake doesn't fully support.
func shellCommand(script string) *exec.Cmd {
	return exec.Command("cmd", "/C", script)
}
