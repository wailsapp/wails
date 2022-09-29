package systray

import (
	"errors"
	"github.com/wailsapp/wails/v2/internal/platform/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

func displayMenu(hwnd win32.HWND, menuItems []*menu.MenuItem) error {
	popupMenu := win32.CreatePopupMenu()

	for index, item := range menuItems {
		var ret bool
		itemID := win32.MenuItemMsgID + index
		flags := win32.MF_STRING
		if item.Disabled {
			flags = flags | win32.MF_GRAYED
		}
		if item.Checked {
			flags = flags | win32.MF_CHECKED
		}
		//if item.BarBreak {
		//	flags = flags | win32.MF_MENUBARBREAK
		//}
		if item.IsSeparator() {
			flags = flags | win32.MF_SEPARATOR
		}

		ret = win32.AppendMenu(popupMenu, uintptr(flags), uintptr(itemID), item.Label)
		if ret == false {
			return errors.New("AppendMenu failed")
		}
	}

	x, y, ok := win32.GetCursorPos()
	if ok == false {
		return errors.New("GetCursorPos failed")
	}

	if win32.SetForegroundWindow(hwnd) == false {
		return errors.New("SetForegroundWindow failed")
	}

	if win32.TrackPopupMenu(popupMenu, win32.TPM_LEFTALIGN, x, y-5, hwnd) == false {
		return errors.New("TrackPopupMenu failed")
	}

	if win32.PostMessage(hwnd, win32.WM_NULL, 0, 0) == 0 {
		return errors.New("PostMessage failed")
	}

	return nil
}
