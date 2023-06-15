//go:build linux

package application

import (
	"fmt"
	"runtime"
)

type linuxMenuItem struct {
	menuItem  *MenuItem
	native    pointer
	handlerId uint
}

func (l linuxMenuItem) setTooltip(tooltip string) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		menuItemSetToolTip(l.native, tooltip)
	})
}

func (l linuxMenuItem) blockSignal() {
	if l.handlerId != 0 {
		menuItemSignalBlock(l.native, l.handlerId, true)
	}
}

func (l linuxMenuItem) unBlockSignal() {
	if l.handlerId != 0 {
		menuItemSignalBlock(l.native, l.handlerId, false)
	}
}

func (l linuxMenuItem) setLabel(s string) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		menuItemSetLabel(l.native, s)
	})
}

func (l linuxMenuItem) isChecked() bool {
	return menuItemChecked(l.native)
}

func (l linuxMenuItem) setDisabled(disabled bool) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		menuItemSetDisabled(l.native, disabled)
	})
}

func (l linuxMenuItem) setChecked(checked bool) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		menuItemSetChecked(l.native, checked)
	})
}

func (l linuxMenuItem) setHidden(hidden bool) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		widgetSetVisible(l.native, hidden)
	})
}

func (l linuxMenuItem) setAccelerator(accelerator *accelerator) {
	fmt.Println("setAccelerator", accelerator)
	// Set the keyboard shortcut of the menu item
	//	var modifier C.int
	//	var key *C.char
	if accelerator != nil {
		//		modifier = C.int(toMacModifier(accelerator.Modifiers))
		//		key = C.CString(accelerator.Key)
	}

	// Convert the key to a string
	//	C.setMenuItemKeyEquivalent(m.nsMenuItem, key, modifier)
}

func newMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
	}
	switch item.itemType {
	case text:
		result.native = menuItemNew(item.label)

	case checkbox:
		result.native = menuCheckItemNew(item.label)
		result.setChecked(item.checked)
		if item.accelerator != nil {
			result.setAccelerator(item.accelerator)
		}
	case submenu:
		result.native = menuItemNew(item.label)

	default:
		panic(fmt.Sprintf("Unknown menu type: %v", item.itemType))
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newRadioItemImpl(item *MenuItem, group GSListPointer) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuRadioItemNew(group, item.label),
	}
	result.setChecked(item.checked)
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newSpeechMenu() *MenuItem {
	speechMenu := NewMenu()
	speechMenu.Add("Start Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.").
		OnClick(func(ctx *Context) {
			//			C.startSpeaking()
		})
	speechMenu.Add("Stop Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,").
		OnClick(func(ctx *Context) {
			//			C.stopSpeaking()
		})
	subMenu := newSubMenuItem("Speech")
	subMenu.submenu = speechMenu
	return subMenu
}

func newHideMenuItem() *MenuItem {
	return newMenuItem("Hide " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+h").
		OnClick(func(ctx *Context) {

			//			C.hideApplication()
		})
}

func newHideOthersMenuItem() *MenuItem {
	return newMenuItem("Hide Others").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+h").
		OnClick(func(ctx *Context) {
			//			C.hideOthers()
		})
}

func newUnhideMenuItem() *MenuItem {
	return newMenuItem("Show All").
		OnClick(func(ctx *Context) {
			//			C.showAll()
		})
}

func newUndoMenuItem() *MenuItem {
	return newMenuItem("Undo").
		SetAccelerator("CmdOrCtrl+z").
		OnClick(func(ctx *Context) {
			//			C.undo()
		})
}

// newRedoMenuItem creates a new menu item for redoing the last action
func newRedoMenuItem() *MenuItem {
	return newMenuItem("Redo").
		SetAccelerator("CmdOrCtrl+Shift+z").
		OnClick(func(ctx *Context) {
			//			C.redo()
		})
}

func newCutMenuItem() *MenuItem {
	return newMenuItem("Cut").
		SetAccelerator("CmdOrCtrl+x").
		OnClick(func(ctx *Context) {
			//			C.cut()
		})
}

func newCopyMenuItem() *MenuItem {
	return newMenuItem("Copy").
		SetAccelerator("CmdOrCtrl+c").
		OnClick(func(ctx *Context) {
			//			C.copy()
		})
}

func newPasteMenuItem() *MenuItem {
	return newMenuItem("Paste").
		SetAccelerator("CmdOrCtrl+v").
		OnClick(func(ctx *Context) {
			//			C.paste()
		})
}

func newPasteAndMatchStyleMenuItem() *MenuItem {
	return newMenuItem("Paste and Match Style").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v").
		OnClick(func(ctx *Context) {
			//			C.pasteAndMatchStyle()
		})
}

func newDeleteMenuItem() *MenuItem {
	return newMenuItem("Delete").
		SetAccelerator("backspace").
		OnClick(func(ctx *Context) {
			//			C.delete()
		})
}

func newQuitMenuItem() *MenuItem {
	return newMenuItem("Quit " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+q").
		OnClick(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func newSelectAllMenuItem() *MenuItem {
	return newMenuItem("Select All").
		SetAccelerator("CmdOrCtrl+a").
		OnClick(func(ctx *Context) {
			//			C.selectAll()
		})
}

func newAboutMenuItem() *MenuItem {
	return newMenuItem("About " + globalApplication.options.Name).
		OnClick(func(ctx *Context) {
			globalApplication.ShowAboutDialog()
		})
}

func newCloseMenuItem() *MenuItem {
	return newMenuItem("Close").
		SetAccelerator("CmdOrCtrl+w").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Close()
			}
		})
}

func newReloadMenuItem() *MenuItem {
	return newMenuItem("Reload").
		SetAccelerator("CmdOrCtrl+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Reload()
			}
		})
}

func newForceReloadMenuItem() *MenuItem {
	return newMenuItem("Force Reload").
		SetAccelerator("CmdOrCtrl+Shift+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ForceReload()
			}
		})
}

func newToggleFullscreenMenuItem() *MenuItem {
	result := newMenuItem("Toggle Full Screen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleFullscreen()
			}
		})
	if runtime.GOOS == "darwin" {
		result.SetAccelerator("Ctrl+Command+F")
	} else {
		result.SetAccelerator("F11")
	}
	return result
}

func newToggleDevToolsMenuItem() *MenuItem {
	return newMenuItem("Toggle Developer Tools").
		SetAccelerator("Alt+Command+I").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleDevTools()
			}
		})
}

func newZoomResetMenuItem() *MenuItem {
	// reset zoom menu item
	return newMenuItem("Actual Size").
		SetAccelerator("CmdOrCtrl+0").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomReset()
			}
		})
}

func newZoomInMenuItem() *MenuItem {
	return newMenuItem("Zoom In").
		SetAccelerator("CmdOrCtrl+plus").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomIn()
			}
		})
}

func newZoomOutMenuItem() *MenuItem {
	return newMenuItem("Zoom Out").
		SetAccelerator("CmdOrCtrl+-").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomOut()
			}
		})
}

func newMinimizeMenuItem() *MenuItem {
	return newMenuItem("Minimize").
		SetAccelerator("CmdOrCtrl+M").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Minimise()
			}
		})
}

func newZoomMenuItem() *MenuItem {
	return newMenuItem("Zoom").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Zoom()
			}
		})
}

func newFullScreenMenuItem() *MenuItem {
	return newMenuItem("Fullscreen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Fullscreen()
			}
		})
}
