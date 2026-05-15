//go:build windows

package updater

import "os"

// syscallSignalZero returns the probe signal used by isAlive. On Windows the
// signal-0 trick is not portable; instead os.FindProcess returns a valid
// handle only for processes the caller can open. We return nil so the
// proc.Signal call performs the cheap permission check that doubles as a
// liveness test on Windows.
func syscallSignalZero() os.Signal { return nil }
