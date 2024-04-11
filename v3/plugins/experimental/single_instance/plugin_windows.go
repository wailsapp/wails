//go:build windows

package single_instance

import (
	"github.com/wailsapp/wails/v3/pkg/w32"
	"syscall"
)

type enumWindowsCallback func(hwnd syscall.Handle, lParam uintptr) uintptr

func enumWindowsProc(hwnd syscall.Handle, lParam uintptr) uintptr {
	_, processID := w32.GetWindowThreadProcessId(uintptr(hwnd))
	targetProcessID := uint32(lParam)
	if uint32(processID) == targetProcessID {
		// Bring the window forward
		w32.SetForegroundWindow(w32.HWND(hwnd))
	}

	// Continue enumeration
	return 1
}

func (p *Plugin) activeInstance(pid int) error {

	// Get the window associated with the process ID.
	targetProcessID := uint32(pid) // Replace with the desired process ID

	w32.EnumWindows(
		syscall.NewCallback(enumWindowsCallback(enumWindowsProc)),
		uintptr(targetProcessID),
	)
	return nil
}
