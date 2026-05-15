//go:build !windows

package updater

import (
	"os"
	"syscall"
)

// syscallSignalZero returns the "is alive" probe signal used by isAlive on
// Unix-like systems. Sending signal 0 to a pid either succeeds (process is
// running and we have permission) or fails (process is gone / no permission).
func syscallSignalZero() os.Signal { return syscall.Signal(0) }
