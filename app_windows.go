// +build windows !linux !darwin

package wails

import (
	"fmt"
	"log"
	"syscall"
)

func platformInit() {
	err := SetProcessDPIAware()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// SetProcessDPIAware via user32.dll
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setprocessdpiaware
// Also, thanks Jack Mordaunt! https://github.com/wailsapp/wails/issues/293
func SetProcessDPIAware() error {
	status, r, err := syscall.NewLazyDLL("user32.dll").NewProc("SetProcessDPIAware").Call()
	if status == 0 {
		return fmt.Errorf("exit status %d: %v %v", status, r, err)
	}
	return nil
}
