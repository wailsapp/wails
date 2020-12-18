package runtime

import (
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Tray defines all Tray related operations
type Tray interface {
	On(menuID string, callback func(*menu.MenuItem))
	Update()
	GetByID(menuID string) *menu.MenuItem
	RemoveByID(id string) bool
}

type trayRuntime struct {
	bus      *servicebus.ServiceBus
	trayMenu *menu.Menu
}

// newTray creates a new Menu struct
func newTray(bus *servicebus.ServiceBus, menu *menu.Menu) Tray {
	return &trayRuntime{
		bus:      bus,
		trayMenu: menu,
	}
}

// On registers a listener for a particular event
func (t *trayRuntime) On(menuID string, callback func(*menu.MenuItem)) {
	t.bus.Publish("tray:on", &message.TrayOnMessage{
		MenuID:   menuID,
		Callback: callback,
	})
}

func (t *trayRuntime) Update() {
	t.bus.Publish("tray:update", t.trayMenu)
}

func (t *trayRuntime) GetByID(menuID string) *menu.MenuItem {
	return t.trayMenu.GetByID(menuID)
}

func (t *trayRuntime) RemoveByID(id string) bool {
	return t.trayMenu.RemoveByID(id)
}
