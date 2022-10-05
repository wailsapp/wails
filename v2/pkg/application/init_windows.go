//go:build windows

package application

import (
	"fmt"
	"syscall"
)

func applicationInit() error {
	status, r, err := syscall.NewLazyDLL("user32.dll").NewProc("SetProcessDPIAware").Call()
	if status == 0 {
		return fmt.Errorf("exit status %d: %v %v", status, r, err)
	}
	return nil
}
