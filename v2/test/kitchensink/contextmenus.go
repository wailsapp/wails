package main

import (
	"fmt"
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

	// Setup Menu Listeners
	c.runtime.ContextMenu.On("Test Context Menu", func(mi *menu.MenuItem, contextData string) {
		fmt.Printf("\n\nContext Data = '%s'\n\n", contextData)
		c.lock.Lock()
		c.counter++
		mi.Label = fmt.Sprintf("Clicked %d times", c.counter)
		c.lock.Unlock()
		c.runtime.ContextMenu.Update()
	})

	return nil
}

func createContextMenus() *menu.ContextMenus {
	result := menu.NewContextMenus()
	result.AddMenu("test", menu.NewMenuFromItems(menu.Text("Clicked 0 times", "Test Context Menu")))
	return result
}
