//go:build !linux && !darwin

package benchmark

import "os"

func peakRSSBytes(_ *os.ProcessState) int64 {
	return 0
}
