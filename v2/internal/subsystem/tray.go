package subsystem

import (
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Tray is the subsystem that handles the operation of the tray menu.
// It manages all service bus messages starting with "tray".
type Tray struct {
	quitChannel <-chan *servicebus.Message
	trayChannel <-chan *servicebus.Message
	running     bool

	// Event listeners
	listeners  map[string][]func(*menu.MenuItem)
	menuItems  map[string]*menu.MenuItem
	notifyLock sync.RWMutex

	// logger
	logger logger.CustomLogger

	// The tray menu
	trayMenu *menu.Menu

	// Service Bus
	bus *servicebus.ServiceBus
}

// NewTray creates a new menu subsystem
func NewTray(trayMenu *menu.Menu, bus *servicebus.ServiceBus,
	logger *logger.Logger) (*Tray, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to menu messages
	trayChannel, err := bus.Subscribe("tray:")
	if err != nil {
		return nil, err
	}

	result := &Tray{
		quitChannel: quitChannel,
		trayChannel: trayChannel,
		logger:      logger.CustomLogger("Tray Subsystem"),
		listeners:   make(map[string][]func(*menu.MenuItem)),
		menuItems:   make(map[string]*menu.MenuItem),
		trayMenu:    trayMenu,
		bus:         bus,
	}

	// Build up list of item/id pairs
	result.processMenu(trayMenu)

	return result, nil
}

// Start the subsystem
func (t *Tray) Start() error {

	t.logger.Trace("Starting")

	t.running = true

	// Spin off a go routine
	go func() {
		for t.running {
			select {
			case <-t.quitChannel:
				t.running = false
				break
			case menuMessage := <-t.trayChannel:
				splitTopic := strings.Split(menuMessage.Topic(), ":")
				menuMessageType := splitTopic[1]
				switch menuMessageType {
				case "clicked":
					if len(splitTopic) != 2 {
						t.logger.Error("Received clicked message with invalid topic format. Expected 2 sections in topic, got %s", splitTopic)
						continue
					}
					t.logger.Trace("Got Tray Menu clicked Message: %s %+v", menuMessage.Topic(), menuMessage.Data())
					menuid := menuMessage.Data().(string)

					// Get the menu item
					menuItem := t.menuItems[menuid]
					if menuItem == nil {
						t.logger.Trace("Cannot process menuid %s - unknown", menuid)
						return
					}

					// Is the menu item a checkbox?
					if menuItem.Type == menu.CheckboxType {
						// Toggle state
						menuItem.Checked = !menuItem.Checked
					}

					// Notify listeners
					t.notifyListeners(menuid, menuItem)
				case "on":
					listenerDetails := menuMessage.Data().(*message.TrayOnMessage)
					id := listenerDetails.MenuID
					t.listeners[id] = append(t.listeners[id], listenerDetails.Callback)

				// Make sure we catch any menu updates
				case "update":
					updatedMenu := menuMessage.Data().(*menu.Menu)
					t.processMenu(updatedMenu)

					// Notify frontend of menu change
					t.bus.Publish("trayfrontend:update", updatedMenu)

				default:
					t.logger.Error("unknown tray message: %+v", menuMessage)
				}
			}
		}

		// Call shutdown
		t.shutdown()
	}()

	return nil
}

func (t *Tray) processMenu(trayMenu *menu.Menu) {
	// Initialise the variables
	t.menuItems = make(map[string]*menu.MenuItem)
	t.trayMenu = trayMenu

	for _, item := range trayMenu.Items {
		t.processMenuItem(item)
	}
}

func (t *Tray) processMenuItem(item *menu.MenuItem) {

	if item.SubMenu != nil {
		for _, submenuitem := range item.SubMenu {
			t.processMenuItem(submenuitem)
		}
		return
	}

	if item.ID != "" {
		if t.menuItems[item.ID] != nil {
			t.logger.Error("Menu id '%s' is used by multiple menu items: %s %s", t.menuItems[item.ID].Label, item.Label)
			return
		}
		t.menuItems[item.ID] = item
	}
}

// Notifies listeners that the given menu was clicked
func (t *Tray) notifyListeners(menuid string, menuItem *menu.MenuItem) {

	// Get list of menu listeners
	listeners := t.listeners[menuid]
	if listeners == nil {
		t.logger.Trace("No listeners for MenuItem with ID '%s'", menuid)
		return
	}

	// Lock the listeners
	t.notifyLock.Lock()

	// Callback in goroutine
	for _, listener := range listeners {
		go listener(menuItem)
	}

	// Unlock
	t.notifyLock.Unlock()
}

func (t *Tray) shutdown() {
	t.logger.Trace("Shutdown")
}
