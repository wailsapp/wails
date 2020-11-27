package subsystem

import (
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// eventListener holds a callback function which is invoked when
// the event listened for is emitted. It has a counter which indicates
// how the total number of events it is interested in. A value of zero
// means it does not expire (default).
// type eventListener struct {
// 	callback func(...interface{}) // Function to call with emitted event data
// 	counter  int                  // The number of times this callback may be called. -1 = infinite
// 	delete   bool                 // Flag to indicate that this listener should be deleted
// }

// Menu is the subsystem that handles the operation of menus. It manages all service bus messages
// starting with "menu".
type Menu struct {
	quitChannel <-chan *servicebus.Message
	menuChannel <-chan *servicebus.Message
	running     bool

	// Event listeners
	listeners  map[string][]func(*menu.MenuItem)
	menuItems  map[string]*menu.MenuItem
	notifyLock sync.RWMutex

	// logger
	logger logger.CustomLogger
}

// NewMenu creates a new menu subsystem
func NewMenu(initialMenu *menu.Menu, bus *servicebus.ServiceBus, logger *logger.Logger) (*Menu, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to menu messages
	menuChannel, err := bus.Subscribe("menu")
	if err != nil {
		return nil, err
	}

	result := &Menu{
		quitChannel: quitChannel,
		menuChannel: menuChannel,
		logger:      logger.CustomLogger("Menu Subsystem"),
		listeners:   make(map[string][]func(*menu.MenuItem)),
		menuItems:   make(map[string]*menu.MenuItem),
	}

	// Build up list of item/id pairs
	result.processMenu(initialMenu)

	return result, nil
}

// Start the subsystem
func (m *Menu) Start() error {

	m.logger.Trace("Starting")

	m.running = true

	// Spin off a go routine
	go func() {
		for m.running {
			select {
			case <-m.quitChannel:
				m.running = false
				break
			case menuMessage := <-m.menuChannel:
				splitTopic := strings.Split(menuMessage.Topic(), ":")
				menuMessageType := splitTopic[1]
				switch menuMessageType {
				case "clicked":
					if len(splitTopic) != 2 {
						m.logger.Error("Received clicked message with invalid topic format. Expected 2 sections in topic, got %s", splitTopic)
						continue
					}
					m.logger.Trace("Got Menu clicked Message: %s %+v", menuMessage.Topic(), menuMessage.Data())
					menuid := menuMessage.Data().(string)

					// Get the menu item
					menuItem := m.menuItems[menuid]
					if menuItem == nil {
						m.logger.Trace("Cannot process menuid %s - unknown", menuid)
						return
					}

					// Is the menu item a checkbox?
					if menuItem.Type == menu.CheckboxType {
						// Toggle state
						menuItem.Checked = !menuItem.Checked
					}

					// Notify listeners
					m.notifyListeners(menuid, menuItem)
				case "on":
					listenerDetails := menuMessage.Data().(*message.MenuOnMessage)
					id := listenerDetails.MenuID
					// Check we have a menu with that id
					if m.menuItems[id] == nil {
						m.logger.Error("cannot register listener for unknown menu id '%s'", id)
						continue
					}
					// We do! Append the callback
					m.listeners[id] = append(m.listeners[id], listenerDetails.Callback)
				default:
					m.logger.Error("unknown menu message: %+v", menuMessage)
				}
			}
		}

		// Call shutdown
		m.shutdown()
	}()

	return nil
}

func (m *Menu) processMenu(menu *menu.Menu) {
	for _, item := range menu.Items {
		m.processMenuItem(item)
	}
}

func (m *Menu) processMenuItem(item *menu.MenuItem) {

	if item.SubMenu != nil {
		for _, submenuitem := range item.SubMenu {
			m.processMenuItem(submenuitem)
		}
		return
	}

	if item.ID != "" {
		if m.menuItems[item.ID] != nil {
			m.logger.Error("Menu id '%s' is used by multiple menu items: %s %s", m.menuItems[item.ID].Label, item.Label)
			return
		}
		m.menuItems[item.ID] = item
	}
}

// Notifies listeners that the given menu was clicked
func (m *Menu) notifyListeners(menuid string, menuItem *menu.MenuItem) {

	// Get list of menu listeners
	listeners := m.listeners[menuid]
	if listeners == nil {
		m.logger.Trace("No listeners for MenuItem with ID '%s'", menuid)
		return
	}

	// Lock the listeners
	m.notifyLock.Lock()

	// Callback in goroutine
	for _, listener := range listeners {
		go listener(menuItem)
	}

	// Unlock
	m.notifyLock.Unlock()
}

func (m *Menu) shutdown() {
	m.logger.Trace("Shutdown")
}
