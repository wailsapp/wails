package win32

func CreatePopupMenu() HMENU {
	ret, _, _ := procCreatePopupMenu.Call(0, 0, 0, 0)
	return HMENU(ret)
}

func TrackPopupMenu(menu HMENU, flags uint, x, y int, wnd HWND) bool {
	ret, _, _ := procTrackPopupMenu.Call(
		uintptr(menu),
		uintptr(flags),
		uintptr(x),
		uintptr(y),
		0,
		uintptr(wnd),
		0,
	)
	return ret != 0
}

func AppendMenu(menu HMENU, flags uintptr, id uintptr, text string) bool {
	ret, _, _ := procAppendMenuW.Call(
		uintptr(menu),
		flags,
		id,
		MustStringToUTF16uintptr(text),
	)
	return ret != 0
}
