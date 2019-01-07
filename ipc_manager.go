package wails

import (
	"fmt"
)

type ipcManager struct {
	renderer     Renderer // The renderer
	messageQueue chan *ipcMessage
	// quitChannel  chan struct{}
	// signals      chan os.Signal
	log            *CustomLogger
	eventManager   *eventManager
	bindingManager *bindingManager
}

func newIPCManager() *ipcManager {
	result := &ipcManager{
		messageQueue: make(chan *ipcMessage, 100),
		// 		quitChannel:  make(chan struct{}),
		// 		signals:      make(chan os.Signal, 1),
		log: newCustomLogger("IPC"),
	}
	return result
}

// Sets the renderer, returns the dispatch function
func (i *ipcManager) bindRenderer(renderer Renderer) {
	i.renderer = renderer
}

func (i *ipcManager) start(eventManager *eventManager, bindingManager *bindingManager) {

	// Store manager references
	i.eventManager = eventManager
	i.bindingManager = bindingManager

	i.log.Info("Starting")
	// signal.Notify(manager.signals, os.Interrupt)
	go func() {
		running := true
		for running {
			select {
			case incomingMessage := <-i.messageQueue:
				i.log.DebugFields("Processing message", Fields{
					"1D": &incomingMessage,
				})
				switch incomingMessage.Type {
				case "call":
					callData := incomingMessage.Payload.(*callData)
					i.log.DebugFields("Processing call", Fields{
						"1D":          &incomingMessage,
						"bindingName": callData.BindingName,
						"data":        callData.Data,
					})
					go func() {
						result, err := bindingManager.processCall(callData)
						i.log.DebugFields("processed call", Fields{"result": result, "err": err})
						if err != nil {
							incomingMessage.ReturnError(err.Error())
						} else {
							incomingMessage.ReturnSuccess(result)
						}
						i.log.DebugFields("Finished processing call", Fields{
							"1D": &incomingMessage,
						})
					}()
				case "event":

					// Extract event data
					eventData := incomingMessage.Payload.(*eventData)

					// Log
					i.log.DebugFields("Processing event", Fields{
						"name": eventData.Name,
						"data": eventData.Data,
					})

					// Push the event to the event manager
					i.eventManager.PushEvent(eventData)

					// Log
					i.log.DebugFields("Finished processing event", Fields{
						"name": eventData.Name,
					})
				case "log":
					logdata := incomingMessage.Payload.(*logData)
					switch logdata.Level {
					case "info":
						logger.Info(logdata.Message)
					case "debug":
						logger.Debug(logdata.Message)
					case "warning":
						logger.Warning(logdata.Message)
					case "error":
						logger.Error(logdata.Message)
					case "fatal":
						logger.Fatal(logdata.Message)
					default:
						i.log.ErrorFields("Invalid log level sent", Fields{
							"level":   logdata.Level,
							"message": logdata.Message,
						})
					}
				default:
					i.log.Debugf("bad message sent to MessageQueue! Unknown type: %s", incomingMessage.Type)
				}

				// Log
				i.log.DebugFields("Finished processing message", Fields{
					"1D": &incomingMessage,
				})
				// 			case <-manager.quitChannel:
				// 				Debug("[MessageQueue] Quit caught")
				// 				running = false
				// 			case <-manager.signals:
				// 				Debug("[MessageQueue] Signal caught")
				// 				running = false
			}
		}
		i.log.Debug("Stopping")
	}()
}

// Dispatch receives JSON encoded messages from the renderer.
// It processes the message to ensure that it is valid and places
// the processed message on the message queue
func (i *ipcManager) Dispatch(message string) {

	// Create a new IPC Message
	incomingMessage, err := newIPCMessage(message, i.SendResponse)
	if err != nil {
		i.log.ErrorFields("Could not understand incoming message! ", map[string]interface{}{
			"message": message,
			"error":   err,
		})
		return
	}

	// Put message on queue
	i.log.DebugFields("Message received", map[string]interface{}{
		"type":    incomingMessage.Type,
		"payload": incomingMessage.Payload,
	})

	// Put incoming message on the message queue
	i.messageQueue <- incomingMessage
}

// SendResponse sends the given response back to the frontend
func (i *ipcManager) SendResponse(response *ipcResponse) error {

	// Serialise the Message
	data, err := response.Serialise()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	// Call back to the front end
	return i.renderer.Callback(data)
}
