//go:build linux

package benchmark

import (
	"os"
	"syscall"
)

func peakRSSBytes(state *os.ProcessState) int64 {
	usage := state.SysUsage().(*syscall.Rusage)
	return usage.Maxrss * 1024
}
