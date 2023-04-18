//go:build windows

package application

import (
	"unsafe"
)

type windowsMenuItem struct {
	menuItem *MenuItem

	menuItemImpl unsafe.Pointer
}

func (m windowsMenuItem) setTooltip(tooltip string) {
	//C.setMenuItemTooltip(m.nsMenuItem, C.CString(tooltip))
}

func (m windowsMenuItem) setLabel(s string) {
	//C.setMenuItemLabel(m.nsMenuItem, C.CString(s))
}

func (m windowsMenuItem) setDisabled(disabled bool) {
	//C.setMenuItemDisabled(m.nsMenuItem, C.bool(disabled))
}

func (m windowsMenuItem) setChecked(checked bool) {
	//C.setMenuItemChecked(m.nsMenuItem, C.bool(checked))
}

func (m windowsMenuItem) setAccelerator(accelerator *accelerator) {
	//// Set the keyboard shortcut of the menu item
	//var modifier C.int
	//var key *C.char
	//if accelerator != nil {
	//	modifier = C.int(toMacModifier(accelerator.Modifiers))
	//	key = C.CString(accelerator.Key)
	//}
	//
	//// Convert the key to a string
	//C.setMenuItemKeyEquivalent(m.nsMenuItem, key, modifier)
}

func newMenuItemImpl(item *MenuItem) *windowsMenuItem {
	result := &windowsMenuItem{
		menuItem: item,
	}
	//
	//switch item.itemType {
	//case text, checkbox, submenu, radio:
	//	result.nsMenuItem = unsafe.Pointer(C.newMenuItem(C.uint(item.id), C.CString(item.label), C.bool(item.disabled), C.CString(item.tooltip)))
	//	if item.itemType == checkbox || item.itemType == radio {
	//		C.setMenuItemChecked(result.nsMenuItem, C.bool(item.checked))
	//	}
	//	if item.accelerator != nil {
	//		result.setAccelerator(item.accelerator)
	//	}
	//default:
	//	panic("WTF")
	//}
	return result
}

func newSpeechMenu() *MenuItem {
	panic("implement me")
}

func newHideMenuItem() *MenuItem {
	panic("implement me")

}

func newHideOthersMenuItem() *MenuItem {
	panic("implement me")

}

func newUnhideMenuItem() *MenuItem {
	panic("implement me")

}

func newUndoMenuItem() *MenuItem {
	panic("implement me")

}

// newRedoMenuItem creates a new menu item for redoing the last action
func newRedoMenuItem() *MenuItem {
	panic("implement me")

}

func newCutMenuItem() *MenuItem {
	panic("implement me")

}

func newCopyMenuItem() *MenuItem {
	panic("implement me")

}

func newPasteMenuItem() *MenuItem {
	panic("implement me")

}

func newPasteAndMatchStyleMenuItem() *MenuItem {
	panic("implement me")

}

func newDeleteMenuItem() *MenuItem {
	panic("implement me")

}

func newQuitMenuItem() *MenuItem {
	panic("implement me")

}

func newSelectAllMenuItem() *MenuItem {
	panic("implement me")

}

func newAboutMenuItem() *MenuItem {
	panic("implement me")

}

func newCloseMenuItem() *MenuItem {
	panic("implement me")

}

func newReloadMenuItem() *MenuItem {
	panic("implement me")

}

func newForceReloadMenuItem() *MenuItem {
	panic("implement me")

}

func newToggleFullscreenMenuItem() *MenuItem {
	panic("implement me")

}

func newToggleDevToolsMenuItem() *MenuItem {
	panic("implement me")
}

func newZoomResetMenuItem() *MenuItem {
	panic("implement me")

}

func newZoomInMenuItem() *MenuItem {
	panic("implement me")

}

func newZoomOutMenuItem() *MenuItem {
	panic("implement me")

}

func newMinimizeMenuItem() *MenuItem {
	panic("implement me")
}

func newZoomMenuItem() *MenuItem {
	panic("implement me")
}
