package main

import (
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// ContextMenu struct
type ContextMenu struct {
	runtime *wails.Runtime
	counter int
	lock    sync.Mutex
}

// WailsInit is called at application startup
func (c *ContextMenu) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	c.runtime = runtime

	return nil
}

func createContextMenus() *menu.ContextMenus {
	result := menu.NewContextMenus()
	result.AddMenu("test", menu.NewMenuFromItems(
		menu.Text("Clicked 0 times", nil, nil),
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
	return result
}
