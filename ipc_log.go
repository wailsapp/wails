package wails

type logData struct {
	Level   string `json:"level"`
	Message string `json:"string"`
}

// Register the message handler
func init() {
	messageProcessors["log"] = processLogData
}

// This processes the given log message
func processLogData(message *ipcMessage) (*ipcMessage, error) {

	var payload logData

	// Decode event data
	payloadMap := message.Payload.(map[string]interface{})
	payload.Level = payloadMap["level"].(string)
	payload.Message = payloadMap["message"].(string)

	// Reassign payload to decoded data
	message.Payload = &payload

	return message, nil
}
