//go:build windows

package w32

var procGetParent = moduser32.NewProc("GetParent")

// GetParent returns the parent window of hwnd, or zero for a top-level window.
func GetParent(hwnd HWND) HWND {
	parent, _, _ := procGetParent.Call(hwnd)
	return parent
}
