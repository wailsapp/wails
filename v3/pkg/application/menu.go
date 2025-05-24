package application

import "sync"

type menuImpl interface {
	update()
}

type ContextMenu struct {
	*Menu
	name string
}

func NewContextMenu(name string) *ContextMenu {
	result := &ContextMenu{
		Menu: NewMenu(),
		name: name,
	}
	result.Update()
	return result
}

func (m *ContextMenu) Update() {
	m.Menu.Update()
	globalApplication.registerContextMenu(m)
}

func (m *ContextMenu) Destroy() {
	globalApplication.unregisterContextMenu(m.name)
}

type Menu struct {
	items []*MenuItem
	label string

	impl menuImpl

	// Track which windows this menu is attached to
	attachedWindows     map[uint]bool
	attachedWindowsLock sync.RWMutex
}

func NewMenu() *Menu {
	return &Menu{
		attachedWindows: make(map[uint]bool),
	}
}

func (m *Menu) Add(label string) *MenuItem {
	result := NewMenuItem(label)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddSeparator() {
	result := NewMenuItemSeparator()
	m.items = append(m.items, result)
}

func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	result := NewMenuItemCheckbox(label, enabled)
	m.items = append(m.items, result)
	return result
}

func (m *Menu) AddRadio(label string, enabled bool) *MenuItem {
	result := NewMenuItemRadio(label, enabled)
	m.items = append(m.items, result)
	return result
}

// Update updates the menu and optionally propagates the changes to the application and all windows
// that have the menu attached.
//
// If propagate is true, the changes will be propagated to the application and all windows
// that have the menu attached. This is equivalent to calling UpdateAll().
func (m *Menu) Update(propagate ...bool) *Menu {
	// Process radio groups and update the menu
	m.processRadioGroups()
	if m.impl == nil {
		m.impl = newMenuImpl(m)
	}
	m.impl.update()

	// Check if we should propagate the changes
	shouldPropagate := false
	if len(propagate) > 0 {
		shouldPropagate = propagate[0]
	}

	if shouldPropagate {
		// If globalApplication is not initialized yet, just return
		if globalApplication == nil {
			return m
		}

		// Get a copy of the attached windows
		m.attachedWindowsLock.RLock()
		attachedWindowsCopy := make(map[uint]bool)
		for id, attached := range m.attachedWindows {
			attachedWindowsCopy[id] = attached
		}
		m.attachedWindowsLock.RUnlock()

		// Check if this menu is the application menu (ID 0)
		if attachedWindowsCopy[0] {
			// Update the application menu
			globalApplication.SetMenu(m)
		}

		// Update all windows that use this menu
		globalApplication.windowsLock.RLock()
		defer globalApplication.windowsLock.RUnlock()

		for id, window := range globalApplication.windows {
			if attachedWindowsCopy[id] {
				window.SetMenu(m)
			}
		}
	}

	return m
}

// Clear all menu items
func (m *Menu) Clear() {
	for _, item := range m.items {
		removeMenuItemByID(item.id)
	}
	m.items = nil
}

func (m *Menu) Destroy() {
	for _, item := range m.items {
		item.Destroy()
	}
	m.items = nil
}

func (m *Menu) AddSubmenu(s string) *Menu {
	result := NewSubMenuItem(s)
	m.items = append(m.items, result)
	return result.submenu
}

func (m *Menu) AddRole(role Role) *Menu {
	result := NewRole(role)
	if result != nil {
		m.items = append(m.items, result)
	}
	return m
}

func (m *Menu) processRadioGroups() {
	var radioGroup []*MenuItem

	closeOutRadioGroups := func() {
		if len(radioGroup) > 0 {
			for _, item := range radioGroup {
				item.radioGroupMembers = radioGroup
			}
			radioGroup = []*MenuItem{}
		}
	}

	for _, item := range m.items {
		if item.itemType != radio {
			closeOutRadioGroups()
		}
		if item.itemType == submenu {
			item.submenu.processRadioGroups()
			continue
		}
		if item.itemType == radio {
			radioGroup = append(radioGroup, item)
		}
	}
	closeOutRadioGroups()
}

func (m *Menu) SetLabel(label string) {
	m.label = label
	// If this menu is a submenu, update the parent MenuItem's label as well
	m.updateParentMenuItemLabel(label)
}

// updateParentMenuItemLabel finds and updates the parent MenuItem that contains this submenu
func (m *Menu) updateParentMenuItemLabel(label string) {
	// Find all menu items that have this menu as their submenu
	menuItemMapLock.Lock()
	defer menuItemMapLock.Unlock()

	for _, item := range menuItemMap {
		if item.submenu == m {
			item.SetLabel(label)
			break // There should only be one parent MenuItem per submenu
		}
	}
}

func (m *Menu) setContextData(data *ContextMenuData) {
	for _, item := range m.items {
		item.setContextData(data)
	}
}

// FindByLabel recursively searches for a menu item with the given label
// and returns the first match, or nil if not found.
func (m *Menu) FindByLabel(label string) *MenuItem {
	for _, item := range m.items {
		if item.label == label {
			return item
		}
		if item.submenu != nil {
			found := item.submenu.FindByLabel(label)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// FindByRole recursively searches for a menu item with the given role
// and returns the first match, or nil if not found.
func (m *Menu) FindByRole(role Role) *MenuItem {
	for _, item := range m.items {
		if item.role == role {
			return item
		}
		if item.submenu != nil {
			found := item.submenu.FindByRole(role)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

func (m *Menu) RemoveMenuItem(target *MenuItem) {
	for i, item := range m.items {
		if item == target {
			// Remove the item from the slice
			m.items = append(m.items[:i], m.items[i+1:]...)
			break
		}
		if item.submenu != nil {
			item.submenu.RemoveMenuItem(target)
		}
	}
}

// ItemAt returns the menu item at the given index, or nil if the index is out of bounds.
func (m *Menu) ItemAt(index int) *MenuItem {
	if index < 0 || index >= len(m.items) {
		return nil
	}
	return m.items[index]
}

// InsertAt inserts a new menu item with the given label at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertAt(index int, label string) *MenuItem {
	result := NewMenuItem(label)
	m.InsertItemAt(index, result)
	return result
}

// InsertItemAt inserts an existing menu item at the specified index.
// Returns the menu for chaining.
func (m *Menu) InsertItemAt(index int, item *MenuItem) *Menu {
	if index < 0 {
		index = 0
	}
	if index > len(m.items) {
		index = len(m.items)
	}

	if index == len(m.items) {
		m.items = append(m.items, item)
	} else {
		// Create a new slice with the item inserted at the specified index
		m.items = append(m.items[:index], append([]*MenuItem{item}, m.items[index:]...)...)
	}

	return m
}

// InsertSeparatorAt inserts a separator at the specified index.
// Returns the menu for chaining.
func (m *Menu) InsertSeparatorAt(index int) *Menu {
	result := NewMenuItemSeparator()
	return m.InsertItemAt(index, result)
}

// InsertCheckboxAt inserts a checkbox menu item at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertCheckboxAt(index int, label string, enabled bool) *MenuItem {
	result := NewMenuItemCheckbox(label, enabled)
	m.InsertItemAt(index, result)
	return result
}

// InsertRadioAt inserts a radio menu item at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertRadioAt(index int, label string, enabled bool) *MenuItem {
	result := NewMenuItemRadio(label, enabled)
	m.InsertItemAt(index, result)
	return result
}

// InsertSubmenuAt inserts a submenu at the specified index.
// Returns the newly created submenu.
func (m *Menu) InsertSubmenuAt(index int, label string) *Menu {
	result := NewSubMenuItem(label)
	m.InsertItemAt(index, result)
	return result.submenu
}

// Clone recursively clones the menu and all its submenus.
func (m *Menu) Clone() *Menu {
	result := &Menu{
		label:           m.label,
		attachedWindows: make(map[uint]bool),
	}
	for _, item := range m.items {
		result.items = append(result.items, item.Clone())
	}
	return result
}

// Append menu to an existing menu
func (m *Menu) Append(in *Menu) {
	if in == nil {
		return
	}
	m.items = append(m.items, in.items...)
}

// AppendItem appends a menu item to the menu
func (m *Menu) AppendItem(item *MenuItem) *Menu {
	if item == nil {
		return m
	}
	m.items = append(m.items, item)
	return m
}

// Remove removes a menu item from the menu
func (m *Menu) Remove(index int) *Menu {
	if index < 0 || index >= len(m.items) {
		return m
	}
	m.items = append(m.items[:index], m.items[index+1:]...)
	return m
}

// registerWindow registers a window as using this menu
func (m *Menu) registerWindow(windowID uint) {
	m.attachedWindowsLock.Lock()
	defer m.attachedWindowsLock.Unlock()
	m.attachedWindows[windowID] = true
}

// unregisterWindow unregisters a window from using this menu
func (m *Menu) unregisterWindow(windowID uint) {
	m.attachedWindowsLock.Lock()
	defer m.attachedWindowsLock.Unlock()
	delete(m.attachedWindows, windowID)
}

// UpdateAll updates the menu and propagates the changes to the application and all windows
// This is a convenience method that handles all the necessary update steps after modifying a menu
//
// Deprecated: Use Update(true) instead.
func (m *Menu) UpdateAll() *Menu {
	return m.Update(true)
}

// Prepend menu before an existing menu
func (m *Menu) Prepend(in *Menu) {
	m.items = append(in.items, m.items...)
}

func (a *App) NewMenu() *Menu {
	return &Menu{
		attachedWindows: make(map[uint]bool),
	}
}

func NewMenuFromItems(item *MenuItem, items ...*MenuItem) *Menu {
	result := &Menu{
		items:           []*MenuItem{item},
		attachedWindows: make(map[uint]bool),
	}
	result.items = append(result.items, items...)
	return result
}

func NewSubmenu(s string, items *Menu) *MenuItem {
	result := NewSubMenuItem(s)
	result.submenu = items
	return result
}
