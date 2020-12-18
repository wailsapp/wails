package subsystem

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// ContextMenus is the subsystem that handles the operation of context menus. It manages all service bus messages
// starting with "contextmenus".
type ContextMenus struct {
	quitChannel <-chan *servicebus.Message
	menuChannel <-chan *servicebus.Message
	running     bool

	// Event listeners
	listeners  map[string][]func(*menu.MenuItem, string)
	menuItems  map[string]*menu.MenuItem
	notifyLock sync.RWMutex

	// logger
	logger logger.CustomLogger

	// The context menus
	contextMenus *menu.ContextMenus

	// Service Bus
	bus *servicebus.ServiceBus
}

// NewContextMenus creates a new context menu subsystem
func NewContextMenus(contextMenus *menu.ContextMenus, bus *servicebus.ServiceBus, logger *logger.Logger) (*ContextMenus, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to menu messages
	menuChannel, err := bus.Subscribe("contextmenus:")
	if err != nil {
		return nil, err
	}

	result := &ContextMenus{
		quitChannel:  quitChannel,
		menuChannel:  menuChannel,
		logger:       logger.CustomLogger("Context Menu Subsystem"),
		listeners:    make(map[string][]func(*menu.MenuItem, string)),
		menuItems:    make(map[string]*menu.MenuItem),
		contextMenus: contextMenus,
		bus:          bus,
	}

	// Build up list of item/id pairs
	result.processContextMenus(contextMenus)

	return result, nil
}

type contextMenuData struct {
	MenuItemID string `json:"menuItemID"`
	Data       string `json:"data"`
}

// Start the subsystem
func (c *ContextMenus) Start() error {

	c.logger.Trace("Starting")

	c.running = true

	// Spin off a go routine
	go func() {
		for c.running {
			select {
			case <-c.quitChannel:
				c.running = false
				break
			case menuMessage := <-c.menuChannel:
				splitTopic := strings.Split(menuMessage.Topic(), ":")
				menuMessageType := splitTopic[1]
				switch menuMessageType {
				case "clicked":
					if len(splitTopic) != 2 {
						c.logger.Error("Received clicked message with invalid topic format. Expected 2 sections in topic, got %s", splitTopic)
						continue
					}
					c.logger.Trace("Got Context Menu clicked Message: %s %+v", menuMessage.Topic(), menuMessage.Data())
					contextMenuDataJSON := menuMessage.Data().(string)

					var data contextMenuData
					err := json.Unmarshal([]byte(contextMenuDataJSON), &data)
					if err != nil {
						c.logger.Trace("Cannot process contextMenuDataJSON %s", string(contextMenuDataJSON))
						return
					}

					// Get the menu item
					menuItem := c.menuItems[data.MenuItemID]
					if menuItem == nil {
						c.logger.Trace("Cannot process menuitem id %s - unknown", data.MenuItemID)
						return
					}

					// Is the menu item a checkbox?
					if menuItem.Type == menu.CheckboxType {
						// Toggle state
						menuItem.Checked = !menuItem.Checked
					}

					// Notify listeners
					c.notifyListeners(data, menuItem)
				case "on":
					listenerDetails := menuMessage.Data().(*message.ContextMenusOnMessage)
					id := listenerDetails.MenuID
					c.listeners[id] = append(c.listeners[id], listenerDetails.Callback)

				// Make sure we catch any menu updates
				case "update":
					updatedMenu := menuMessage.Data().(*menu.ContextMenus)
					c.processContextMenus(updatedMenu)

					// Notify frontend of menu change
					c.bus.Publish("contextmenufrontend:update", updatedMenu)

				default:
					c.logger.Error("unknown context menu message: %+v", menuMessage)
				}
			}
		}

		// Call shutdown
		c.shutdown()
	}()

	return nil
}

func (c *ContextMenus) processContextMenus(contextMenus *menu.ContextMenus) {
	// Initialise the variables
	c.menuItems = make(map[string]*menu.MenuItem)
	c.contextMenus = contextMenus

	for _, contextMenu := range contextMenus.Items {
		for _, item := range contextMenu.Items {
			c.processMenuItem(item)
		}
	}
}

func (c *ContextMenus) processMenuItem(item *menu.MenuItem) {

	if item.SubMenu != nil {
		for _, submenuitem := range item.SubMenu {
			c.processMenuItem(submenuitem)
		}
		return
	}

	if item.ID != "" {
		if c.menuItems[item.ID] != nil {
			c.logger.Error("Context Menu id '%s' is used by multiple menu items: %s %s", c.menuItems[item.ID].Label, item.Label)
			return
		}
		c.menuItems[item.ID] = item
	}
}

// Notifies listeners that the given menu was clicked
func (c *ContextMenus) notifyListeners(contextData contextMenuData, menuItem *menu.MenuItem) {

	// Get list of menu listeners
	listeners := c.listeners[contextData.MenuItemID]
	if listeners == nil {
		c.logger.Trace("No listeners for MenuItem with ID '%s'", contextData.MenuItemID)
		return
	}

	// Lock the listeners
	c.notifyLock.Lock()

	// Callback in goroutine
	for _, listener := range listeners {
		go listener(menuItem, contextData.Data)
	}

	// Unlock
	c.notifyLock.Unlock()
}

func (c *ContextMenus) shutdown() {
	c.logger.Trace("Shutdown")
}
