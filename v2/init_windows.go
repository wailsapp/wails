package wails

import (
	"fmt"
	"syscall"
)

// Init is called at the start of the application
func Init() error {
	status, r, err := syscall.NewLazyDLL("user32.dll").NewProc("SetProcessDPIAware").Call()
	if status == 0 {
		return fmt.Errorf("exit status %d: %v %v", status, r, err)
	}
	return nil
}
