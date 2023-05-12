package application

import (
	"fmt"
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

type PopupMenu struct {
	menu          w32.PopupMenu
	parent        w32.HWND
	menuMapping   map[int]*MenuItem
	checkboxItems map[*MenuItem][]int
	radioGroups   map[*MenuItem][]*RadioGroup
	menuData      *Menu
	currentMenuID int
	onMenuClose   func()
	onMenuOpen    func()
}

func (p *PopupMenu) buildMenu(parentMenu w32.PopupMenu, inputMenu *Menu) {
	var currentRadioGroup RadioGroup
	for _, item := range inputMenu.items {
		if item.Hidden() {
			continue
		}
		p.currentMenuID++
		itemID := p.currentMenuID
		p.menuMapping[itemID] = item

		menuItemImpl := newMenuItemImpl(item, w32.HWND(parentMenu), itemID)

		flags := uint32(w32.MF_STRING)
		if item.disabled {
			flags = flags | w32.MF_GRAYED
		}
		if item.checked && item.IsCheckbox() {
			flags = flags | w32.MF_CHECKED
		}
		if item.IsSeparator() {
			flags = flags | w32.MF_SEPARATOR
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
			newSubmenu := w32.CreatePopupMenu()
			p.buildMenu(newSubmenu, item.submenu)
			itemID = int(newSubmenu)
			menuItemImpl.submenu = w32.HWND(newSubmenu)
		}

		var menuText = item.Label()

		ok := parentMenu.Append(flags, uintptr(itemID), menuText)
		if !ok {
			w32.Fatal(fmt.Sprintf("Error adding menu item: %s", menuText))
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

func (p *PopupMenu) Update() {
	p.menu = w32.CreatePopupMenu()
	p.menuMapping = make(map[int]*MenuItem)
	p.currentMenuID = MenuItemMsgID
	p.buildMenu(p.menu, p.menuData)
	p.updateRadioGroups()
}

func NewPopupMenu(parent w32.HWND, inputMenu *Menu) *PopupMenu {
	result := &PopupMenu{
		parent:        parent,
		menuData:      inputMenu,
		checkboxItems: make(map[*MenuItem][]int),
		radioGroups:   make(map[*MenuItem][]*RadioGroup),
	}
	result.Update()
	return result
}

func (p *PopupMenu) ShowAtCursor() {
	x, y, ok := w32.GetCursorPos()
	if ok == false {
		w32.Fatal("GetCursorPos failed")
	}

	w32.SetForegroundWindow(p.parent)

	if p.onMenuOpen != nil {
		p.onMenuOpen()
	}

	if p.menu.Track(p.parent, w32.TPM_LEFTALIGN, int32(x), int32(y-5)) == false {
		w32.Fatal("TrackPopupMenu failed")
	}

	if p.onMenuClose != nil {
		p.onMenuClose()
	}

	if !w32.PostMessage(p.parent, w32.WM_NULL, 0, 0) {
		w32.Fatal("PostMessage failed")
	}

}

func (p *PopupMenu) ProcessCommand(cmdMsgID int) {
	item := p.menuMapping[cmdMsgID]
	if item == nil {
		return
	}
	if item.IsRadio() {
		item.checked = true
		p.updateRadioGroup(item)
	}
	if item.callback != nil {
		item.handleClick()
	}
}

func (p *PopupMenu) Destroy() {
	p.menu.Destroy()
}

func (p *PopupMenu) UpdateMenuItem(item *MenuItem) {
	if item.IsCheckbox() {
		for _, itemID := range p.checkboxItems[item] {
			p.menu.Check(uintptr(itemID), item.checked)
		}
		return
	}
	if item.IsRadio() && item.checked == true {
		p.updateRadioGroup(item)
	}
}

func (p *PopupMenu) updateRadioGroups() {
	for menuItem := range p.radioGroups {
		if menuItem.checked {
			p.updateRadioGroup(menuItem)
		}
	}
}

func (p *PopupMenu) updateRadioGroup(item *MenuItem) {
	for _, radioGroup := range p.radioGroups[item] {
		thisMenuID := radioGroup.MenuID(item)
		startID, endID := radioGroup.Bounds()
		p.menu.CheckRadio(startID, endID, thisMenuID)
	}
}

func (p *PopupMenu) OnMenuOpen(fn func()) {
	p.onMenuOpen = fn
}

func (p *PopupMenu) OnMenuClose(fn func()) {
	p.onMenuClose = fn
}
