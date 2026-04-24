//go:build windows && !server

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

	// bitmaps tracks HBITMAP handles allocated by SetMenuIcons during
	// processMenu so they can be released when the menu is rebuilt.
	bitmaps []w32.HBITMAP
}

func (w *windowsMenu) freeBitmaps() {
	for _, h := range w.bitmaps {
		w32.DeleteObject(w32.HGDIOBJ(h))
	}
	w.bitmaps = nil
	// HBITMAPs allocated at runtime via MenuItem.SetBitmap live on the
	// windowsMenuItem impl, not in w.bitmaps. Walk the menuMapping to
	// release them before DestroyMenu — otherwise every rebuild leaks
	// one HBITMAP per item that had a runtime SetBitmap call.
	for _, item := range w.menuMapping {
		impl, ok := item.impl.(*windowsMenuItem)
		if !ok || impl.bitmap == 0 {
			continue
		}
		w32.DeleteObject(w32.HGDIOBJ(impl.bitmap))
		impl.bitmap = 0
	}
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
		w.freeBitmaps()
		w32.DestroyMenu(w.hMenu)
	}
	w.hMenu = w32.NewPopupMenu()
	// menuMapping and currentMenuID index items built during processMenu.
	// Without resetting, repeated update calls accumulate stale *MenuItem
	// entries keyed by IDs that reference HMENUs that no longer exist.
	w.menuMapping = make(map[int]*MenuItem)
	w.currentMenuID = 0
	w.processMenu(w.hMenu, w.menu)
}

func (w *windowsMenu) processMenu(parentMenu w32.HMENU, inputMenu *Menu) {
	for _, item := range inputMenu.items {
		w.currentMenuID++
		itemID := w.currentMenuID
		w.menuMapping[itemID] = item

		menuItemImpl := newMenuItemImpl(item, parentMenu, itemID)
		menuItemImpl.parent = inputMenu
		item.impl = menuItemImpl

		if item.Hidden() {
			if item.accelerator != nil && item.callback != nil {
				if w.parentWindow != nil {
					w.parentWindow.parent.removeMenuBinding(item.accelerator)
				} else {
					globalApplication.KeyBinding.Remove(item.accelerator.String())
				}
			}
		}

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
				globalApplication.KeyBinding.Add(item.accelerator.String(), func(w Window) {
					item.handleClick()
				})
			}
			thisText = thisText + "\t" + item.accelerator.String()
		}
		var menuText = w32.MustStringToUTF16Ptr(thisText)

		// If the item is hidden, don't append
		if item.Hidden() {
			continue
		}

		w32.AppendMenu(parentMenu, flags, uintptr(itemID), menuText)
		if item.bitmap != nil {
			handles, err := w32.SetMenuIcons(parentMenu, itemID, item.bitmap, nil)
			if err != nil {
				globalApplication.error("SetMenuIcons failed: %v", err)
				continue
			}
			w.bitmaps = append(w.bitmaps, handles...)
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
