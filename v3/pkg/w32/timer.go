//go:build windows

package w32

func SetTimer(hwnd HWND, nIDEvent uintptr, uElapse uint32, lpTimerFunc uintptr) uintptr {
	ret, _, _ := procSetTimer.Call(
		uintptr(hwnd),
		nIDEvent,
		uintptr(uElapse),
		lpTimerFunc)
	return ret
}

func KillTimer(hwnd HWND, nIDEvent uintptr) bool {
	ret, _, _ := procKillTimer.Call(
		uintptr(hwnd),
		nIDEvent)
	return ret != 0
}
