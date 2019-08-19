package ipc

import (
	"encoding/json"
	"fmt"
)

// Message handler
type messageProcessorFunc func(*ipcMessage) (*ipcMessage, error)

var messageProcessors = make(map[string]messageProcessorFunc)

// ipcMessage is the struct version of the Message sent from the frontend.
// The payload has the specialised message data
type ipcMessage struct {
	Type         string      `json:"type"`
	Payload      interface{} `json:"payload"`
	CallbackID   string      `json:"callbackid,omitempty"`
	sendResponse func(*ipcResponse) error
}

func parseMessage(incomingMessage string) (*ipcMessage, error) {
	// Parse message
	var message ipcMessage
	err := json.Unmarshal([]byte(incomingMessage), &message)
	return &message, err
}

func newIPCMessage(incomingMessage string, responseFunction func(*ipcResponse) error) (*ipcMessage, error) {

	// Parse the Message
	message, err := parseMessage(incomingMessage)
	if err != nil {
		return nil, err
	}

	// Check message type is valid
	messageProcessor := messageProcessors[message.Type]
	if messageProcessor == nil {
		return nil, fmt.Errorf("unknown message type: %s", message.Type)
	}

	// Process message payload
	message, err = messageProcessor(message)
	if err != nil {
		return nil, err
	}

	// Set the response function
	message.sendResponse = responseFunction

	return message, nil
}

// hasCallbackID checks if the message can send an error back to the frontend
func (m *ipcMessage) hasCallbackID() error {
	if m.CallbackID == "" {
		return fmt.Errorf("attempted to return error to message with no Callback ID")
	}
	return nil
}

// ReturnError returns an error back to the frontend
func (m *ipcMessage) ReturnError(format string, args ...interface{}) error {

	// Ignore ReturnError if no callback ID given
	err := m.hasCallbackID()
	if err != nil {
		return err
	}

	// Create response
	response := newErrorResponse(m.CallbackID, fmt.Sprintf(format, args...))

	// Send response
	return m.sendResponse(response)
}

// ReturnSuccess returns a success message back with the given data
func (m *ipcMessage) ReturnSuccess(data interface{}) error {

	// Ignore ReturnSuccess if no callback ID given
	err := m.hasCallbackID()
	if err != nil {
		return err
	}

	// Create the response
	response := newSuccessResponse(m.CallbackID, data)

	// Send response
	return m.sendResponse(response)
}
