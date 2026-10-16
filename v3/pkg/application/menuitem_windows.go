//go:build windows && !server

package application

import (
	"fmt"
	"syscall"
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

	// detachedSubmenus is shared with the builder that owns this item. A
	// submenu is detached while its parent row is hidden, so destroying the
	// root HMENU will not release it.
	detachedSubmenus detachedSubmenuSet

	// bitmap holds the HBITMAP handle installed by the most recent
	// setBitmap call so it can be released before a new one is installed.
	bitmap w32.HBITMAP
}

func (m *windowsMenuItem) setHidden(hidden bool) {
	if hidden && !m.hidden {
		m.hidden = true
		// Remove from parent menu
		if w32.RemoveMenu(m.hMenu, m.id, w32.MF_BYCOMMAND) {
			m.detachedSubmenus.add(m.submenu)
		}
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
		if w32.InsertMenuItem(m.hMenu, uint32(pos), true, m.getMenuInfo()) {
			m.detachedSubmenus.remove(m.submenu)
		}
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

type detachedSubmenuSet map[w32.HMENU]struct{}

func (s detachedSubmenuSet) add(submenu w32.HMENU) {
	if submenu != 0 && s != nil {
		s[submenu] = struct{}{}
	}
}

func (s detachedSubmenuSet) remove(submenu w32.HMENU) {
	delete(s, submenu)
}

func (s detachedSubmenuSet) destroyAll() {
	for submenu := range s {
		w32.DestroyMenu(submenu)
		delete(s, submenu)
	}
}

// assignCommandID gives the last already-appended menu item an explicit
// command ID.
func assignCommandID(parentMenu w32.HMENU, id int) error {
	itemCount := w32.GetMenuItemCount(parentMenu)
	if itemCount <= 0 {
		return fmt.Errorf("GetMenuItemCount returned %d: %v", itemCount, syscall.GetLastError())
	}

	var mii w32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = w32.MIIM_ID
	mii.WID = uint32(id)
	position := uint32(itemCount - 1)
	if !w32.SetMenuItemInfo(parentMenu, position, true, &mii) {
		return fmt.Errorf("SetMenuItemInfo failed at position %d: %v", position, syscall.GetLastError())
	}
	return nil
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

	handles, err := w32.SetMenuIcons(m.hMenu, m.id, bitmap, nil)
	if err != nil {
		globalApplication.error("unable to set bitmap on menu item: %w", err)
		return
	}
	// Release the previous HBITMAP, if any, before replacing it.
	if m.bitmap != 0 {
		w32.DeleteObject(w32.HGDIOBJ(m.bitmap))
	}
	m.bitmap = 0
	if len(handles) > 0 {
		m.bitmap = handles[0]
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
