package updater

import "os/exec"

// SetSelfExecutableForTest replaces the package-level selfExecutable
// resolver for the duration of a test. Returns a restore function the
// caller should defer.
func SetSelfExecutableForTest(f func() (string, error)) (restore func()) {
	prev := selfExecutable
	selfExecutable = f
	return func() { selfExecutable = prev }
}

// SetNewDetachedCommandForTest replaces the package-level command builder
// for the duration of a test. Used to substitute a benign command (e.g.
// /usr/bin/true) so Restart's spawn succeeds without re-execing the test
// binary into helper mode.
func SetNewDetachedCommandForTest(f func(path string) *exec.Cmd) (restore func()) {
	prev := newDetachedCommand
	newDetachedCommand = f
	return func() { newDetachedCommand = prev }
}
