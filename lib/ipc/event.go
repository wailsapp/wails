package ipc

import (
	"encoding/json"

	"github.com/wailsapp/wails/lib/messages"
)

// Register the message handler
func init() {
	messageProcessors["event"] = processEventData
}

// This processes the given event message
func processEventData(message *ipcMessage) (*ipcMessage, error) {

	// TODO: Is it worth double checking this is actually an event message,
	// even though that's done by the caller?
	var payload messages.EventData

	// Decode event data
	payloadMap := message.Payload.(map[string]interface{})
	payload.Name = payloadMap["name"].(string)

	// decode the payload data
	var data []interface{}
	err := json.Unmarshal([]byte(payloadMap["data"].(string)), &data)
	if err != nil {
		return nil, err
	}
	payload.Data = data

	// Reassign payload to decoded data
	message.Payload = &payload

	return message, nil
}
