package systray

import (
	"errors"
	"github.com/wailsapp/wails/v2/internal/platform/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type PopupMenu struct {
	menu        win32.HMENU
	parent      win32.HWND
	menuMapping map[int]*menu.MenuItem
	menuData    *menu.Menu
}

func buildMenu(parentMenu win32.HMENU, inputMenu *menu.Menu) (map[int]*menu.MenuItem, error) {
	menuMapping := make(map[int]*menu.MenuItem)
	for index, item := range inputMenu.Items {
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

		menuMapping[itemID] = item
		ret = win32.AppendMenu(parentMenu, uintptr(flags), uintptr(itemID), item.Label)
		if ret == false {
			return nil, errors.New("AppendMenu failed")
		}
	}
	return menuMapping, nil
}

func (p *PopupMenu) Update() error {
	var err error
	p.menu = win32.CreatePopupMenu()
	p.menuMapping, err = buildMenu(p.menu, p.menuData)
	return err
}

func NewPopupMenu(parent win32.HWND, inputMenu *menu.Menu) (*PopupMenu, error) {
	result := &PopupMenu{
		parent:   parent,
		menuData: inputMenu,
	}
	err := result.Update()
	return result, err
}

func (p *PopupMenu) ShowAtCursor() error {
	x, y, ok := win32.GetCursorPos()
	if ok == false {
		return errors.New("GetCursorPos failed")
	}

	if win32.SetForegroundWindow(p.parent) == false {
		return errors.New("SetForegroundWindow failed")
	}

	if win32.TrackPopupMenu(p.menu, win32.TPM_LEFTALIGN, x, y-5, p.parent) == false {
		return errors.New("TrackPopupMenu failed")
	}

	if win32.PostMessage(p.parent, win32.WM_NULL, 0, 0) == 0 {
		return errors.New("PostMessage failed")
	}

	return nil
}

func (p *PopupMenu) ProcessCommand(cmdMsgID int) {
	item := p.menuMapping[cmdMsgID]
	if item != nil {
		item.Click(&menu.CallbackData{MenuItem: item})
	}
}

func (p *PopupMenu) Destroy() {
	win32.DestroyMenu(p.menu)
}
