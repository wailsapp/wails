package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Menu defines all Menu related operations
type Menu interface {
	UpdateApplicationMenu()
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
