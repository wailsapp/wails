package systray

import (
	"errors"
	"github.com/wailsapp/wails/v2/internal/platform/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type PopupMenu struct {
	menu          win32.PopupMenu
	parent        win32.HWND
	menuMapping   map[int]*menu.MenuItem
	checkboxItems map[*menu.MenuItem][]int
	menuData      *menu.Menu
}

func (p *PopupMenu) buildMenu(parentMenu win32.PopupMenu, inputMenu *menu.Menu, startindex int) error {
	for index, item := range inputMenu.Items {
		if item.Hidden {
			continue
		}
		var ret bool
		itemID := index + startindex
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
		if item.SubMenu != nil {
			flags = flags | win32.MF_POPUP
			submenu := win32.CreatePopupMenu()
			err := p.buildMenu(submenu, item.SubMenu, itemID)
			if err != nil {
				return err
			}
			ret = parentMenu.Append(uintptr(flags), uintptr(submenu), item.Label)
			if ret == false {
				return errors.New("AppendMenu failed")
			}
			continue
		}

		p.menuMapping[itemID] = item
		if item.IsCheckbox() {
			p.checkboxItems[item] = append(p.checkboxItems[item], itemID)
		}
		ret = parentMenu.Append(uintptr(flags), uintptr(itemID), item.Label)
		if ret == false {
			return errors.New("AppendMenu failed")
		}
	}
	return nil
}

func (p *PopupMenu) Update() error {
	p.menu = win32.CreatePopupMenu()
	return p.buildMenu(p.menu, p.menuData, win32.MenuItemMsgID)
}

func NewPopupMenu(parent win32.HWND, inputMenu *menu.Menu) (*PopupMenu, error) {
	result := &PopupMenu{
		parent:        parent,
		menuData:      inputMenu,
		menuMapping:   make(map[int]*menu.MenuItem),
		checkboxItems: make(map[*menu.MenuItem][]int),
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

	if p.menu.Track(win32.TPM_LEFTALIGN, x, y-5, p.parent) == false {
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
		if item.Type == menu.CheckboxType {
			item.Checked = !item.Checked
			for _, menuID := range p.checkboxItems[item] {
				p.menu.Check(uintptr(menuID), item.Checked)
			}
			// TODO: Check duplicate menu items
		}
		if item.Click != nil {
			item.Click(&menu.CallbackData{MenuItem: item})
		}
	}
}

func (p *PopupMenu) Destroy() {
	p.menu.Destroy()
}
