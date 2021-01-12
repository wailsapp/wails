package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Menu defines all Menu related operations
type Menu interface {
	UpdateApplicationMenu()
}

type menuRuntime struct {
	bus *servicebus.ServiceBus
}

// newMenu creates a new Menu struct
func newMenu(bus *servicebus.ServiceBus) Menu {
	return &menuRuntime{
		bus: bus,
	}
}

func (m *menuRuntime) UpdateApplicationMenu() {
	m.bus.Publish("menu:updateappmenu", nil)
}
