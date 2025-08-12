//go:build windows

package application

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsMenuItem struct {
	parent   *Menu
	menuItem *MenuItem

	hMenu    w32.HMENU
	id       int
	label    string
	disabled bool
	checked  bool
	itemType menuItemType
	hidden   bool
	submenu  w32.HMENU
}

func (m *windowsMenuItem) setHidden(hidden bool) {
	if hidden && !m.hidden {
		m.hidden = true
		// Remove from parent menu
		w32.RemoveMenu(m.hMenu, m.id, w32.MF_BYCOMMAND)
	} else if !hidden && m.hidden {
		m.hidden = false
		// Reinsert into parent menu at correct visible position
		var pos int
		for _, item := range m.parent.items {
			if item == m.menuItem {
				break
			}
			if item.hidden == false {
				pos++
			}
		}
		w32.InsertMenuItem(m.hMenu, uint32(pos), true, m.getMenuInfo())
	}
}

func (m *windowsMenuItem) Checked() bool {
	return m.checked
}

func (m *windowsMenuItem) IsSeparator() bool {
	return m.itemType == separator
}

func (m *windowsMenuItem) IsCheckbox() bool {
	return m.itemType == checkbox
}

func (m *windowsMenuItem) IsRadio() bool {
	return m.itemType == radio
}

func (m *windowsMenuItem) Enabled() bool {
	return !m.disabled
}

func (m *windowsMenuItem) update() {
	w32.SetMenuItemInfo(m.hMenu, uint32(m.id), false, m.getMenuInfo())
}

func (m *windowsMenuItem) setLabel(label string) {
	m.label = label
	m.update()
}

func (m *windowsMenuItem) setDisabled(disabled bool) {
	m.disabled = disabled
	m.update()
}

func (m *windowsMenuItem) setChecked(checked bool) {
	m.checked = checked
	m.update()
}

func (m *windowsMenuItem) destroy() {
	w32.RemoveMenu(m.hMenu, m.id, w32.MF_BYCOMMAND)
}

func (m *windowsMenuItem) setAccelerator(accelerator *accelerator) {
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

func (m *windowsMenuItem) setBitmap(bitmap []byte) {
	if m.menuItem.bitmap == nil {
		return
	}

	// Set the icon
	err := w32.SetMenuIcons(m.hMenu, m.id, bitmap, nil)
	if err != nil {
		globalApplication.error("unable to set bitmap on menu item: %w", err)
		return
	}
	m.update()
}

func newMenuItemImpl(item *MenuItem, parentMenu w32.HMENU, ID int) *windowsMenuItem {
	result := &windowsMenuItem{
		menuItem: item,
		hMenu:    parentMenu,
		id:       ID,
		disabled: item.disabled,
		checked:  item.checked,
		itemType: item.itemType,
		label:    item.label,
		hidden:   item.hidden,
	}

	return result
}

func (m *windowsMenuItem) setTooltip(_ string) {
	// Unsupported
}

func (m *windowsMenuItem) getMenuInfo() *w32.MENUITEMINFO {
	var mii w32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = w32.MIIM_FTYPE | w32.MIIM_ID | w32.MIIM_STATE | w32.MIIM_STRING
	if m.IsSeparator() {
		mii.FType = w32.MFT_SEPARATOR
	} else {
		mii.FType = w32.MFT_STRING
		if m.IsRadio() {
			mii.FType |= w32.MFT_RADIOCHECK
		}
		thisText := m.label
		if m.menuItem.accelerator != nil {
			thisText += "\t" + m.menuItem.accelerator.String()
		}
		mii.DwTypeData = w32.MustStringToUTF16Ptr(thisText)
		mii.Cch = uint32(len([]rune(thisText)))
	}
	mii.WID = uint32(m.id)
	if m.Enabled() {
		mii.FState &^= w32.MFS_DISABLED
	} else {
		mii.FState |= w32.MFS_DISABLED
	}

	if m.IsCheckbox() || m.IsRadio() {
		mii.FMask |= w32.MIIM_CHECKMARKS
	}
	if m.Checked() {
		mii.FState |= w32.MFS_CHECKED
	}

	if m.menuItem.submenu != nil {
		mii.FMask |= w32.MIIM_SUBMENU
		mii.HSubMenu = m.submenu
	}
	return &mii
}
