package subsystem

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Menu is the subsystem that handles the operation of menus. It manages all service bus messages
// starting with "menu".
type Menu struct {
	menuChannel <-chan *servicebus.Message

	// shutdown flag
	shouldQuit bool

	// logger
	logger logger.CustomLogger

	// Service Bus
	bus *servicebus.ServiceBus

	// Menu Manager
	menuManager *menumanager.Manager

	// ctx
	ctx context.Context

	// parent waitgroup
	wg *sync.WaitGroup
}

// NewMenu creates a new menu subsystem
func NewMenu(ctx context.Context, bus *servicebus.ServiceBus, logger *logger.Logger, menuManager *menumanager.Manager) (*Menu, error) {

	// Subscribe to menu messages
	menuChannel, err := bus.Subscribe("menu:")
	if err != nil {
		return nil, err
	}

	result := &Menu{
		menuChannel: menuChannel,
		logger:      logger.CustomLogger("Menu Subsystem"),
		bus:         bus,
		menuManager: menuManager,
		ctx:         ctx,
		wg:          ctx.Value("waitgroup").(*sync.WaitGroup),
	}

	return result, nil
}

// Start the subsystem
func (m *Menu) Start() error {

	m.logger.Trace("Starting")

	m.wg.Add(1)

	// Spin off a go routine
	go func() {
		defer m.logger.Trace("Shutdown")
		for {
			select {
			case <-m.ctx.Done():
				m.wg.Done()
				return
			case menuMessage := <-m.menuChannel:
				splitTopic := strings.Split(menuMessage.Topic(), ":")
				menuMessageType := splitTopic[1]
				switch menuMessageType {
				case "ontrayopen":
					trayID := menuMessage.Data().(string)
					m.menuManager.OnTrayMenuOpen(trayID)
				case "ontrayclose":
					trayID := menuMessage.Data().(string)
					m.menuManager.OnTrayMenuClose(trayID)
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
						ParentID   string `json:"parentID"`
					}

					var callbackData ClickCallbackMessage
					payload := []byte(menuMessage.Data().(string))
					err := json.Unmarshal(payload, &callbackData)
					if err != nil {
						m.logger.Error("%s", err.Error())
						return
					}

					err = m.menuManager.ProcessClick(callbackData.MenuItemID, callbackData.Data, callbackData.MenuType, callbackData.ParentID)
					if err != nil {
						m.logger.Trace("%s", err.Error())
					}

				// Make sure we catch any menu updates
				case "updateappmenu":
					updatedMenu, err := m.menuManager.UpdateApplicationMenu()
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}

					// Notify frontend of menu change
					m.bus.Publish("menufrontend:updateappmenu", updatedMenu)

				case "updatecontextmenu":
					contextMenu := menuMessage.Data().(*menu.ContextMenu)
					updatedMenu, err := m.menuManager.UpdateContextMenu(contextMenu)
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}

					// Notify frontend of menu change
					m.bus.Publish("menufrontend:updatecontextmenu", updatedMenu)

				case "settraymenu":
					trayMenu := menuMessage.Data().(*menu.TrayMenu)
					updatedMenu, err := m.menuManager.SetTrayMenu(trayMenu)
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}

					// Notify frontend of menu change
					m.bus.Publish("menufrontend:settraymenu", updatedMenu)

				case "deletetraymenu":
					trayMenu := menuMessage.Data().(*menu.TrayMenu)
					trayID, err := m.menuManager.GetTrayID(trayMenu)
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}

					// Notify frontend of menu change
					m.bus.Publish("menufrontend:deletetraymenu", trayID)

				case "updatetraymenulabel":
					trayMenu := menuMessage.Data().(*menu.TrayMenu)
					updatedLabel, err := m.menuManager.UpdateTrayMenuLabel(trayMenu)
					if err != nil {
						m.logger.Trace("%s", err.Error())
						return
					}

					// Notify frontend of menu change
					m.bus.Publish("menufrontend:updatetraymenulabel", updatedLabel)

				default:
					m.logger.Error("unknown menu message: %+v", menuMessage)
				}
			}
		}
	}()

	return nil
}
