//go:build linux && !android && !gtk3 && !server

package application

import (
	"fmt"
)

type linuxMenuItem struct {
	menuItem   *MenuItem
	native     pointer
	handlerId  uint
	parentMenu pointer
	menuIndex  int
	isHidden   bool
}

// menuItemIdAttribute names the GMenuItem attribute that carries a Wails item
// id into the GMenu model. Both GTK backends set it on every native item.
const menuItemIdAttribute = "wails-item-id"

// radioItemTarget is the value a radio group's action state holds when the
// item with this id is the selected one.
func radioItemTarget(itemId uint) string {
	return fmt.Sprintf("%d", itemId)
}

// acceleratorAction names the action an item's accelerator activates. A
// radio item has no action of its own: its group's stateful action does the
// work, so the accelerator activates that action with the item's target, the
// detailed name "radio_group_N::<target>".
func acceleratorAction(itemId uint) (string, bool) {
	if name, ok := lookupActionName(itemId); ok {
		return name, true
	}
	if group, ok := lookupRadioGroup(itemId); ok {
		return group + "::" + radioItemTarget(itemId), true
	}
	return "", false
}

func (l linuxMenuItem) setTooltip(tooltip string) {
	InvokeSync(func() {
		menuItemSetToolTip(l.native, tooltip)
	})
}

func (l linuxMenuItem) destroy() {
	InvokeSync(func() {
		menuItemDestroy(l.native)
	})
}

func (l linuxMenuItem) blockSignal() {
}

func (l linuxMenuItem) setBitmap(data []byte) {
	InvokeSync(func() {
		menuItemSetBitmap(l.native, data)
	})
}

func (l linuxMenuItem) unBlockSignal() {
}

func (l linuxMenuItem) setLabel(s string) {
	InvokeSync(func() {
		menuItemSetLabel(l.native, s)
		// GMenu copied this item when it was appended, so the change above
		// reaches only the original. Replace the model's copy to show it.
		// A hidden item is not in the model; it gets the new label when it is
		// shown again, because that reinserts the original.
		if l.menuItem != nil && l.parentMenu != nilPointer && !l.isHidden {
			menuModelReplaceItem(l.parentMenu, l.native, l.menuItem.id)
		}
	})
}

func (l linuxMenuItem) isChecked() bool {
	return menuItemChecked(l.native)
}

func (l linuxMenuItem) setDisabled(disabled bool) {
	InvokeSync(func() {
		menuItemSetDisabled(l.native, disabled)
	})
}

func (l linuxMenuItem) setChecked(checked bool) {
	InvokeSync(func() {
		menuItemSetChecked(l.native, checked)
	})
}

func (l *linuxMenuItem) setHidden(hidden bool) {
	if l.isHidden == hidden {
		return
	}
	InvokeSync(func() {
		menuItemSetHidden(l, hidden)
		// A hidden item's shortcut does nothing, as on macOS and Windows: its
		// action stays registered, so its accelerator goes and comes back.
		if l.menuItem != nil && l.menuItem.accelerator != nil {
			if hidden {
				clearMenuItemAccelerator(l.menuItem.id)
			} else {
				setMenuItemAccelerator(l.menuItem.id, l.menuItem.accelerator)
			}
		}
	})
	l.isHidden = hidden
}

// setAccelerator binds the shortcut to the action; nil removes the binding.
func (l linuxMenuItem) setAccelerator(accelerator *accelerator) {
	if l.menuItem == nil {
		return
	}
	InvokeSync(func() {
		if accelerator == nil {
			clearMenuItemAccelerator(l.menuItem.id)
			return
		}
		setMenuItemAccelerator(l.menuItem.id, accelerator)
	})
}

func newMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
	}
	switch item.itemType {
	case text:
		result.native = menuItemNewWithId(item.label, item.bitmap, item.id)
	case submenu:
		result.native = menuItemNewWithId(item.label, item.bitmap, item.id)
	default:
		panic(fmt.Sprintf("Unknown menu type for newMenuItemImpl: %v", item.itemType))
	}
	menuItemSetIdAttribute(result.native, item.id)
	if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newCheckMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuCheckItemNewWithId(item.label, item.bitmap, item.id, item.checked),
	}
	menuItemSetIdAttribute(result.native, item.id)
	if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newRadioMenuItemImpl(item *MenuItem, groupId uint, checkedId uint) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuRadioItemNewWithGroup(item.label, item.id, groupId, checkedId),
	}
	menuItemSetIdAttribute(result.native, item.id)
	if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}
