package main

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// ContextMenu struct
type ContextMenu struct {
	runtime         *wails.Runtime
	counter         int
	lock            sync.Mutex
	testContextMenu *menu.ContextMenu
	clickedMenu     *menu.MenuItem
}

// WailsInit is called at application startup
func (c *ContextMenu) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	c.runtime = runtime

	return nil
}

// Setup Menu Listeners
func (c *ContextMenu) updateContextMenu(_ *menu.CallbackData) {
	c.lock.Lock()
	c.counter++
	c.clickedMenu.Label = fmt.Sprintf("Clicked %d times", c.counter)
	c.lock.Unlock()
	c.runtime.Menu.UpdateContextMenu(c.testContextMenu)
}

func (c *ContextMenu) createContextMenus() []*menu.ContextMenu {
	c.clickedMenu = menu.Text("Clicked 0 times", nil, c.updateContextMenu)
	c.testContextMenu = menu.NewContextMenu("test", menu.NewMenuFromItems(
		c.clickedMenu,
		menu.Separator(),
		menu.Checkbox("I am a checkbox", false, nil, nil),
		menu.Separator(),
		menu.Radio("Radio Option 1", true, nil, nil),
		menu.Radio("Radio Option 2", false, nil, nil),
		menu.Radio("Radio Option 3", false, nil, nil),
		menu.Separator(),
		menu.SubMenu("A Submenu", menu.NewMenuFromItems(
			menu.Text("Hello", nil, nil),
		)),
	))
	return []*menu.ContextMenu{
		c.testContextMenu,
	}
}
