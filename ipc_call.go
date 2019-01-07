package wails

import (
	"fmt"
)

type callData struct {
	BindingName string `json:"bindingName"`
	Data        string `json:"data,omitempty"`
}

func init() {
	messageProcessors["call"] = processCallData
}

func processCallData(message *ipcMessage) (*ipcMessage, error) {

	var payload callData

	// Decode binding call data
	payloadMap := message.Payload.(map[string]interface{})

	// Check for binding name
	if payloadMap["bindingName"] == nil {
		return nil, fmt.Errorf("bindingName not given in call")
	}
	payload.BindingName = payloadMap["bindingName"].(string)

	// Check for data
	if payloadMap["data"] != nil {
		payload.Data = payloadMap["data"].(string)
	}

	// Reassign payload to decoded data
	message.Payload = &payload

	return message, nil
}
