//go:build linux && purego

package application

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

type linuxMenuItem struct {
	menuItem  *MenuItem
	native    uintptr
	handlerId uint
}

func (l linuxMenuItem) setTooltip(tooltip string) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		var setToolTip func(uintptr, string)
		purego.RegisterLibFunc(&setToolTip, gtk, "gtk_widget_set_tooltip_text")

		setToolTip(l.native, tooltip)
	})
}

func (l linuxMenuItem) blockSignal() {
	var block func(uintptr, uint)
	purego.RegisterLibFunc(&block, gtk, "g_signal_handler_block")

	if l.handlerId != 0 {
		block(l.native, l.handlerId)
	}
}

func (l linuxMenuItem) unBlockSignal() {
	var unblock func(uintptr, uint)
	purego.RegisterLibFunc(&unblock, gtk, "g_signal_handler_unblock")

	if l.handlerId != 0 {
		unblock(l.native, l.handlerId)
	}
}

func (l linuxMenuItem) setLabel(s string) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()
		var setLabel func(uintptr, string)
		purego.RegisterLibFunc(&setLabel, gtk, "gtk_menu_item_set_label")
		setLabel(l.native, s)
	})
}

func (l linuxMenuItem) isChecked() bool {
	var getActive func(uintptr) int
	purego.RegisterLibFunc(&getActive, gtk, "gtk_check_menu_item_get_active")

	if getActive(l.native) == 1 {
		return true
	}
	return false
}

func (l linuxMenuItem) setDisabled(disabled bool) {

	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()

		var setSensitive func(uintptr, int)
		purego.RegisterLibFunc(&setSensitive, gtk, "gtk_widget_set_sensitive")

		value := 1
		if disabled {
			value = 0
		}
		setSensitive(l.native, value)
	})
}

func (l linuxMenuItem) setChecked(checked bool) {
	globalApplication.dispatchOnMainThread(func() {
		l.blockSignal()
		defer l.unBlockSignal()

		var setActive func(uintptr, int)
		purego.RegisterLibFunc(&setActive, gtk, "gtk_check_menu_item_set_active")

		value := 0
		if checked {
			value = 1
		}
		setActive(l.native, value)
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
	var newWithLabel func(string) uintptr
	purego.RegisterLibFunc(&newWithLabel, gtk, "gtk_menu_item_new_with_label")
	var newCBWithLabel func(string) uintptr
	purego.RegisterLibFunc(&newCBWithLabel, gtk, "gtk_check_menu_item_new_with_label")

	switch item.itemType {
	case text:
		result.native = newWithLabel(item.label)

	case checkbox:
		result.native = newCBWithLabel(item.label)
		result.setChecked(item.checked)
		if item.accelerator != nil {
			result.setAccelerator(item.accelerator)
		}
	case radio:
		panic("Shouldn't get here with a radio item")

	case submenu:
		result.native = newWithLabel(item.label)

	default:
		panic("WTF")
	}
	result.setDisabled(result.menuItem.disabled)

	return result
}

func newRadioItemImpl(item *MenuItem, group uintptr) *linuxMenuItem {
	var newWithLabel func(uintptr, string) uintptr
	purego.RegisterLibFunc(&newWithLabel, gtk, "gtk_radio_menu_item_new_with_label")

	result := &linuxMenuItem{
		menuItem: item,
		native:   newWithLabel(group, item.label),
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
				currentWindow.Minimize()
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
