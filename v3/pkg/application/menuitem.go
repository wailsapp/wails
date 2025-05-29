package application

import (
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

func removeMenuItemByID(id uint) {
	menuItemMapLock.Lock()
	defer menuItemMapLock.Unlock()
	delete(menuItemMap, id)
}

type menuItemImpl interface {
	setTooltip(s string)
	setLabel(s string)
	setDisabled(disabled bool)
	setChecked(checked bool)
	setAccelerator(accelerator *accelerator)
	setHidden(hidden bool)
	setBitmap(bitmap []byte)
	destroy()
}

type MenuItem struct {
	id              uint
	elementID       string // New ID for MenuElementInterface
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

	// Parent reference
	parentMenu *Menu
}

func NewMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		label:     label,
		itemType:  text,
	}
	addToMenuItemMap(result)
	registerElement(result)
	return result
}

func NewMenuItemSeparator() *MenuItem {
	result := &MenuItem{
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		itemType:  separator,
	}
	registerElement(result)
	return result
}

func NewMenuItemCheckbox(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		label:     label,
		checked:   checked,
		itemType:  checkbox,
	}
	addToMenuItemMap(result)
	registerElement(result)
	return result
}

func NewMenuItemRadio(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		label:     label,
		checked:   checked,
		itemType:  radio,
	}
	addToMenuItemMap(result)
	registerElement(result)
	return result
}

func NewSubMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		label:     label,
		itemType:  submenu,
		submenu: &Menu{
			label: label,
		},
	}
	addToMenuItemMap(result)
	registerElement(result)
	return result
}

// MenuItem implementation of MenuElementInterface

func (m *MenuItem) ID() string {
	return m.elementID
}

func (m *MenuItem) Label() string {
	return m.label
}

func (m *MenuItem) SetLabel(s string) MenuElementInterface {
	m.label = s
	if m.impl != nil {
		m.impl.setLabel(s)
	}
	// If this is a submenu, update the submenu's label as well
	if m.itemType == submenu && m.submenu != nil {
		m.submenu.label = s
	}
	return m
}

func (m *MenuItem) SetEnabled(enabled bool) MenuElementInterface {
	m.disabled = !enabled
	if m.impl != nil {
		m.impl.setDisabled(m.disabled)
	}
	return m
}

func (m *MenuItem) Enabled() bool {
	return !m.disabled
}

func (m *MenuItem) SetHidden(hidden bool) MenuElementInterface {
	m.hidden = hidden
	if m.impl != nil {
		m.impl.setHidden(m.hidden)
	}
	return m
}

func (m *MenuItem) Hidden() bool {
	return m.hidden
}

func (m *MenuItem) OnClick(f func(*Context)) MenuElementInterface {
	m.callback = f
	return m
}

func (m *MenuItem) SetAccelerator(shortcut string) MenuElementInterface {
	accelerator, err := parseAccelerator(shortcut)
	if err != nil {
		globalApplication.error("invalid accelerator: %w", err)
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

func (m *MenuItem) IsSubmenu() bool {
	return m.itemType == submenu
}

func (m *MenuItem) GetSubmenu() *Menu {
	return m.submenu
}

func (m *MenuItem) Update() MenuElementInterface {
	if m.impl != nil {
		// Update platform-specific implementation
		m.impl.setLabel(m.label)
		m.impl.setDisabled(m.disabled)
		m.impl.setHidden(m.hidden)
		m.impl.setChecked(m.checked)
		if m.accelerator != nil {
			m.impl.setAccelerator(m.accelerator)
		}
	}
	// Update submenu if this is a submenu
	if m.itemType == submenu && m.submenu != nil {
		m.submenu.Update()
	}
	return m
}

// Legacy MenuItem methods that still work

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

func (m *MenuItem) Checked() bool {
	return m.checked
}

func (m *MenuItem) IsSeparator() bool {
	return m.itemType == separator
}

func (m *MenuItem) IsCheckbox() bool {
	return m.itemType == checkbox
}

func (m *MenuItem) IsRadio() bool {
	return m.itemType == radio
}

func (m *MenuItem) Tooltip() string {
	return m.tooltip
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
		id:        uint(atomic.AddUintptr(&menuItemID, 1)),
		elementID: generateElementID(),
		label:     m.label,
		tooltip:   m.tooltip,
		disabled:  m.disabled,
		checked:   m.checked,
		hidden:    m.hidden,
		bitmap:    m.bitmap,
		callback:  m.callback,
		itemType:  m.itemType,
		role:      m.role,
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
	addToMenuItemMap(result)
	registerElement(result)
	return result
}

func (m *MenuItem) Destroy() {
	removeMenuItemByID(m.id)
	unregisterElement(m.elementID)

	// Clean up resources
	if m.impl != nil {
		m.impl.destroy()
	}
	if m.submenu != nil {
		m.submenu.Destroy()
		m.submenu = nil
	}

	if m.contextMenuData != nil {
		m.contextMenuData = nil
	}

	if m.accelerator != nil {
		m.accelerator = nil
	}

	m.callback = nil
	m.radioGroupMembers = nil
}

// NewRole creates menu items based on predefined roles
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
		globalApplication.error("no support for role: %v", role)
	}

	if result != nil {
		result.role = role
	}

	return result
}

func NewServicesMenu() *MenuItem {
	serviceMenu := NewSubMenuItem("Services")
	serviceMenu.role = ServicesMenu
	return serviceMenu
}

// AsMenuItem is a helper function to cast MenuElementInterface back to *MenuItem
// Returns nil if the element is not a MenuItem
func AsMenuItem(element MenuElementInterface) *MenuItem {
	if menuItem, ok := element.(*MenuItem); ok {
		return menuItem
	}
	return nil
}

// AsMenu is a helper function to cast MenuElementInterface back to *Menu
// Returns nil if the element is not a Menu
func AsMenu(element MenuElementInterface) *Menu {
	if menu, ok := element.(*Menu); ok {
		return menu
	}
	return nil
}

// Convenience methods for MenuItem to allow fluent chaining while returning *MenuItem

// SetLabelItem sets the label and returns *MenuItem for chaining
func (m *MenuItem) SetLabelItem(label string) *MenuItem {
	m.SetLabel(label)
	return m
}

// SetEnabledItem sets the enabled state and returns *MenuItem for chaining
func (m *MenuItem) SetEnabledItem(enabled bool) *MenuItem {
	m.SetEnabled(enabled)
	return m
}

// SetHiddenItem sets the hidden state and returns *MenuItem for chaining
func (m *MenuItem) SetHiddenItem(hidden bool) *MenuItem {
	m.SetHidden(hidden)
	return m
}

// OnClickItem sets the click callback and returns *MenuItem for chaining
func (m *MenuItem) OnClickItem(callback func(*Context)) *MenuItem {
	m.OnClick(callback)
	return m
}

// SetAcceleratorItem sets the accelerator and returns *MenuItem for chaining
func (m *MenuItem) SetAcceleratorItem(shortcut string) *MenuItem {
	m.SetAccelerator(shortcut)
	return m
}

// SetRoleItem sets the role and returns *MenuItem for chaining
func (m *MenuItem) SetRoleItem(role Role) *MenuItem {
	m.SetRole(role)
	return m
}
