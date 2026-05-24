//go:build windows && !server

package application

import (
	"fmt"
	"syscall"

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
	releaseMenuBitmaps(w.bitmaps, w.menuMapping)
	w.bitmaps = nil
}

func newMenuImpl(menu *Menu) *windowsMenu {
	result := &windowsMenu{
		menu:        menu,
		menuMapping: make(map[int]*MenuItem),
	}

	return result
}

func (w *windowsMenu) update() {
	// Stage the rebuild into a fresh HMENU and fresh mapping, only swapping
	// the old state out once processMenu returns successfully. See
	// Win32Menu.Update for the item.impl caveat on the failure path.
	newHMENU := w32.NewPopupMenu()
	oldHMENU := w.hMenu
	oldMapping := w.menuMapping
	oldBitmaps := w.bitmaps

	// Transfer runtime SetBitmap handles off the old impls now, before
	// processMenu reassigns item.impl and makes them unreachable via the
	// mapping walk. Every handle lives in oldBitmaps from here on.
	for _, item := range oldMapping {
		if impl, ok := item.impl.(*windowsMenuItem); ok && impl.bitmap != 0 {
			oldBitmaps = append(oldBitmaps, impl.bitmap)
			impl.bitmap = 0
		}
	}

	w.hMenu = newHMENU
	w.menuMapping = make(map[int]*MenuItem)
	w.currentMenuID = 0
	w.bitmaps = nil

	if err := w.processMenu(newHMENU, w.menu); err != nil {
		globalApplication.error("menu rebuild failed, keeping previous menu: %v", err)
		w.freeBitmaps()
		w32.DestroyMenu(newHMENU)
		w.hMenu = oldHMENU
		w.menuMapping = oldMapping
		w.bitmaps = oldBitmaps
		return
	}

	if oldHMENU != 0 {
		releaseMenuBitmaps(oldBitmaps, oldMapping)
		w32.DestroyMenu(oldHMENU)
	}
}

// processMenu populates parentMenu from inputMenu. Any native AppendMenu or
// SetMenuIcons failure returns an error; recursive submenu builds propagate
// the error so the outer update can back out cleanly instead of attaching a
// half-built submenu via MF_POPUP.
func (w *windowsMenu) processMenu(parentMenu w32.HMENU, inputMenu *Menu) error {
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
			if err := w.processMenu(newSubmenu, item.submenu); err != nil {
				// Submenu was allocated but never attached, so the outer
				// DestroyMenu on parentMenu won't reach it. Free it here.
				w32.DestroyMenu(newSubmenu)
				return err
			}
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

		if ok := w32.AppendMenu(parentMenu, flags, uintptr(itemID), menuText); !ok {
			return fmt.Errorf("AppendMenu failed for %q: %v", thisText, syscall.GetLastError())
		}
		if item.bitmap != nil {
			handles, err := w32.SetMenuIcons(parentMenu, itemID, item.bitmap, nil)
			if err != nil {
				return fmt.Errorf("SetMenuIcons failed for %q: %w", thisText, err)
			}
			w.bitmaps = append(w.bitmaps, handles...)
		}
	}
	return nil
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
