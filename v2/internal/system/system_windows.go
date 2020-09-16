// +build windows

package system

import (
	"fmt"
	"syscall"
)

func (i *Info) discover() {
	dll := syscall.MustLoadDLL("kernel32.dll")
	p := dll.MustFindProc("GetVersion")
	v, _, _ := p.Call()
	fmt.Printf("Windows version %d.%d (Build %d)\n", byte(v), uint8(v>>8), uint16(v>>16))
}
