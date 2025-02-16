package application

import (
	"github.com/wailsapp/wails/v3/pkg/w32"
)

const (
	MenuItemMsgID = w32.WM_APP + 1024
)

type RadioGroupMember struct {
	ID       int
	MenuItem *MenuItem
}

type RadioGroup []*RadioGroupMember

func (r *RadioGroup) Add(id int, item *MenuItem) {
	*r = append(*r, &RadioGroupMember{
		ID:       id,
		MenuItem: item,
	})
}

func (r *RadioGroup) Bounds() (int, int) {
	p := *r
	return p[0].ID, p[len(p)-1].ID
}

func (r *RadioGroup) MenuID(item *MenuItem) int {
	for _, member := range *r {
		if member.MenuItem == item {
			return member.ID
		}
	}
	panic("RadioGroup.MenuID: item not found:")
}

type Win32Menu struct {
	isPopup       bool
	menu          w32.HMENU
	parentWindow  *windowsWebviewWindow
	parent        w32.HWND
	menuMapping   map[int]*MenuItem
	checkboxItems map[*MenuItem][]int
	radioGroups   map[*MenuItem][]*RadioGroup
	menuData      *Menu
	currentMenuID int
	onMenuClose   func()
	onMenuOpen    func()
}

func (p *Win32Menu) newMenu() w32.HMENU {
	if p.isPopup {
		return w32.NewPopupMenu()
	}
	return w32.CreateMenu()
}

func (p *Win32Menu) buildMenu(parentMenu w32.HMENU, inputMenu *Menu) {
	currentRadioGroup := RadioGroup{}
	for _, item := range inputMenu.items {
		if item.Hidden() {
			if item.accelerator != nil {
				if p.parentWindow != nil {
					// Remove the accelerator from the keybindings
					p.parentWindow.parent.removeMenuBinding(item.accelerator)
				} else {
					// Remove the global keybindings
					globalApplication.removeKeyBinding(item.accelerator.String())
				}
			}
			continue
		}
		p.currentMenuID++
		itemID := p.currentMenuID
		p.menuMapping[itemID] = item

		menuItemImpl := newMenuItemImpl(item, parentMenu, itemID)
		menuItemImpl.parent = inputMenu

		flags := uint32(w32.MF_STRING)
		if item.disabled {
			flags = flags | w32.MF_GRAYED
		}
		if item.checked {
			flags = flags | w32.MF_CHECKED
		}
		if item.IsSeparator() {
			flags = flags | w32.MF_SEPARATOR
		}

		if item.checked && item.IsRadio() {
			flags = flags | w32.MFT_RADIOCHECK
		}

		if item.IsCheckbox() {
			p.checkboxItems[item] = append(p.checkboxItems[item], itemID)
		}
		if item.IsRadio() {
			currentRadioGroup.Add(itemID, item)
		} else {
			if len(currentRadioGroup) > 0 {
				for _, radioMember := range currentRadioGroup {
					currentRadioGroup := currentRadioGroup
					p.radioGroups[radioMember.MenuItem] = append(p.radioGroups[radioMember.MenuItem], &currentRadioGroup)
				}
				currentRadioGroup = RadioGroup{}
			}
		}

		if item.submenu != nil {
			flags = flags | w32.MF_POPUP
			newSubmenu := p.newMenu()
			p.buildMenu(newSubmenu, item.submenu)
			itemID = int(newSubmenu)
			menuItemImpl.submenu = newSubmenu
		}

		var menuText = item.Label()
		if item.accelerator != nil {
			menuText = menuText + "\t" + item.accelerator.String()
			if item.callback != nil {
				if p.parentWindow != nil {
					p.parentWindow.parent.addMenuBinding(item.accelerator, item)
				} else {
					globalApplication.addKeyBinding(item.accelerator.String(), func(w *WebviewWindow) {
						item.handleClick()
					})
				}
			}
		}
		ok := w32.AppendMenu(parentMenu, flags, uintptr(itemID), w32.MustStringToUTF16Ptr(menuText))
		if !ok {
			globalApplication.fatal("error adding menu item '%s'", menuText)
		}
		if item.bitmap != nil {
			err := w32.SetMenuIcons(parentMenu, itemID, item.bitmap, nil)
			if err != nil {
				globalApplication.fatal("error setting menu icons: %w", err)
			}
		}

		item.impl = menuItemImpl
	}
	if len(currentRadioGroup) > 0 {
		for _, radioMember := range currentRadioGroup {
			currentRadioGroup := currentRadioGroup
			p.radioGroups[radioMember.MenuItem] = append(p.radioGroups[radioMember.MenuItem], &currentRadioGroup)
		}
		currentRadioGroup = RadioGroup{}
	}
}

func (p *Win32Menu) Update() {
	p.menu = p.newMenu()
	p.menuMapping = make(map[int]*MenuItem)
	p.currentMenuID = MenuItemMsgID
	p.buildMenu(p.menu, p.menuData)
	p.updateRadioGroups()
}

func NewPopupMenu(parent w32.HWND, inputMenu *Menu) *Win32Menu {
	result := &Win32Menu{
		isPopup:       true,
		parent:        parent,
		menuData:      inputMenu,
		checkboxItems: make(map[*MenuItem][]int),
		radioGroups:   make(map[*MenuItem][]*RadioGroup),
	}
	result.Update()
	return result
}
func NewApplicationMenu(parent *windowsWebviewWindow, inputMenu *Menu) *Win32Menu {
	result := &Win32Menu{
		parentWindow:  parent,
		parent:        parent.hwnd,
		menuData:      inputMenu,
		checkboxItems: make(map[*MenuItem][]int),
		radioGroups:   make(map[*MenuItem][]*RadioGroup),
	}
	result.Update()
	return result
}

func (p *Win32Menu) ShowAt(x int, y int) {

	w32.SetForegroundWindow(p.parent)

	if p.onMenuOpen != nil {
		p.onMenuOpen()
	}

	if !w32.TrackPopupMenuEx(p.menu, w32.TPM_LEFTALIGN, int32(x), int32(y-5), p.parent, nil) {
		globalApplication.fatal("TrackPopupMenu failed")
	}

	if p.onMenuClose != nil {
		p.onMenuClose()
	}

	if !w32.PostMessage(p.parent, w32.WM_NULL, 0, 0) {
		globalApplication.fatal("PostMessage failed")
	}

}

func (p *Win32Menu) ShowAtCursor() {
	x, y, ok := w32.GetCursorPos()
	if ok == false {
		globalApplication.fatal("GetCursorPos failed")
	}

	p.ShowAt(x, y)
}

func (p *Win32Menu) ProcessCommand(cmdMsgID int) bool {
	item := p.menuMapping[cmdMsgID]
	if item == nil {
		return false
	}
	if item.IsRadio() {
		if item.checked {
			return true
		}
		item.checked = true
		p.updateRadioGroup(item)
	}
	if item.callback != nil {
		item.handleClick()
	}
	return true
}

func (p *Win32Menu) Destroy() {
	w32.DestroyMenu(p.menu)
}

func (p *Win32Menu) UpdateMenuItem(item *MenuItem) {
	if item.IsCheckbox() {
		for _, itemID := range p.checkboxItems[item] {
			var checkState uint = w32.MF_UNCHECKED
			if item.checked {
				checkState = w32.MF_CHECKED
			}
			w32.CheckMenuItem(p.menu, uintptr(itemID), checkState)
		}
		return
	}
	if item.IsRadio() && item.checked == true {
		p.updateRadioGroup(item)
	}
}

func (p *Win32Menu) updateRadioGroups() {
	for menuItem := range p.radioGroups {
		if menuItem.checked {
			p.updateRadioGroup(menuItem)
		}
	}
}

func (p *Win32Menu) updateRadioGroup(item *MenuItem) {
	for _, radioGroup := range p.radioGroups[item] {
		thisMenuID := radioGroup.MenuID(item)
		startID, endID := radioGroup.Bounds()
		w32.CheckRadio(p.menu, startID, endID, thisMenuID)

	}
}

func (p *Win32Menu) OnMenuOpen(fn func()) {
	p.onMenuOpen = fn
}

func (p *Win32Menu) OnMenuClose(fn func()) {
	p.onMenuClose = fn
}
