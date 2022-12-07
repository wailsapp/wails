//go:build windows

package win32

type Menu HMENU
type PopupMenu Menu

func CreatePopupMenu() PopupMenu {
	ret, _, _ := procCreatePopupMenu.Call(0, 0, 0, 0)
	return PopupMenu(ret)
}

func (m Menu) Destroy() bool {
	ret, _, _ := procDestroyMenu.Call(uintptr(m))
	return ret != 0
}

func (p PopupMenu) Destroy() bool {
	return Menu(p).Destroy()
}

func (p PopupMenu) Track(flags uint, x, y int, wnd HWND) bool {
	ret, _, _ := procTrackPopupMenu.Call(
		uintptr(p),
		uintptr(flags),
		uintptr(x),
		uintptr(y),
		0,
		uintptr(wnd),
		0,
	)
	return ret != 0
}

func (p PopupMenu) Append(flags uintptr, id uintptr, text string) bool {
	return Menu(p).Append(flags, id, text)
}

func (m Menu) Append(flags uintptr, id uintptr, text string) bool {
	ret, _, _ := procAppendMenuW.Call(
		uintptr(m),
		flags,
		id,
		MustStringToUTF16uintptr(text),
	)
	return ret != 0
}

func (p PopupMenu) Check(id uintptr, checked bool) bool {
	return Menu(p).Check(id, checked)
}

func (m Menu) Check(id uintptr, check bool) bool {
	var checkState uint = MF_UNCHECKED
	if check {
		checkState = MF_CHECKED
	}
	return CheckMenuItem(HMENU(m), id, checkState) != 0
}

func (m Menu) CheckRadio(startID int, endID int, selectedID int) bool {
	ret, _, _ := procCheckMenuRadioItem.Call(
		uintptr(m),
		uintptr(startID),
		uintptr(endID),
		uintptr(selectedID),
		MF_BYCOMMAND)
	return ret != 0
}

func CheckMenuItem(menu HMENU, id uintptr, flags uint) uint {
	ret, _, _ := procCheckMenuItem.Call(
		uintptr(menu),
		id,
		uintptr(flags),
	)
	return uint(ret)
}

func (p PopupMenu) CheckRadio(startID, endID, selectedID int) bool {
	return Menu(p).CheckRadio(startID, endID, selectedID)
}
