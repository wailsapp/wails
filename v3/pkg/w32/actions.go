//go:build windows

package w32

func Undo(hwnd HWND) {
	SendMessage(hwnd, WM_UNDO, 0, 0)
}

func Cut(hwnd HWND) {
	SendMessage(hwnd, WM_CUT, 0, 0)
}

func Copy(hwnd HWND) {
	SendMessage(hwnd, WM_COPY, 0, 0)
}

func Paste(hwnd HWND) {
	SendMessage(hwnd, WM_PASTE, 0, 0)
}

func Delete(hwnd HWND) {
	SendMessage(hwnd, WM_CLEAR, 0, 0)
}

func SelectAll(hwnd HWND) {
	SendMessage(hwnd, WM_SELECTALL, 0, 0)
}
