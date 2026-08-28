//go:build windows

package commands

import "os"

const statusControlCExit = 0xC000013A

func isInterruptProcessState(state *os.ProcessState) bool {
	return uint32(state.ExitCode()) == statusControlCExit
}
