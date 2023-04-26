package w32

// MenuItemFromPoint determines which menu item, if any, is at the specified
// location.
func MenuItemFromPoint(w HWND, m HMENU, screen POINT) int {
	ret, _, _ := menuItemFromPoint.Call(
		uintptr(w),
		uintptr(m),
		uintptr(screen.X),
		uintptr(screen.Y),
	)
	return int(ret)
}

func GetClassLongPtr(w HWND, index int) uintptr {
	return uintptr(GetClassLong(w, index))
}

func SetClassLongPtr(w HWND, index int, value uintptr) uintptr {
	return uintptr(SetClassLong(w, index, int32(value)))
}

func GetWindowLongPtr(hwnd HWND, index int) uintptr {
	return uintptr(GetWindowLong(hwnd, index))
}

func SetWindowLongPtr(hwnd HWND, index int, value uintptr) uintptr {
	return uintptr(SetWindowLong(hwnd, index, int32(value)))
}
