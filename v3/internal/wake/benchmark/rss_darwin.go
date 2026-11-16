//go:build darwin

package benchmark

import (
	"os"
	"syscall"
)

func peakRSSBytes(state *os.ProcessState) int64 {
	usage := state.SysUsage().(*syscall.Rusage)
	return usage.Maxrss
}
