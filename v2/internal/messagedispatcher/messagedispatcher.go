package messagedispatcher

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/options/dialog"

	"github.com/wailsapp/wails/v2/internal/crypto"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Dispatcher translates messages received from the frontend
// and publishes them onto the service bus
type Dispatcher struct {
	quitChannel   <-chan *servicebus.Message
	resultChannel <-chan *servicebus.Message
	eventChannel  <-chan *servicebus.Message
	windowChannel <-chan *servicebus.Message
	dialogChannel <-chan *servicebus.Message
	systemChannel <-chan *servicebus.Message
	menuChannel   <-chan *servicebus.Message

	servicebus *servicebus.ServiceBus
	logger     logger.CustomLogger

	// Clients
	clients map[string]*DispatchClient
	lock    sync.RWMutex

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// internal wait group
	wg sync.WaitGroup
}

// New dispatcher. Needs a service bus to send to.
func New(servicebus *servicebus.ServiceBus, logger *logger.Logger) (*Dispatcher, error) {
	// Subscribe to call result messages
	resultChannel, err := servicebus.Subscribe("call:result")
	if err != nil {
		return nil, err
	}

	// Subscribe to event messages
	eventChannel, err := servicebus.Subscribe("event:emit")
	if err != nil {
		return nil, err
	}

	// Subscribe to quit messages
	quitChannel, err := servicebus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to window messages
	windowChannel, err := servicebus.Subscribe("window")
	if err != nil {
		return nil, err
	}

	// Subscribe to dialog events
	dialogChannel, err := servicebus.Subscribe("dialog:select")
	if err != nil {
		return nil, err
	}

	systemChannel, err := servicebus.Subscribe("system:")
	if err != nil {
		return nil, err
	}

	menuChannel, err := servicebus.Subscribe("menufrontend:")
	if err != nil {
		return nil, err
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	result := &Dispatcher{
		servicebus:    servicebus,
		eventChannel:  eventChannel,
		logger:        logger.CustomLogger("Message Dispatcher"),
		clients:       make(map[string]*DispatchClient),
		resultChannel: resultChannel,
		quitChannel:   quitChannel,
		windowChannel: windowChannel,
		dialogChannel: dialogChannel,
		systemChannel: systemChannel,
		menuChannel:   menuChannel,
		ctx:           ctx,
		cancel:        cancel,
	}

	return result, nil
}

// Start the subsystem
func (d *Dispatcher) Start() error {

	d.logger.Trace("Starting")

	d.wg.Add(1)

	// Spin off a go routine
	go func() {
		defer d.logger.Trace("Shutdown")
		for {
			select {
			case <-d.ctx.Done():
				d.wg.Done()
				return
			case <-d.quitChannel:
				d.processQuit()
			case resultMessage := <-d.resultChannel:
				d.processCallResult(resultMessage)
			case eventMessage := <-d.eventChannel:
				d.processEvent(eventMessage)
			case windowMessage := <-d.windowChannel:
				d.processWindowMessage(windowMessage)
			case dialogMessage := <-d.dialogChannel:
				d.processDialogMessage(dialogMessage)
			case systemMessage := <-d.systemChannel:
				d.processSystemMessage(systemMessage)
			case menuMessage := <-d.menuChannel:
				d.processMenuMessage(menuMessage)
			}
		}
	}()

	return nil
}

func (d *Dispatcher) processQuit() {
	d.lock.RLock()
	defer d.lock.RUnlock()
	for _, client := range d.clients {
		client.frontend.Quit()
	}
}

// RegisterClient will register the given callback with the dispatcher
// and return a DispatchClient that the caller can use to send messages
func (d *Dispatcher) RegisterClient(client Client) *DispatchClient {
	d.lock.Lock()
	defer d.lock.Unlock()

	// Create ID
	id := d.getUniqueID()
	d.clients[id] = newDispatchClient(id, client, d.logger, d.servicebus)

	return d.clients[id]
}

// RemoveClient will remove the registered client
func (d *Dispatcher) RemoveClient(dc *DispatchClient) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.clients, dc.id)
}

func (d *Dispatcher) getUniqueID() string {
	var uid string
	for {
		uid = crypto.RandomID()

		if d.clients[uid] == nil {
			break
		}
	}

	return uid
}

func (d *Dispatcher) processCallResult(result *servicebus.Message) {
	target := result.Target()

	if target == "" {
		// This is an error. Calls are 1:1!
		d.logger.Fatal("No target for call result: %+v", result)
	}

	d.lock.RLock()
	client := d.clients[target]
	d.lock.RUnlock()
	if client == nil {
		// This is fatal - unknown target!
		d.logger.Fatal("Unknown target for call result: %+v", result)
		return
	}

	d.logger.Trace("Sending message to client %s: R%s", target, result.Data().(string))
	client.frontend.CallResult(result.Data().(string))
}

// processSystem
func (d *Dispatcher) processSystemMessage(result *servicebus.Message) {

	d.logger.Trace("Got system in message dispatcher: %+v", result)

	splitTopic := strings.Split(result.Topic(), ":")
	command := splitTopic[1]
	callbackID := splitTopic[2]
	switch command {
	case "isdarkmode":
		d.lock.RLock()
		for _, client := range d.clients {
			client.frontend.DarkModeEnabled(callbackID)
			break
		}
		d.lock.RUnlock()

	default:
		d.logger.Error("Unknown system command: %s", command)
	}
}

// processEvent will
func (d *Dispatcher) processEvent(result *servicebus.Message) {

	d.logger.Trace("Got event in message dispatcher: %+v", result)

	splitTopic := strings.Split(result.Topic(), ":")
	eventType := splitTopic[1]
	switch eventType {
	case "emit":
		eventFrom := splitTopic[3]
		if eventFrom == "g" {
			// This was sent from Go - notify frontend
			eventData := result.Data().(*message.EventMessage)
			// Unpack event
			payload, err := json.Marshal(eventData)
			if err != nil {
				d.logger.Error("Unable to marshal eventData: %s", err.Error())
				return
			}
			d.lock.RLock()
			for _, client := range d.clients {
				client.frontend.NotifyEvent(string(payload))
			}
			d.lock.RUnlock()
		}
	default:
		d.logger.Error("Unknown event type: %s", eventType)
	}
}

// processWindowMessage processes messages intended for the window
func (d *Dispatcher) processWindowMessage(result *servicebus.Message) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	splitTopic := strings.Split(result.Topic(), ":")
	command := splitTopic[1]
	switch command {
	case "settitle":
		title, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid title for 'window:settitle' : %#v", result.Data())
			return
		}
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowSetTitle(title)
		}
	case "fullscreen":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowFullscreen()
		}
	case "unfullscreen":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowUnFullscreen()
		}
	case "setcolour":
		colour, ok := result.Data().(int)
		if !ok {
			d.logger.Error("Invalid colour for 'window:setcolour' : %#v", result.Data())
			return
		}
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowSetColour(colour)
		}
	case "show":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowShow()
		}
	case "hide":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowHide()
		}
	case "center":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowCenter()
		}
	case "maximise":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowMaximise()
		}
	case "unmaximise":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowUnmaximise()
		}
	case "minimise":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowMinimise()
		}
	case "unminimise":
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowUnminimise()
		}
	case "position":
		// We need 2 arguments
		if len(splitTopic) != 4 {
			d.logger.Error("Invalid number of parameters for 'window:position' : %#v", result.Data())
			return
		}
		x, err1 := strconv.Atoi(splitTopic[2])
		y, err2 := strconv.Atoi(splitTopic[3])
		if err1 != nil || err2 != nil {
			d.logger.Error("Invalid integer parameters for 'window:position' : %#v", result.Data())
			return
		}
		// Notify clients
		for _, client := range d.clients {
			client.frontend.WindowPosition(x, y)
		}
	case "size":
		// We need 2 arguments
		if len(splitTopic) != 4 {
			d.logger.Error("Invalid number of parameters for 'window:size' : %#v", result.Data())
			return
		}
		w, err1 := strconv.Atoi(splitTopic[2])
		h, err2 := strconv.Atoi(splitTopic[3])
		if err1 != nil || err2 != nil {
			d.logger.Error("Invalid integer parameters for 'window:size' : %#v", result.Data())
			return
		}
		// Notifh clients
		for _, client := range d.clients {
			client.frontend.WindowSize(w, h)
		}
	case "minsize":
		// We need 2 arguments
		if len(splitTopic) != 4 {
			d.logger.Error("Invalid number of parameters for 'window:minsize' : %#v", result.Data())
			return
		}
		w, err1 := strconv.Atoi(splitTopic[2])
		h, err2 := strconv.Atoi(splitTopic[3])
		if err1 != nil || err2 != nil {
			d.logger.Error("Invalid integer parameters for 'window:minsize' : %#v", result.Data())
			return
		}
		// Notifh clients
		for _, client := range d.clients {
			client.frontend.WindowSetMinSize(w, h)
		}
	case "maxsize":
		// We need 2 arguments
		if len(splitTopic) != 4 {
			d.logger.Error("Invalid number of parameters for 'window:maxsize' : %#v", result.Data())
			return
		}
		w, err1 := strconv.Atoi(splitTopic[2])
		h, err2 := strconv.Atoi(splitTopic[3])
		if err1 != nil || err2 != nil {
			d.logger.Error("Invalid integer parameters for 'window:maxsize' : %#v", result.Data())
			return
		}
		// Notifh clients
		for _, client := range d.clients {
			client.frontend.WindowSetMaxSize(w, h)
		}
	default:
		d.logger.Error("Unknown window command: %s", command)
	}
	d.logger.Trace("Got window in message dispatcher: %+v", result)

}

// processDialogMessage processes dialog messages
func (d *Dispatcher) processDialogMessage(result *servicebus.Message) {
	splitTopic := strings.Split(result.Topic(), ":")
	if len(splitTopic) < 4 {
		d.logger.Error("Invalid dialog message : %#v", result.Data())
		return
	}

	command := splitTopic[1]
	switch command {
	case "select":
		dialogType := splitTopic[2]
		switch dialogType {
		case "open":
			dialogOptions, ok := result.Data().(*dialog.OpenDialog)
			if !ok {
				d.logger.Error("Invalid data for 'dialog:select:open' : %#v", result.Data())
				return
			}
			// This is hardcoded in the sender too
			callbackID := splitTopic[3]

			// TODO: Work out what we mean in a multi window environment...
			// For now we will just pick the first one
			for _, client := range d.clients {
				client.frontend.OpenDialog(dialogOptions, callbackID)
			}
		case "save":
			dialogOptions, ok := result.Data().(*dialog.SaveDialog)
			if !ok {
				d.logger.Error("Invalid data for 'dialog:select:save' : %#v", result.Data())
				return
			}
			// This is hardcoded in the sender too
			callbackID := splitTopic[3]

			// TODO: Work out what we mean in a multi window environment...
			// For now we will just pick the first one
			for _, client := range d.clients {
				client.frontend.SaveDialog(dialogOptions, callbackID)
			}
		case "message":
			dialogOptions, ok := result.Data().(*dialog.MessageDialog)
			if !ok {
				d.logger.Error("Invalid data for 'dialog:select:message' : %#v", result.Data())
				return
			}
			// This is hardcoded in the sender too
			callbackID := splitTopic[3]

			// TODO: Work out what we mean in a multi window environment...
			// For now we will just pick the first one
			for _, client := range d.clients {
				client.frontend.MessageDialog(dialogOptions, callbackID)
			}
		default:
			d.logger.Error("Unknown dialog type: %s", dialogType)
		}

	default:
		d.logger.Error("Unknown dialog command: %s", command)
	}

}

func (d *Dispatcher) processMenuMessage(result *servicebus.Message) {
	splitTopic := strings.Split(result.Topic(), ":")
	if len(splitTopic) < 2 {
		d.logger.Error("Invalid menu message : %#v", result.Data())
		return
	}

	command := splitTopic[1]
	switch command {
	case "updateappmenu":

		updatedMenu, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid data for 'menufrontend:updateappmenu' : %#v",
				result.Data())
			return
		}

		// TODO: Work out what we mean in a multi window environment...
		// For now we will just pick the first one
		for _, client := range d.clients {
			client.frontend.SetApplicationMenu(updatedMenu)
		}

	case "settraymenu":
		trayMenuJSON, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid data for 'menufrontend:settraymenu' : %#v",
				result.Data())
			return
		}

		// TODO: Work out what we mean in a multi window environment...
		// For now we will just pick the first one
		for _, client := range d.clients {
			client.frontend.SetTrayMenu(trayMenuJSON)
		}

	case "updatecontextmenu":
		updatedContextMenu, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid data for 'menufrontend:updatecontextmenu' : %#v",
				result.Data())
			return
		}

		// TODO: Work out what we mean in a multi window environment...
		// For now we will just pick the first one
		for _, client := range d.clients {
			client.frontend.UpdateContextMenu(updatedContextMenu)
		}

	case "updatetraymenulabel":
		updatedTrayMenuLabel, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid data for 'menufrontend:updatetraymenulabel' : %#v",
				result.Data())
			return
		}

		// TODO: Work out what we mean in a multi window environment...
		// For now we will just pick the first one
		for _, client := range d.clients {
			client.frontend.UpdateTrayMenuLabel(updatedTrayMenuLabel)
		}
	case "deletetraymenu":
		traymenuid, ok := result.Data().(string)
		if !ok {
			d.logger.Error("Invalid data for 'menufrontend:updatetraymenulabel' : %#v",
				result.Data())
			return
		}

		for _, client := range d.clients {
			client.frontend.DeleteTrayMenuByID(traymenuid)
		}

	default:
		d.logger.Error("Unknown menufrontend command: %s", command)
	}
}

func (d *Dispatcher) Close() {
	d.cancel()
	d.wg.Wait()
}
