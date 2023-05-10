//go:build windows

package application

import (
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsMenu struct {
	menu *Menu

	hWnd          w32.HWND
	hMenu         w32.HMENU
	currentMenuID int
	menuMapping   map[int]*MenuItem
	checkboxItems []*Menu
}

func newMenuImpl(menu *Menu) *windowsMenu {
	result := &windowsMenu{
		menu:        menu,
		menuMapping: make(map[int]*MenuItem),
	}

	return result
}

func (w *windowsMenu) update() {
	if w.hMenu != 0 {
		w32.DestroyMenu(w.hMenu)
	}
	w.hMenu = w32.NewPopupMenu()
	w.processMenu(w.hMenu, w.menu)
}

func (w *windowsMenu) processMenu(parentMenu w32.HMENU, inputMenu *Menu) {
	for _, item := range inputMenu.items {
		if item.Hidden() {
			continue
		}
		w.currentMenuID++
		itemID := w.currentMenuID
		w.menuMapping[itemID] = item

		flags := uint32(w32.MF_STRING)
		if item.disabled {
			flags = flags | w32.MF_GRAYED
		}
		if item.checked {
			flags = flags | w32.MF_CHECKED
		}
		if item.IsSeparator() {
			flags = flags | w32.MF_SEPARATOR
		}
		//
		//if item.IsCheckbox() {
		//	w.checkboxItems[item] = append(w.checkboxItems[item], itemID)
		//}
		//if item.IsRadio() {
		//	currentRadioGroup.Add(itemID, item)
		//} else {
		//	if len(currentRadioGroup) > 0 {
		//		for _, radioMember := range currentRadioGroup {
		//			currentRadioGroup := currentRadioGroup
		//			p.radioGroups[radioMember.MenuItem] = append(p.radioGroups[radioMember.MenuItem], &currentRadioGroup)
		//		}
		//		currentRadioGroup = RadioGroup{}
		//	}
		//}

		if item.submenu != nil {
			flags = flags | w32.MF_POPUP
			newSubmenu := w32.CreateMenu()
			w.processMenu(newSubmenu, item.submenu)
			itemID = int(newSubmenu)
		}

		var menuText = w32.MustStringToUTF16Ptr(item.Label())

		w32.AppendMenu(parentMenu, flags, uintptr(itemID), menuText)
	}
}

func (w *windowsMenu) ShowAtCursor() {
	invokeSync(func() {
		x, y, ok := w32.GetCursorPos()
		if !ok {
			return
		}
		w.ShowAt(x, y)
	})
}

func (w *windowsMenu) ShowAt(x int, y int) {
	w.update()
	w32.TrackPopupMenuEx(w.hMenu,
		w32.TPM_LEFTALIGN,
		int32(x),
		int32(y),
		w.hWnd,
		nil)
	w32.PostMessage(w.hWnd, w32.WM_NULL, 0, 0)
}

func (w *windowsMenu) ProcessCommand(cmdMsgID int) {
	item := w.menuMapping[cmdMsgID]
	if item == nil {
		return
	}
	item.handleClick()
}

func defaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
