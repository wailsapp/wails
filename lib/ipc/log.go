package ipc

import "github.com/wailsapp/wails/lib/messages"

// Register the message handler
func init() {
	messageProcessors["log"] = processLogData
}

// This processes the given log message
func processLogData(message *ipcMessage) (*ipcMessage, error) {

	var payload messages.LogData

	// Decode event data
	payloadMap := message.Payload.(map[string]interface{})
	payload.Level = payloadMap["level"].(string)
	payload.Message = payloadMap["message"].(string)

	// Reassign payload to decoded data
	message.Payload = &payload

	return message, nil
}
