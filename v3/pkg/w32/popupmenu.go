//go:build windows

package w32

type Menu HMENU
type PopupMenu Menu

func (m Menu) destroy() bool {
	ret, _, _ := procDestroyMenu.Call(uintptr(m))
	return ret != 0
}

func (p PopupMenu) destroy() bool {
	return Menu(p).destroy()
}

func (p PopupMenu) Track(hwnd HWND, flags uint32, x, y int32) bool {
	return TrackPopupMenuEx(
		HMENU(p),
		flags,
		x,
		y,
		hwnd,
		nil)
}

func RemoveMenu(m HMENU, pos, flags int) bool {
	ret, _, _ := procRemoveMenu.Call(
		uintptr(m),
		uintptr(pos),
		uintptr(flags))
	return ret != 0
}

func (p PopupMenu) Append(flags uint32, id uintptr, text string) bool {
	return Menu(p).Append(flags, id, text)
}

func (m Menu) Append(flags uint32, id uintptr, text string) bool {
	return AppendMenu(HMENU(m), flags, id, MustStringToUTF16Ptr(text))
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

func CheckRadio(m HMENU, startID int, endID int, selectedID int) bool {
	ret, _, _ := procCheckMenuRadioItem.Call(
		m,
		uintptr(startID),
		uintptr(endID),
		uintptr(selectedID),
		MF_BYCOMMAND)
	return ret != 0
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
		menu,
		id,
		uintptr(flags),
	)
	return uint(ret)
}

func (p PopupMenu) CheckRadio(startID, endID, selectedID int) bool {
	return Menu(p).CheckRadio(startID, endID, selectedID)
}

func NewMenu() HMENU {
	ret, _, _ := procCreateMenu.Call()
	return HMENU(ret)
}

func NewPopupMenu() HMENU {
	ret, _, _ := procCreatePopupMenu.Call()
	return ret
}
