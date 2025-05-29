package application

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type menuImpl interface {
	update()
}

// MenuElementInterface defines the common interface for both Menu and MenuItem
type MenuElementInterface interface {
	// Identity methods
	ID() string
	Label() string
	SetLabel(string) MenuElementInterface

	// State methods
	SetEnabled(bool) MenuElementInterface
	Enabled() bool
	SetHidden(bool) MenuElementInterface
	Hidden() bool

	// Behavior methods
	OnClick(func(*Context)) MenuElementInterface
	SetAccelerator(string) MenuElementInterface
	GetAccelerator() string

	// Submenu capability
	IsSubmenu() bool
	GetSubmenu() *Menu

	// For updating the menu system
	Update() MenuElementInterface
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

// Menu represents a menu that can also act as a MenuItem
type Menu struct {
	id    string
	label string
	items []MenuElementInterface

	// MenuItem properties when this menu acts as a menu item
	enabled     bool
	hidden      bool
	accelerator *accelerator
	callback    func(*Context)
	role        Role

	impl        menuImpl
	parentMenu  *Menu
	contextData *ContextMenuData

	// Track which windows this menu is attached to
	attachedWindows     map[uint]bool
	attachedWindowsLock sync.RWMutex

	// ID management
	itemRegistry  map[string]MenuElementInterface
	registryMutex sync.RWMutex
}

// Global registry for quick lookups
var (
	globalMenuRegistry      = make(map[string]MenuElementInterface)
	globalMenuRegistryMutex sync.RWMutex
	globalElementIDCounter  uint64
)

// Helper functions for ID generation
func generateElementID() string {
	return fmt.Sprintf("element_%d", atomic.AddUint64(&globalElementIDCounter, 1))
}

func registerElement(element MenuElementInterface) {
	globalMenuRegistryMutex.Lock()
	defer globalMenuRegistryMutex.Unlock()
	globalMenuRegistry[element.ID()] = element
}

func unregisterElement(id string) {
	globalMenuRegistryMutex.Lock()
	defer globalMenuRegistryMutex.Unlock()
	delete(globalMenuRegistry, id)
}

func GetElementByID(id string) MenuElementInterface {
	globalMenuRegistryMutex.RLock()
	defer globalMenuRegistryMutex.RUnlock()
	return globalMenuRegistry[id]
}

func GetElementByLabel(label string) MenuElementInterface {
	globalMenuRegistryMutex.RLock()
	defer globalMenuRegistryMutex.RUnlock()
	for _, element := range globalMenuRegistry {
		if element.Label() == label {
			return element
		}
	}
	return nil
}

func NewMenu() *Menu {
	menu := &Menu{
		id:              generateElementID(),
		enabled:         true,
		attachedWindows: make(map[uint]bool),
		itemRegistry:    make(map[string]MenuElementInterface),
	}
	registerElement(menu)
	return menu
}

// Menu implementation of MenuElementInterface

func (m *Menu) ID() string {
	return m.id
}

func (m *Menu) Label() string {
	return m.label
}

func (m *Menu) SetLabel(label string) MenuElementInterface {
	m.label = label
	return m
}

func (m *Menu) SetEnabled(enabled bool) MenuElementInterface {
	m.enabled = enabled
	return m
}

func (m *Menu) Enabled() bool {
	return m.enabled
}

func (m *Menu) SetHidden(hidden bool) MenuElementInterface {
	m.hidden = hidden
	return m
}

func (m *Menu) Hidden() bool {
	return m.hidden
}

func (m *Menu) OnClick(callback func(*Context)) MenuElementInterface {
	m.callback = callback
	return m
}

func (m *Menu) SetAccelerator(shortcut string) MenuElementInterface {
	accelerator, err := parseAccelerator(shortcut)
	if err != nil {
		globalApplication.error("invalid accelerator: %w", err)
		return m
	}
	m.accelerator = accelerator
	return m
}

func (m *Menu) GetAccelerator() string {
	if m.accelerator == nil {
		return ""
	}
	return m.accelerator.String()
}

func (m *Menu) IsSubmenu() bool {
	return true // Menus are always submenus when used as menu items
}

func (m *Menu) GetSubmenu() *Menu {
	return m // A menu returns itself as the submenu
}

func (m *Menu) Add(label string) *MenuItem {
	result := NewMenuItem(label)
	m.AddElement(result)
	return result
}

// AddElement adds an existing menu element (item or submenu)
func (m *Menu) AddElement(element MenuElementInterface) *Menu {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	m.items = append(m.items, element)
	m.itemRegistry[element.ID()] = element

	// Set parent reference
	if menu, ok := element.(*Menu); ok {
		menu.parentMenu = m
	}

	return m
}

func (m *Menu) AddSeparator() {
	result := NewMenuItemSeparator()
	m.AddElement(result)
}

func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	result := NewMenuItemCheckbox(label, enabled)
	m.AddElement(result)
	return result
}

func (m *Menu) AddRadio(label string, enabled bool) *MenuItem {
	result := NewMenuItemRadio(label, enabled)
	m.AddElement(result)
	return result
}

// AddSubmenu creates and adds a new submenu, returning the submenu for chaining
func (m *Menu) AddSubmenu(s string) *Menu {
	submenu := NewMenu()
	submenu.SetLabel(s)
	m.AddElement(submenu)
	return submenu
}

func (m *Menu) AddRole(role Role) *Menu {
	result := NewRole(role)
	if result != nil {
		m.AddElement(result)
	}
	return m
}

// Update updates the menu and propagates the changes to the application and all windows
// that have the menu attached.
func (m *Menu) Update() MenuElementInterface {
	// Update all child items
	for _, item := range m.items {
		item.Update()
	}

	// Process radio groups and update the menu
	m.processRadioGroups()
	if m.impl == nil {
		m.impl = newMenuImpl(m.convertToOldMenu())
	}
	m.impl.update()

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
		globalApplication.SetMenu(m.convertToOldMenu())
	}

	// Update all windows that use this menu
	globalApplication.windowsLock.RLock()
	defer globalApplication.windowsLock.RUnlock()

	for id, window := range globalApplication.windows {
		if attachedWindowsCopy[id] {
			window.SetMenu(m.convertToOldMenu())
		}
	}

	return m
}

// FindByID finds an element by its ID
func (m *Menu) FindByID(id string) MenuElementInterface {
	m.registryMutex.RLock()
	defer m.registryMutex.RUnlock()

	if element, exists := m.itemRegistry[id]; exists {
		return element
	}

	// Search recursively in submenus
	for _, item := range m.items {
		if submenu := item.GetSubmenu(); submenu != nil && submenu != item {
			if found := submenu.FindByID(id); found != nil {
				return found
			}
		}
	}

	return nil
}

// FindByLabel finds an element by its label - updated to use MenuElementInterface
func (m *Menu) FindByLabel(label string) MenuElementInterface {
	m.registryMutex.RLock()
	defer m.registryMutex.RUnlock()

	for _, item := range m.items {
		if item.Label() == label {
			return item
		}
		// Search recursively in submenus
		if submenu := item.GetSubmenu(); submenu != nil && submenu != item {
			if found := submenu.FindByLabel(label); found != nil {
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
		if menuItem, ok := item.(*MenuItem); ok {
			if menuItem.role == role {
				return menuItem
			}
		}
		if submenu := item.GetSubmenu(); submenu != nil && submenu != item {
			found := submenu.FindByRole(role)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// RemoveByID removes an element by ID
func (m *Menu) RemoveByID(id string) bool {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	for i, item := range m.items {
		if item.ID() == id {
			// Remove from slice
			m.items = append(m.items[:i], m.items[i+1:]...)
			// Remove from registry
			delete(m.itemRegistry, id)
			unregisterElement(id)
			return true
		}
	}
	return false
}

// RemoveByLabel removes an element by label
func (m *Menu) RemoveByLabel(label string) bool {
	if element := m.FindByLabel(label); element != nil {
		return m.RemoveByID(element.ID())
	}
	return false
}

func (m *Menu) RemoveMenuItem(target *MenuItem) {
	m.RemoveByID(target.ID())
}

// InsertAtIndex inserts an element at a specific index
func (m *Menu) InsertAtIndex(index int, element MenuElementInterface) *Menu {
	m.registryMutex.Lock()
	defer m.registryMutex.Unlock()

	if index < 0 {
		index = 0
	}
	if index > len(m.items) {
		index = len(m.items)
	}

	if index == len(m.items) {
		m.items = append(m.items, element)
	} else {
		m.items = append(m.items[:index], append([]MenuElementInterface{element}, m.items[index:]...)...)
	}

	m.itemRegistry[element.ID()] = element

	// Set parent reference
	if item, ok := element.(*MenuItem); ok {
		item.parentMenu = m
	} else if menu, ok := element.(*Menu); ok {
		menu.parentMenu = m
	}

	return m
}

// InsertBeforeID inserts an element before the element with the given ID
func (m *Menu) InsertBeforeID(beforeID string, element MenuElementInterface) bool {
	for i, item := range m.items {
		if item.ID() == beforeID {
			m.InsertAtIndex(i, element)
			return true
		}
	}
	return false
}

// InsertAfterID inserts an element after the element with the given ID
func (m *Menu) InsertAfterID(afterID string, element MenuElementInterface) bool {
	for i, item := range m.items {
		if item.ID() == afterID {
			m.InsertAtIndex(i+1, element)
			return true
		}
	}
	return false
}

// InsertBeforeLabel inserts an element before the element with the given label
func (m *Menu) InsertBeforeLabel(beforeLabel string, element MenuElementInterface) bool {
	if target := m.FindByLabel(beforeLabel); target != nil {
		return m.InsertBeforeID(target.ID(), element)
	}
	return false
}

// InsertAfterLabel inserts an element after the element with the given label
func (m *Menu) InsertAfterLabel(afterLabel string, element MenuElementInterface) bool {
	if target := m.FindByLabel(afterLabel); target != nil {
		return m.InsertAfterID(target.ID(), element)
	}
	return false
}

// Legacy methods that still work

// ItemAt returns the menu item at the given index, or nil if the index is out of bounds.
func (m *Menu) ItemAt(index int) *MenuItem {
	if index < 0 || index >= len(m.items) {
		return nil
	}
	if item, ok := m.items[index].(*MenuItem); ok {
		return item
	}
	return nil
}

// InsertAt inserts a new menu item with the given label at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertAt(index int, label string) *MenuItem {
	result := NewMenuItem(label)
	m.InsertAtIndex(index, result)
	return result
}

// InsertItemAt inserts an existing menu item at the specified index.
// Returns the menu for chaining.
func (m *Menu) InsertItemAt(index int, item *MenuItem) *Menu {
	m.InsertAtIndex(index, item)
	return m
}

// InsertSeparatorAt inserts a separator at the specified index.
// Returns the menu for chaining.
func (m *Menu) InsertSeparatorAt(index int) *Menu {
	result := NewMenuItemSeparator()
	m.InsertAtIndex(index, result)
	return m
}

// InsertCheckboxAt inserts a checkbox menu item at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertCheckboxAt(index int, label string, enabled bool) *MenuItem {
	result := NewMenuItemCheckbox(label, enabled)
	m.InsertAtIndex(index, result)
	return result
}

// InsertRadioAt inserts a radio menu item at the specified index.
// Returns the newly created menu item.
func (m *Menu) InsertRadioAt(index int, label string, enabled bool) *MenuItem {
	result := NewMenuItemRadio(label, enabled)
	m.InsertAtIndex(index, result)
	return result
}

// InsertSubmenuAt inserts a submenu at the specified index.
// Returns the newly created submenu.
func (m *Menu) InsertSubmenuAt(index int, label string) *Menu {
	submenu := NewMenu()
	submenu.SetLabel(label)
	m.InsertAtIndex(index, submenu)
	return submenu
}

// Clear all menu items
func (m *Menu) Clear() {
	for _, item := range m.items {
		if menuItem, ok := item.(*MenuItem); ok {
			removeMenuItemByID(menuItem.id)
		}
		unregisterElement(item.ID())
	}
	m.items = nil
}

func (m *Menu) Destroy() {
	for _, item := range m.items {
		if menuItem, ok := item.(*MenuItem); ok {
			menuItem.Destroy()
		} else if submenu, ok := item.(*Menu); ok {
			submenu.Destroy()
		}
	}
	m.items = nil
}

// Clone recursively clones the menu and all its submenus.
func (m *Menu) Clone() *Menu {
	result := &Menu{
		id:              generateElementID(),
		label:           m.label,
		attachedWindows: make(map[uint]bool),
		itemRegistry:    make(map[string]MenuElementInterface),
	}
	registerElement(result)
	for _, item := range m.items {
		if menuItem, ok := item.(*MenuItem); ok {
			clonedItem := menuItem.Clone()
			result.items = append(result.items, clonedItem)
		} else if submenu, ok := item.(*Menu); ok {
			clonedSubmenu := submenu.Clone()
			result.items = append(result.items, clonedSubmenu)
		}
	}
	return result
}

// Append menu to an existing menu
func (m *Menu) Append(in *Menu) {
	if in == nil {
		return
	}
	for _, item := range in.items {
		m.AddElement(item)
	}
}

// AppendItem appends a menu item to the menu
func (m *Menu) AppendItem(item *MenuItem) *Menu {
	if item == nil {
		return m
	}
	m.AddElement(item)
	return m
}

// Remove removes a menu item from the menu
func (m *Menu) Remove(index int) *Menu {
	if index < 0 || index >= len(m.items) {
		return m
	}
	itemToRemove := m.items[index]
	m.RemoveByID(itemToRemove.ID())
	return m
}

func (m *Menu) processRadioGroups() {
	var radioGroup []MenuElementInterface

	closeOutRadioGroups := func() {
		if len(radioGroup) > 0 {
			for _, item := range radioGroup {
				if radioItem, ok := item.(*MenuItem); ok {
					// Convert MenuElementInterface slice to MenuItem slice for compatibility
					radioMembers := make([]*MenuItem, 0, len(radioGroup))
					for _, member := range radioGroup {
						if menuItem, ok := member.(*MenuItem); ok {
							radioMembers = append(radioMembers, menuItem)
						}
					}
					radioItem.radioGroupMembers = radioMembers
				}
			}
			radioGroup = []MenuElementInterface{}
		}
	}

	for _, item := range m.items {
		if menuItem, ok := item.(*MenuItem); ok {
			if menuItem.itemType != radio {
				closeOutRadioGroups()
			}
			if menuItem.itemType == radio {
				radioGroup = append(radioGroup, item)
			}
		} else {
			closeOutRadioGroups()
		}

		// Process submenus recursively
		if submenu := item.GetSubmenu(); submenu != nil && submenu != item {
			submenu.processRadioGroups()
		}
	}
	closeOutRadioGroups()
}

func (m *Menu) setContextData(data *ContextMenuData) {
	for _, item := range m.items {
		if menuItem, ok := item.(*MenuItem); ok {
			menuItem.setContextData(data)
		} else if submenu, ok := item.(*Menu); ok {
			submenu.setContextData(data)
		}
	}
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

// updateParentMenuItemLabel finds and updates the parent MenuItem that contains this submenu
func (m *Menu) updateParentMenuItemLabel(label string) {
	// This is no longer needed as menus can update their own labels directly
	m.SetLabel(label)
}

// UpdateAll updates the menu and propagates the changes to the application and all windows
// This is a convenience method that handles all the necessary update steps after modifying a menu
//
// Deprecated: Use Update() instead as it now always propagates changes.
func (m *Menu) UpdateAll() *Menu {
	m.Update()
	return m
}

// Prepend menu before an existing menu
func (m *Menu) Prepend(in *Menu) {
	if in == nil {
		return
	}
	// Insert all items from 'in' at the beginning
	for i, item := range in.items {
		m.InsertAtIndex(i, item)
	}
}

// Conversion methods to maintain compatibility with existing platform implementations
func (m *Menu) convertToOldMenu() *Menu {
	// For now, return self since we're replacing the old system
	// This method may need platform-specific implementations
	return m
}

func (a *App) NewMenu() *Menu {
	return NewMenu()
}

func NewMenuFromItems(item *MenuItem, items ...*MenuItem) *Menu {
	result := NewMenu()
	result.AddElement(item)
	for _, itm := range items {
		result.AddElement(itm)
	}
	return result
}

func NewSubmenu(s string, items *Menu) *MenuItem {
	// This creates a MenuItem that wraps a Menu - for backward compatibility
	result := NewMenuItem(s)
	result.itemType = submenu
	result.submenu = items
	return result
}
