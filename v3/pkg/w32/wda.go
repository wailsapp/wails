//go:build windows

package w32

const (
	WDA_NONE               = 0x00000000
	WDA_MONITOR            = 0x00000001
	WDA_EXCLUDEFROMCAPTURE = 0x00000011 // windows 10 2004+
)

func SetWindowDisplayAffinity(hwnd uintptr, affinity uint32) bool {
	if affinity == WDA_EXCLUDEFROMCAPTURE && !IsWindowsVersionAtLeast(10, 0, 19041) {
		// for older windows versions, use WDA_MONITOR
		affinity = WDA_MONITOR
	}
	ret, _, _ := procSetWindowDisplayAffinity.Call(
		hwnd,
		uintptr(affinity),
	)
	return ret != 0
}
