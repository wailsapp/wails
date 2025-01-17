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

func NewMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		itemType: text,
	}
	addToMenuItemMap(result)
	return result
}

func NewMenuItemSeparator() *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		itemType: separator,
	}
	return result
}

func NewMenuItemCheckbox(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		checked:  checked,
		itemType: checkbox,
	}
	addToMenuItemMap(result)
	return result
}

func NewMenuItemRadio(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		checked:  checked,
		itemType: radio,
	}
	addToMenuItemMap(result)
	return result
}

func NewSubMenuItem(label string) *MenuItem {
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

func NewRole(role Role) *MenuItem {
	var result *MenuItem
	switch role {
	case AppMenu:
		result = NewAppMenu()
	case EditMenu:
		result = NewEditMenu()
	case FileMenu:
		result = NewFileMenu()
	case ViewMenu:
		result = NewViewMenu()
	case ServicesMenu:
		return NewServicesMenu()
	case SpeechMenu:
		result = NewSpeechMenu()
	case WindowMenu:
		result = NewWindowMenu()
	case HelpMenu:
		result = NewHelpMenu()
	case Hide:
		result = NewHideMenuItem()
	case Front:
		result = NewFrontMenuItem()
	case HideOthers:
		result = NewHideOthersMenuItem()
	case UnHide:
		result = NewUnhideMenuItem()
	case Undo:
		result = NewUndoMenuItem()
	case Redo:
		result = NewRedoMenuItem()
	case Cut:
		result = NewCutMenuItem()
	case Copy:
		result = NewCopyMenuItem()
	case Paste:
		result = NewPasteMenuItem()
	case PasteAndMatchStyle:
		result = NewPasteAndMatchStyleMenuItem()
	case SelectAll:
		result = NewSelectAllMenuItem()
	case Delete:
		result = NewDeleteMenuItem()
	case Quit:
		result = NewQuitMenuItem()
	case CloseWindow:
		result = NewCloseMenuItem()
	case About:
		result = NewAboutMenuItem()
	case Reload:
		result = NewReloadMenuItem()
	case ForceReload:
		result = NewForceReloadMenuItem()
	case ToggleFullscreen:
		result = NewToggleFullscreenMenuItem()
	case OpenDevTools:
		result = NewOpenDevToolsMenuItem()
	case ResetZoom:
		result = NewZoomResetMenuItem()
	case ZoomIn:
		result = NewZoomInMenuItem()
	case ZoomOut:
		result = NewZoomOutMenuItem()
	case Minimise:
		result = NewMinimiseMenuItem()
	case Zoom:
		result = NewZoomMenuItem()
	case FullScreen:
		result = NewFullScreenMenuItem()
	case Print:
		result = NewPrintMenuItem()
	case PageLayout:
		result = NewPageLayoutMenuItem()
	case NoRole:
	case ShowAll:
		result = NewShowAllMenuItem()
	case BringAllToFront:
		result = NewBringAllToFrontMenuItem()
	case NewFile:
		result = NewNewFileMenuItem()
	case Open:
		result = NewOpenMenuItem()
	case Save:
		result = NewSaveMenuItem()
	case SaveAs:
		result = NewSaveAsMenuItem()
	case StartSpeaking:
		result = NewStartSpeakingMenuItem()
	case StopSpeaking:
		result = NewStopSpeakingMenuItem()
	case Revert:
		result = NewRevertMenuItem()
	case Find:
		result = NewFindMenuItem()
	case FindAndReplace:
		result = NewFindAndReplaceMenuItem()
	case FindNext:
		result = NewFindNextMenuItem()
	case FindPrevious:
		result = NewFindPreviousMenuItem()
	case Help:
		result = NewHelpMenuItem()

	default:
		globalApplication.error(fmt.Sprintf("No support for role: %v", role))
		os.Exit(1)
	}

	if result == nil {
		return nil
	}

	result.role = role
	return result
}

func NewServicesMenu() *MenuItem {
	serviceMenu := NewSubMenuItem("Services")
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
		go func() {
			defer handlePanic()
			m.callback(ctx)
		}()
	}
}

func (m *MenuItem) SetAccelerator(shortcut string) *MenuItem {
	accelerator, err := parseAccelerator(shortcut)
	if err != nil {
		globalApplication.error("invalid accelerator. %v", err.Error())
		return m
	}
	m.accelerator = accelerator
	if m.impl != nil {
		m.impl.setAccelerator(accelerator)
	}
	return m
}

func (m *MenuItem) GetAccelerator() string {
	if m.accelerator == nil {
		return ""
	}
	return m.accelerator.String()
}

func (m *MenuItem) RemoveAccelerator() {
	m.accelerator = nil
}

func (m *MenuItem) SetTooltip(s string) *MenuItem {
	m.tooltip = s
	if m.impl != nil {
		m.impl.setTooltip(s)
	}
	return m
}

func (m *MenuItem) SetRole(role Role) *MenuItem {
	m.role = role
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

// GetSubmenu returns the submenu of the MenuItem.
// If the MenuItem is not a submenu, it returns nil.
func (m *MenuItem) GetSubmenu() *Menu {
	return m.submenu
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

// Clone returns a deep copy of the MenuItem
func (m *MenuItem) Clone() *MenuItem {
	result := &MenuItem{
		id:       m.id,
		label:    m.label,
		tooltip:  m.tooltip,
		disabled: m.disabled,
		checked:  m.checked,
		hidden:   m.hidden,
		bitmap:   m.bitmap,
		callback: m.callback,
		itemType: m.itemType,
		role:     m.role,
	}
	if m.submenu != nil {
		result.submenu = m.submenu.Clone()
	}
	if m.accelerator != nil {
		result.accelerator = m.accelerator.clone()
	}
	if m.contextMenuData != nil {
		result.contextMenuData = m.contextMenuData.clone()
	}
	return result
}
