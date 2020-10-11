package subsystem

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Call is the Call subsystem. It manages all service bus messages
// starting with "call".
type Call struct {
	quitChannel <-chan *servicebus.Message
	callChannel <-chan *servicebus.Message
	running     bool

	// bindings DB
	DB *binding.DB

	// ServiceBus
	bus *servicebus.ServiceBus

	// logger
	logger logger.CustomLogger

	// runtime
	runtime *runtime.Runtime
}

// NewCall creates a new call subsystem
func NewCall(bus *servicebus.ServiceBus, logger *logger.Logger, DB *binding.DB, runtime *runtime.Runtime) (*Call, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to event messages
	callChannel, err := bus.Subscribe("call:invoke")
	if err != nil {
		return nil, err
	}

	result := &Call{
		quitChannel: quitChannel,
		callChannel: callChannel,
		logger:      logger.CustomLogger("Call Subsystem"),
		DB:          DB,
		bus:         bus,
		runtime:     runtime,
	}

	return result, nil
}

// Start the subsystem
func (c *Call) Start() error {

	c.running = true

	// Spin off a go routine
	go func() {
		for c.running {
			select {
			case <-c.quitChannel:
				c.running = false
			case callMessage := <-c.callChannel:
				// TODO: Check if this works ok in a goroutine
				c.processCall(callMessage)
			}
		}

		// Call shutdown
		c.shutdown()
	}()

	return nil
}

func (c *Call) processCall(callMessage *servicebus.Message) {

	c.logger.Trace("Got message: %+v", callMessage)

	// Extract payload
	payload := callMessage.Data().(*message.CallMessage)

	// Lookup method
	registeredMethod := c.DB.GetMethod(payload.Name)

	// Check if it's a system call
	if strings.HasPrefix(payload.Name, ".wails.") {
		c.processSystemCall(payload, callMessage.Target())
		return
	}

	// Check we have it
	if registeredMethod == nil {
		c.sendError(fmt.Errorf("Method not registered"), payload, callMessage.Target())
		return
	}
	c.logger.Trace("Got registered method: %+v", registeredMethod)

	result, err := registeredMethod.Call(payload.Args)
	if err != nil {
		c.sendError(err, payload, callMessage.Target())
		return
	}
	c.logger.Trace("registeredMethod.Call: %+v, %+v", result, err)
	// process result
	c.sendResult(result, payload, callMessage.Target())

}

func (c *Call) processSystemCall(payload *message.CallMessage, clientID string) {
	c.logger.Trace("Got internal System call: %+v", payload)
	callName := strings.TrimPrefix(payload.Name, ".wails.")
	switch callName {
	case "IsDarkMode":
		darkModeEnabled := c.runtime.System.IsDarkMode()
		c.sendResult(darkModeEnabled, payload, clientID)
	}
}

func (c *Call) sendResult(result interface{}, payload *message.CallMessage, clientID string) {
	c.logger.Trace("Sending success result with CallbackID '%s' : %+v\n", payload.CallbackID, result)
	message := &CallbackMessage{
		Result:     result,
		CallbackID: payload.CallbackID,
	}
	messageData, err := json.Marshal(message)
	c.logger.Trace("json message data: %+v\n", string(messageData))
	if err != nil {
		// what now?
		c.logger.Fatal(err.Error())
	}
	c.bus.PublishForTarget("call:result", string(messageData), clientID)
}

func (c *Call) sendError(err error, payload *message.CallMessage, clientID string) {
	c.logger.Trace("Sending error result with CallbackID '%s' : %+v\n", payload.CallbackID, err.Error())
	message := &CallbackMessage{
		Err:        err.Error(),
		CallbackID: payload.CallbackID,
	}

	messageData, err := json.Marshal(message)
	c.logger.Trace("json message data: %+v\n", string(messageData))
	if err != nil {
		// what now?
		c.logger.Fatal(err.Error())
	}
	c.bus.PublishForTarget("call:result", string(messageData), clientID)
}

func (c *Call) shutdown() {
	c.logger.Trace("Shutdown")
}

// CallbackMessage defines a message that contains the result of a call
type CallbackMessage struct {
	Result     interface{} `json:"result"`
	Err        string      `json:"error"`
	CallbackID string      `json:"callbackid"`
}
