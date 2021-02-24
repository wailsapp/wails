package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Menu defines all Menu related operations
type Menu interface {
	UpdateApplicationMenu()
	UpdateContextMenu(contextMenu *menu.ContextMenu)
	SetTrayMenu(trayMenu *menu.TrayMenu)
	DeleteTrayMenu(trayMenu *menu.TrayMenu)
	UpdateTrayMenuLabel(trayMenu *menu.TrayMenu)
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

func (m *menuRuntime) UpdateContextMenu(contextMenu *menu.ContextMenu) {
	m.bus.Publish("menu:updatecontextmenu", contextMenu)
}

func (m *menuRuntime) SetTrayMenu(trayMenu *menu.TrayMenu) {
	m.bus.Publish("menu:settraymenu", trayMenu)
}

func (m *menuRuntime) UpdateTrayMenuLabel(trayMenu *menu.TrayMenu) {
	m.bus.Publish("menu:updatetraymenulabel", trayMenu)
}

func (m *menuRuntime) DeleteTrayMenu(trayMenu *menu.TrayMenu) {
	m.bus.Publish("menu:deletetraymenu", trayMenu)
}
