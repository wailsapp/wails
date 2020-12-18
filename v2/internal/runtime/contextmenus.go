package runtime

import (
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// ContextMenus defines all ContextMenu related operations
type ContextMenus interface {
	On(menuID string, callback func(*menu.MenuItem, string))
	Update()
	GetByID(menuID string) *menu.MenuItem
	RemoveByID(id string) bool
}

type contextMenus struct {
	bus          *servicebus.ServiceBus
	contextmenus *menu.ContextMenus
}

// newContextMenus creates a new ContextMenu struct
func newContextMenus(bus *servicebus.ServiceBus, contextmenus *menu.ContextMenus) ContextMenus {
	return &contextMenus{
		bus:          bus,
		contextmenus: contextmenus,
	}
}

// On registers a listener for a particular event
func (t *contextMenus) On(menuID string, callback func(*menu.MenuItem, string)) {
	t.bus.Publish("contextmenus:on", &message.ContextMenusOnMessage{
		MenuID:   menuID,
		Callback: callback,
	})
}

func (t *contextMenus) Update() {
	t.bus.Publish("contextmenus:update", t.contextmenus)
}

func (t *contextMenus) GetByID(menuItemID string) *menu.MenuItem {
	return t.contextmenus.GetByID(menuItemID)
}

func (t *contextMenus) RemoveByID(menuItemID string) bool {
	return t.contextmenus.RemoveByID(menuItemID)
}
