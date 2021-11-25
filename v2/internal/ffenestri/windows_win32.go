//go:build windows
// +build windows

package ffenestri

import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/menumanager"
	"golang.org/x/sys/windows"
)

var (
	// DLL stuff
	user32                  = windows.NewLazySystemDLL("User32.dll")
	win32CreateMenu         = user32.NewProc("CreateMenu")
	win32DestroyMenu        = user32.NewProc("DestroyMenu")
	win32CreatePopupMenu    = user32.NewProc("CreatePopupMenu")
	win32AppendMenuW        = user32.NewProc("AppendMenuW")
	win32SetMenu            = user32.NewProc("SetMenu")
	win32CheckMenuItem      = user32.NewProc("CheckMenuItem")
	win32GetMenuState       = user32.NewProc("GetMenuState")
	win32CheckMenuRadioItem = user32.NewProc("CheckMenuRadioItem")

	applicationMenu *menumanager.WailsMenu
	menuManager     *menumanager.Manager
)

const MF_BITMAP uint32 = 0x00000004
const MF_CHECKED uint32 = 0x00000008
const MF_DISABLED uint32 = 0x00000002
const MF_ENABLED uint32 = 0x00000000
const MF_GRAYED uint32 = 0x00000001
const MF_MENUBARBREAK uint32 = 0x00000020
const MF_MENUBREAK uint32 = 0x00000040
const MF_OWNERDRAW uint32 = 0x00000100
const MF_POPUP uint32 = 0x00000010
const MF_SEPARATOR uint32 = 0x00000800
const MF_STRING uint32 = 0x00000000
const MF_UNCHECKED uint32 = 0x00000000
const MF_BYCOMMAND uint32 = 0x00000000
const MF_BYPOSITION uint32 = 0x00000400

const WM_SIZE = 5
const WM_GETMINMAXINFO = 36

type Win32Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// ------------------- win32 calls -----------------------

func createWin32Menu() (win32Menu, error) {
	res, _, err := win32CreateMenu.Call()
	if res == 0 {
		return 0, wall(err)
	}
	return win32Menu(res), nil
}

func destroyWin32Menu(menu win32Menu) error {
	res, _, err := win32DestroyMenu.Call(uintptr(menu))
	if res == 0 {
		return wall(err, "Menu:", menu)
	}
	return nil
}

func createWin32PopupMenu() (win32Menu, error) {
	res, _, err := win32CreatePopupMenu.Call()
	if res == 0 {
		return 0, wall(err)
	}
	return win32Menu(res), nil
}

func appendWin32MenuItem(menu win32Menu, flags uintptr, submenuOrID uintptr, label string) error {
	menuText, err := windows.UTF16PtrFromString(label)
	if err != nil {
		return err
	}
	res, _, err := win32AppendMenuW.Call(
		uintptr(menu),
		flags,
		submenuOrID,
		uintptr(unsafe.Pointer(menuText)),
	)
	if res == 0 {
		return wall(err, "Menu", menu, "Flags", flags, "submenuOrID", submenuOrID, "label", label)
	}
	return nil
}

func setWindowMenu(window win32Window, menu win32Menu) error {
	res, _, err := win32SetMenu.Call(uintptr(window), uintptr(menu))
	if res == 0 {
		return wall(err, "window", window, "menu", menu)
	}
	return nil
}

func selectRadioItem(selectedMenuID, startMenuItemID, endMenuItemID win32MenuItemID, parent win32Menu) error {
	res, _, err := win32CheckMenuRadioItem.Call(uintptr(parent), uintptr(startMenuItemID), uintptr(endMenuItemID), uintptr(selectedMenuID), uintptr(MF_BYCOMMAND))
	if int(res) == 0 {
		return wall(err, selectedMenuID, startMenuItemID, endMenuItemID, parent)
	}
	return nil
}

//
//func getWindowRect(window win32Window) (*Win32Rect, error) {
//	var windowRect Win32Rect
//	res, _, err := win32GetWindowRect.Call(uintptr(window), uintptr(unsafe.Pointer(&windowRect)))
//	if res == 0 {
//		return nil, err
//	}
//	return &windowRect, nil
//}
//
//func getClientRect(window win32Window) (*Win32Rect, error) {
//	var clientRect Win32Rect
//	res, _, err := win32GetClientRect.Call(uintptr(window), uintptr(unsafe.Pointer(&clientRect)))
//	if res == 0 {
//		return nil, err
//	}
//	return &clientRect, nil
//}
//
