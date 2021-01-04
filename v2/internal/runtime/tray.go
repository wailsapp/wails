package runtime

import (
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Tray defines all Tray related operations
type Tray interface {
	NewTray(id string) *menu.Tray
	On(menuID string, callback func(*menu.MenuItem))
	Update(tray ...*menu.Tray)
	GetByID(menuID string) *menu.MenuItem
	RemoveByID(id string) bool
	SetLabel(label string)
	SetIcon(name string)
}

type trayRuntime struct {
	bus      *servicebus.ServiceBus
	trayMenu *menu.Tray
}

// newTray creates a new Menu struct
func newTray(bus *servicebus.ServiceBus, menu *menu.Tray) Tray {
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

// NewTray creates a new Tray item
func (t *trayRuntime) NewTray(trayID string) *menu.Tray {
	return &menu.Tray{
		ID: trayID,
	}
}

func (t *trayRuntime) Update(tray ...*menu.Tray) {

	//trayToUpdate := t.trayMenu
	t.bus.Publish("tray:update", t.trayMenu)
}

func (t *trayRuntime) SetLabel(label string) {
	t.bus.Publish("tray:setlabel", label)
}
func (t *trayRuntime) SetIcon(name string) {
	t.bus.Publish("tray:seticon", name)
}

func (t *trayRuntime) GetByID(menuID string) *menu.MenuItem {
	return t.trayMenu.Menu.GetByID(menuID)
}

func (t *trayRuntime) RemoveByID(id string) bool {
	return t.trayMenu.Menu.RemoveByID(id)
}
