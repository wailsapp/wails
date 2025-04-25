//go:build windows

package w32

import (
	"syscall"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procSetWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")

const (
	WDA_NONE               = 0x00000000
	WDA_MONITOR            = 0x00000001
	WDA_EXCLUDEFROMCAPTURE = 0x00000011 // windows 10 2004+
)

func SetWindowDisplayAffinity(hwnd uintptr, affinity uint32) bool {
	ret, _, _ := procSetWindowDisplayAffinity.Call(
		hwnd,
		uintptr(affinity),
	)
	return ret != 0
}
