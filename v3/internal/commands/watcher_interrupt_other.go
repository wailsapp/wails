//go:build !darwin && !linux && !windows

package commands

import "os"

func isInterruptProcessState(*os.ProcessState) bool {
	return false
}
