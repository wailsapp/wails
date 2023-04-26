package w32

// MenuItemFromPoint determines which menu item, if any, is at the specified
// location.
func MenuItemFromPoint(w HWND, m HMENU, screen POINT) int {
	ret, _, _ := menuItemFromPoint.Call(
		uintptr(w),
		uintptr(m),
		uintptr(uint64(screen.X)<<32|uint64(screen.Y)),
	)
	return int(ret)
}

func GetClassLongPtr(w HWND, index int) uintptr {
	ret, _, _ := getClassLongPtr.Call(uintptr(w), uintptr(index))
	return ret
}

func SetClassLongPtr(w HWND, index int, value uintptr) uintptr {
	ret, _, _ := setClassLongPtr.Call(
		uintptr(w),
		uintptr(index),
		value,
	)
	return ret
}

func GetWindowLongPtr(hwnd HWND, index int) uintptr {
	ret, _, _ := getWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(index),
	)
	return ret
}

func SetWindowLongPtr(hwnd HWND, index int, value uintptr) uintptr {
	ret, _, _ := setWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(index),
		value,
	)
	return ret
}
