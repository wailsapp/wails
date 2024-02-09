package application

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

type menuItemType int

const (
	text menuItemType = iota
	separator
	checkbox
	radio
	submenu
)

var menuItemID uintptr
var menuItemMap = make(map[uint]*MenuItem)
var menuItemMapLock sync.Mutex

func addToMenuItemMap(menuItem *MenuItem) {
	menuItemMapLock.Lock()
	menuItemMap[menuItem.id] = menuItem
	menuItemMapLock.Unlock()
}

func getMenuItemByID(id uint) *MenuItem {
	menuItemMapLock.Lock()
	defer menuItemMapLock.Unlock()
	return menuItemMap[id]
}

type menuItemImpl interface {
	setTooltip(s string)
	setLabel(s string)
	setDisabled(disabled bool)
	setChecked(checked bool)
	setAccelerator(accelerator *accelerator)
	setHidden(hidden bool)
	setBitmap(bitmap []byte)
}

type MenuItem struct {
	id              uint
	label           string
	tooltip         string
	disabled        bool
	checked         bool
	hidden          bool
	bitmap          []byte
	submenu         *Menu
	callback        func(*Context)
	itemType        menuItemType
	accelerator     *accelerator
	role            Role
	contextMenuData *ContextMenuData

	impl              menuItemImpl
	radioGroupMembers []*MenuItem
}

func newMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		itemType: text,
	}
	addToMenuItemMap(result)
	return result
}

func newMenuItemSeparator() *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		itemType: separator,
	}
	return result
}

func newMenuItemCheckbox(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		checked:  checked,
		itemType: checkbox,
	}
	addToMenuItemMap(result)
	return result
}

func newMenuItemRadio(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		checked:  checked,
		itemType: radio,
	}
	addToMenuItemMap(result)
	return result
}

func newSubMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		itemType: submenu,
		submenu: &Menu{
			label: label,
		},
	}
	addToMenuItemMap(result)
	return result
}

func newRole(role Role) *MenuItem {
	switch role {
	case AppMenu:
		return newAppMenu()
	case EditMenu:
		return newEditMenu()
	case FileMenu:
		return newFileMenu()
	case ViewMenu:
		return newViewMenu()
	case ServicesMenu:
		return newServicesMenu()
	case SpeechMenu:
		return newSpeechMenu()
	case WindowMenu:
		return newWindowMenu()
	case HelpMenu:
		return newHelpMenu()
	case Hide:
		return newHideMenuItem()
	case HideOthers:
		return newHideOthersMenuItem()
	case UnHide:
		return newUnhideMenuItem()
	case Undo:
		return newUndoMenuItem()
	case Redo:
		return newRedoMenuItem()
	case Cut:
		return newCutMenuItem()
	case Copy:
		return newCopyMenuItem()
	case Paste:
		return newPasteMenuItem()
	case PasteAndMatchStyle:
		return newPasteAndMatchStyleMenuItem()
	case SelectAll:
		return newSelectAllMenuItem()
	case Delete:
		return newDeleteMenuItem()
	case Quit:
		return newQuitMenuItem()
	case Close:
		return newCloseMenuItem()
	case About:
		return newAboutMenuItem()
	case Reload:
		return newReloadMenuItem()
	case ForceReload:
		return newForceReloadMenuItem()
	case ToggleFullscreen:
		return newToggleFullscreenMenuItem()
	case ShowDevTools:
		return newShowDevToolsMenuItem()
	case ResetZoom:
		return newZoomResetMenuItem()
	case ZoomIn:
		return newZoomInMenuItem()
	case ZoomOut:
		return newZoomOutMenuItem()
	case Minimize:
		return newMinimizeMenuItem()
	case Zoom:
		return newZoomMenuItem()
	case FullScreen:
		return newFullScreenMenuItem()

	default:
		globalApplication.error(fmt.Sprintf("No support for role: %v", role))
		os.Exit(1)
	}
	return nil
}

func newServicesMenu() *MenuItem {
	serviceMenu := newSubMenuItem("Services")
	serviceMenu.role = ServicesMenu
	return serviceMenu
}

func (m *MenuItem) handleClick() {
	var ctx = newContext().
		withClickedMenuItem(m).
		withContextMenuData(m.contextMenuData)
	if m.itemType == checkbox {
		m.checked = !m.checked
		ctx.withChecked(m.checked)
		if m.impl != nil {
			m.impl.setChecked(m.checked)
		}
	}
	if m.itemType == radio {
		for _, member := range m.radioGroupMembers {
			member.checked = false
			if member.impl != nil {
				member.impl.setChecked(false)
			}
		}
		m.checked = true
		ctx.withChecked(true)
		if m.impl != nil {
			m.impl.setChecked(true)
		}
	}
	if m.callback != nil {
		go m.callback(ctx)
	}
}

func (m *MenuItem) SetAccelerator(shortcut string) *MenuItem {
	accelerator, err := parseAccelerator(shortcut)
	if err != nil {
		globalApplication.error("invalid accelerator:", err.Error())
		return m
	}
	m.accelerator = accelerator
	if m.impl != nil {
		m.impl.setAccelerator(accelerator)
	}
	return m
}

func (m *MenuItem) SetTooltip(s string) *MenuItem {
	m.tooltip = s
	if m.impl != nil {
		m.impl.setTooltip(s)
	}
	return m
}

func (m *MenuItem) SetLabel(s string) *MenuItem {
	m.label = s
	if m.impl != nil {
		m.impl.setLabel(s)
	}
	return m
}

func (m *MenuItem) SetEnabled(enabled bool) *MenuItem {
	m.disabled = !enabled
	if m.impl != nil {
		m.impl.setDisabled(m.disabled)
	}
	return m
}

func (m *MenuItem) SetBitmap(bitmap []byte) *MenuItem {
	m.bitmap = bitmap
	if m.impl != nil {
		m.impl.setBitmap(bitmap)
	}
	return m
}

func (m *MenuItem) SetChecked(checked bool) *MenuItem {
	m.checked = checked
	if m.impl != nil {
		m.impl.setChecked(m.checked)
	}
	return m
}

func (m *MenuItem) SetHidden(hidden bool) *MenuItem {
	m.hidden = hidden
	if m.impl != nil {
		m.impl.setHidden(m.hidden)
	}
	return m
}

func (m *MenuItem) Checked() bool {
	return m.checked
}

func (m *MenuItem) IsSeparator() bool {
	return m.itemType == separator
}

func (m *MenuItem) IsSubmenu() bool {
	return m.itemType == submenu
}

func (m *MenuItem) IsCheckbox() bool {
	return m.itemType == checkbox
}

func (m *MenuItem) IsRadio() bool {
	return m.itemType == radio
}

func (m *MenuItem) Hidden() bool {
	return m.hidden
}

func (m *MenuItem) OnClick(f func(*Context)) *MenuItem {
	m.callback = f
	return m
}

func (m *MenuItem) Label() string {
	return m.label
}

func (m *MenuItem) Tooltip() string {
	return m.tooltip
}

func (m *MenuItem) Enabled() bool {
	return !m.disabled
}

func (m *MenuItem) setContextData(data *ContextMenuData) {
	m.contextMenuData = data
	if m.submenu != nil {
		m.submenu.setContextData(data)
	}
}
