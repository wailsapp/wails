package subsystem

import (
	"encoding/json"
	"github.com/wailsapp/wails/v2/internal/menumanager"
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

	// Service Bus
	bus *servicebus.ServiceBus

	// Menu Manager
	menuManager *menumanager.Manager
}

// NewMenu creates a new menu subsystem
func NewMenu(bus *servicebus.ServiceBus, logger *logger.Logger, menuManager *menumanager.Manager) (*Menu, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to menu messages
	menuChannel, err := bus.Subscribe("menu:")
	if err != nil {
		return nil, err
	}

	result := &Menu{
		quitChannel: quitChannel,
		menuChannel: menuChannel,
		logger:      logger.CustomLogger("Menu Subsystem"),
		listeners:   make(map[string][]func(*menu.MenuItem)),
		menuItems:   make(map[string]*menu.MenuItem),
		bus:         bus,
		menuManager: menuManager,
	}

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

					type ClickCallbackMessage struct {
						MenuItemID string `json:"menuItemID"`
						MenuType   string `json:"menuType"`
						Data       string `json:"data"`
					}

					var callbackData ClickCallbackMessage
					payload := []byte(menuMessage.Data().(string))
					err := json.Unmarshal(payload, &callbackData)
					if err != nil {
						m.logger.Error("%s", err.Error())
						return
					}

					err = m.menuManager.ProcessClick(callbackData.MenuItemID, callbackData.Data, callbackData.MenuType)
					if err != nil {
						m.logger.Trace("%s", err.Error())
					}

				case "on":
					listenerDetails := menuMessage.Data().(*message.MenuOnMessage)
					id := listenerDetails.MenuID
					m.listeners[id] = append(m.listeners[id], listenerDetails.Callback)

				// Make sure we catch any menu updates
				case "updateappmenu":
					updatedMenu, err := m.menuManager.UpdateApplicationMenu()
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}
					//updatedMenu := menuMessage.Data().(*menu.Menu)
					//m.processMenu(updatedMenu)
					//
					//// Notify frontend of menu change
					m.bus.Publish("menufrontend:updateappmenu", updatedMenu)

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

func (m *Menu) shutdown() {
	m.logger.Trace("Shutdown")
}
