package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Menu defines all Menu related operations
type Menu interface {
	UpdateApplicationMenu()
	GetByID(menuID string) *menu.MenuItem
	RemoveByID(id string) bool
}

type menuRuntime struct {
	bus  *servicebus.ServiceBus
	menu *menu.Menu
}

// newMenu creates a new Menu struct
func newMenu(bus *servicebus.ServiceBus, menu *menu.Menu) Menu {
	return &menuRuntime{
		bus:  bus,
		menu: menu,
	}
}

func (m *menuRuntime) UpdateApplicationMenu() {
	m.bus.Publish("menu:updateappmenu", nil)
}

func (m *menuRuntime) GetByID(menuID string) *menu.MenuItem {
	return m.menu.GetByID(menuID)
}

func (m *menuRuntime) RemoveByID(id string) bool {
	return m.menu.RemoveByID(id)
}
