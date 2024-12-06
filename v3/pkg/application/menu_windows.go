//go:build windows

package application

import (
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsMenu struct {
	menu         *Menu
	parentWindow *windowsWebviewWindow

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
			if item.accelerator != nil && item.callback != nil {
				if w.parentWindow != nil {
					w.parentWindow.parent.removeMenuBinding(item.accelerator)
				} else {
					globalApplication.removeKeyBinding(item.accelerator.String())
				}
			}
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
		if item.itemType == radio {
			flags = flags | w32.MFT_RADIOCHECK
		}

		if item.submenu != nil {
			flags = flags | w32.MF_POPUP
			newSubmenu := w32.CreateMenu()
			w.processMenu(newSubmenu, item.submenu)
			itemID = int(newSubmenu)
		}

		thisText := item.Label()
		if item.accelerator != nil && item.callback != nil {
			if w.parentWindow != nil {
				w.parentWindow.parent.addMenuBinding(item.accelerator, item)
			} else {
				globalApplication.addKeyBinding(item.accelerator.String(), func(w *WebviewWindow) {
					item.handleClick()
				})
			}
			thisText = thisText + "\t" + item.accelerator.String()
		}
		var menuText = w32.MustStringToUTF16Ptr(thisText)

		w32.AppendMenu(parentMenu, flags, uintptr(itemID), menuText)
		if item.bitmap != nil {
			w32.SetMenuIcons(parentMenu, itemID, item.bitmap, nil)
		}
	}
}

func (w *windowsMenu) ShowAtCursor() {
	InvokeSync(func() {
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

func DefaultApplicationMenu() *Menu {
	menu := NewMenu()
	menu.AddRole(FileMenu)
	menu.AddRole(EditMenu)
	menu.AddRole(ViewMenu)
	menu.AddRole(WindowMenu)
	menu.AddRole(HelpMenu)
	return menu
}
