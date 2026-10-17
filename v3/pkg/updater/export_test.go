package updater

import (
	"os/exec"
	"time"
)

// SetSelfExecutableForTest replaces the package-level selfExecutable
// resolver for the duration of a test. Returns a restore function the
// caller should defer.
func SetSelfExecutableForTest(f func() (string, error)) (restore func()) {
	prev := selfExecutable
	selfExecutable = f
	return func() { selfExecutable = prev }
}

// AppliedMarkerForTest exposes the path where the helper records the
// replaced version for target.
var AppliedMarkerForTest = appliedMarker

// SetNewDetachedCommandForTest replaces the package-level command builder
// for the duration of a test. Used to substitute a benign command (e.g.
// /usr/bin/true) so Restart's spawn succeeds without re-execing the test
// binary into helper mode.
func SetNewDetachedCommandForTest(f func(path string) *exec.Cmd) (restore func()) {
	prev := newDetachedCommand
	newDetachedCommand = f
	return func() { newDetachedCommand = prev }
}

// SetWaitForHelperReadyForTest replaces the helper readiness wait for the
// duration of a test. Returns a restore function the caller should defer.
func SetWaitForHelperReadyForTest(f func(path string, timeout time.Duration) error) (restore func()) {
	prev := waitForHelperReady
	waitForHelperReady = f
	return func() { waitForHelperReady = prev }
}
