//go:build windows

package systray

import (
	"errors"
	"fmt"
	platformMenu "github.com/wailsapp/wails/v2/internal/platform/menu"
	"github.com/wailsapp/wails/v2/internal/platform/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

type RadioGroupMember struct {
	ID       int
	MenuItem *menu.MenuItem
}

type RadioGroup []*RadioGroupMember

func (r *RadioGroup) Add(id int, item *menu.MenuItem) {
	*r = append(*r, &RadioGroupMember{
		ID:       id,
		MenuItem: item,
	})
}

func (r *RadioGroup) Bounds() (int, int) {
	p := *r
	return p[0].ID, p[len(p)-1].ID
}

func (r *RadioGroup) MenuID(item *menu.MenuItem) int {
	for _, member := range *r {
		if member.MenuItem == item {
			return member.ID
		}
	}
	panic("RadioGroup.MenuID: item not found:")
}

type PopupMenu struct {
	menu          win32.PopupMenu
	parent        win32.HWND
	menuMapping   map[int]*menu.MenuItem
	checkboxItems map[*menu.MenuItem][]int
	radioGroups   map[*menu.MenuItem][]*RadioGroup
	menuData      *menu.Menu
	currentMenuID int
	onMenuClose   func()
	onMenuOpen    func()
}

func (p *PopupMenu) buildMenu(parentMenu win32.PopupMenu, inputMenu *menu.Menu) error {
	var currentRadioGroup RadioGroup
	for _, item := range inputMenu.Items {
		if item.Hidden {
			continue
		}
		var ret bool
		p.currentMenuID++
		itemID := p.currentMenuID
		p.menuMapping[itemID] = item

		flags := win32.MF_STRING
		if item.Disabled {
			flags = flags | win32.MF_GRAYED
		}
		if item.Checked {
			flags = flags | win32.MF_CHECKED
		}
		//if item.BarBreak {
		//	flags = flags | win32.MF_MENUBARBREAK
		//}
		if item.IsSeparator() {
			flags = flags | win32.MF_SEPARATOR
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

		if item.SubMenu != nil {
			flags = flags | win32.MF_POPUP
			submenu := win32.CreatePopupMenu()
			err := p.buildMenu(submenu, item.SubMenu)
			if err != nil {
				return err
			}
			itemID = int(submenu)
		}

		var menuText = item.Label
		if item.Accelerator != nil {
			shortcut := win32.AcceleratorToShortcut(item.Accelerator)
			menuText = fmt.Sprintf("%s\t%s", menuText, shortcut)
			// Popup Menus don't appear to support accelerators and I'm not
			// sure they make sense either
		}

		ret = parentMenu.Append(uintptr(flags), uintptr(itemID), menuText)
		if ret == false {
			return errors.New("AppendMenu failed")
		}
	}
	if len(currentRadioGroup) > 0 {
		for _, radioMember := range currentRadioGroup {
			currentRadioGroup := currentRadioGroup
			p.radioGroups[radioMember.MenuItem] = append(p.radioGroups[radioMember.MenuItem], &currentRadioGroup)
		}
		currentRadioGroup = RadioGroup{}
	}
	return nil
}

func (p *PopupMenu) Update() error {
	p.menu = win32.CreatePopupMenu()
	p.menuMapping = make(map[int]*menu.MenuItem)
	p.currentMenuID = win32.MenuItemMsgID
	err := p.buildMenu(p.menu, p.menuData)
	if err != nil {
		return err
	}
	p.updateRadioGroups()
	return nil
}

func NewPopupMenu(parent win32.HWND, inputMenu *menu.Menu) (*PopupMenu, error) {
	result := &PopupMenu{
		parent:        parent,
		menuData:      inputMenu,
		checkboxItems: make(map[*menu.MenuItem][]int),
		radioGroups:   make(map[*menu.MenuItem][]*RadioGroup),
	}
	err := result.Update()
	platformMenu.MenuManager.AddMenu(inputMenu, result.UpdateMenuItem)
	return result, err
}

func (p *PopupMenu) ShowAtCursor() error {
	x, y, ok := win32.GetCursorPos()
	if ok == false {
		return errors.New("GetCursorPos failed")
	}

	if win32.SetForegroundWindow(p.parent) == false {
		return errors.New("SetForegroundWindow failed")
	}

	if p.onMenuOpen != nil {
		p.onMenuOpen()
	}

	if p.menu.Track(win32.TPM_LEFTALIGN, x, y-5, p.parent) == false {
		return errors.New("TrackPopupMenu failed")
	}

	if p.onMenuClose != nil {
		p.onMenuClose()
	}

	if win32.PostMessage(p.parent, win32.WM_NULL, 0, 0) == 0 {
		return errors.New("PostMessage failed")
	}

	return nil
}

func (p *PopupMenu) ProcessCommand(cmdMsgID int) {
	item := p.menuMapping[cmdMsgID]
	platformMenu.MenuManager.ProcessClick(item)
}

func (p *PopupMenu) Destroy() {
	p.menu.Destroy()
}

func (p *PopupMenu) UpdateMenuItem(item *menu.MenuItem) {
	if item.IsCheckbox() {
		for _, itemID := range p.checkboxItems[item] {
			p.menu.Check(uintptr(itemID), item.Checked)
		}
		return
	}
	if item.IsRadio() && item.Checked == true {
		p.updateRadioGroup(item)
	}
}

func (p *PopupMenu) updateRadioGroups() {
	for menuItem := range p.radioGroups {
		if menuItem.Checked {
			p.updateRadioGroup(menuItem)
		}
	}
}

func (p *PopupMenu) updateRadioGroup(item *menu.MenuItem) {
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
