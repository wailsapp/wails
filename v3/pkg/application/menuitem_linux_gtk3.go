//go:build linux && !android && gtk3 && !server

package application

import (
	"fmt"
)

type linuxMenuItem struct {
	menuItem  *MenuItem
	owner     *MenuItem
	native    pointer
	handlerId uint
	// The accelerator bound to native, so it can be replaced or removed.
	accelSet  bool
	accelKey  uint32
	accelMods uint32
}

// eachNative propagates changes made through the original item to its live views.
func (l *linuxMenuItem) eachNative(fn func(*linuxMenuItem)) {
	if l.owner == nil || l.menuItem != l.owner {
		if l.native != nil {
			fn(l)
		}
		return
	}
	peers := gtk3ItemPeers(l.owner)
	if len(peers) == 0 && l.native != nil {
		fn(l)
		return
	}
	for _, peer := range peers {
		fn(peer)
	}
}

func (l *linuxMenuItem) withNative(fn func(*linuxMenuItem)) {
	InvokeSync(func() {
		l.eachNative(func(peer *linuxMenuItem) {
			peer.blockSignal()
			defer peer.unBlockSignal()
			fn(peer)
		})
	})
}

func (l *linuxMenuItem) setTooltip(s string) {
	l.withNative(func(p *linuxMenuItem) { menuItemSetToolTip(p.native, s) })
}
func (l *linuxMenuItem) setBitmap(b []byte) {
	l.withNative(func(p *linuxMenuItem) { menuItemSetBitmap(p.native, b) })
}
func (l *linuxMenuItem) setLabel(s string) {
	l.withNative(func(p *linuxMenuItem) { menuItemSetLabel(p.native, s) })
}
func (l *linuxMenuItem) setDisabled(b bool) {
	l.withNative(func(p *linuxMenuItem) { menuItemSetDisabled(p.native, b) })
}
func (l *linuxMenuItem) setChecked(b bool) {
	l.withNative(func(p *linuxMenuItem) { menuItemSetChecked(p.native, b) })
}
func (l *linuxMenuItem) setHidden(b bool) {
	l.withNative(func(p *linuxMenuItem) { widgetSetVisible(p.native, b) })
	// A hidden item's shortcut does nothing, as on macOS and Windows: GTK3
	// still activates the accelerator of a hidden menu item.
	if l.menuItem != nil && l.menuItem.accelerator != nil {
		if b {
			l.setAccelerator(nil)
		} else {
			l.setAccelerator(l.menuItem.accelerator)
		}
	}
}
func (l *linuxMenuItem) isChecked() bool { return l.native != nil && menuItemChecked(l.native) }
func (l *linuxMenuItem) blockSignal() {
	if l.handlerId != 0 {
		menuItemSignalBlock(l.native, l.handlerId, true)
	}
}
func (l *linuxMenuItem) unBlockSignal() {
	if l.handlerId != 0 {
		menuItemSignalBlock(l.native, l.handlerId, false)
	}
}

func (l *linuxMenuItem) destroy() {
	InvokeSync(func() { l.eachNative(func(p *linuxMenuItem) { p.blockSignal(); menuItemDestroy(p.native) }) })
}

func (l *linuxMenuItem) setAccelerator(accelerator *accelerator) {
	InvokeSync(func() {
		l.eachNative(func(p *linuxMenuItem) {
			p.clearAccelerator()
			if accelerator == nil {
				return
			}
			key := gtk3AcceleratorKeyval(accelerator.Key)
			if key == 0 {
				globalApplication.warning("menu accelerator %q has no GTK key", accelerator.String())
				return
			}
			mods := gtk3AcceleratorModifiers(accelerator.Modifiers)
			menuItemAddAccelerator(p.native, key, mods)
			p.accelSet, p.accelKey, p.accelMods = true, key, mods
		})
	})
}

func (l *linuxMenuItem) clearAccelerator() {
	if l.accelSet && l.native != nil {
		menuItemRemoveAccelerator(l.native, l.accelKey, l.accelMods)
	}
	l.accelSet = false
}

func newMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
	}
	switch item.itemType {
	case text:
		result.native = menuItemNew(item.label, item.bitmap)
		if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
			result.setAccelerator(item.accelerator)
		}

	case checkbox:
		result.native = menuCheckItemNew(item.label, item.bitmap)
		result.setChecked(item.checked)
		if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
			result.setAccelerator(item.accelerator)
		}
	case submenu:
		result.native = menuItemNew(item.label, item.bitmap)

	default:
		panic(fmt.Sprintf("Unknown menu type: %v", item.itemType))
	}
	result.setDisabled(result.menuItem.disabled)
	result.trackDestroy()
	registerGTK3ItemView(result)
	return result
}

func newRadioItemImpl(item *MenuItem, group GSListPointer) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuRadioItemNew(group, item.label),
	}
	if item.accelerator != nil && !item.hidden { // a hidden item's shortcut does nothing
		result.setAccelerator(item.accelerator)
	}
	result.setChecked(item.checked)
	result.setDisabled(result.menuItem.disabled)
	result.trackDestroy()
	registerGTK3ItemView(result)
	return result
}

// GTK destroys descendants when their parent menu item is removed. Clear the
// Go implementation too, including items removed from the Go menu before Update.
func (l *linuxMenuItem) trackDestroy() {
	menuOnDestroy(l.native, func() {
		unregisterGTK3ItemView(l)
		l.native = nil
		l.handlerId = 0
		// A newer native item owns the handler entry; a tray binding does not,
		// so the entry must still go or items removed from a tray menu leak.
		if impl, ok := l.menuItem.impl.(*linuxMenuItem); ok && impl != l {
			return
		}
		if l.menuItem.impl == l {
			l.menuItem.impl = nil
		}
		delete(gtkSignalToMenuItem, l.menuItem.id)
	})
}
