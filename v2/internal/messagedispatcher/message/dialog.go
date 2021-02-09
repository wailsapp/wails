package message

import (
	"encoding/json"
	"fmt"
	"strings"
)

// dialogMessageParser does what it says on the tin!
func dialogMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Dialog messages must be at least 4 bytes
	if len(message) < 4 {
		return nil, fmt.Errorf("dialog message was an invalid length")
	}

	var topic = "bad topic from dialogMessageParser"
	var responseMessage *parsedMessage

	// Switch the event type (with or without data)
	switch message[0] {
	// Format of Dialog response messages: D<dialog type><callbackID>|<[]string as json encoded string>
	case 'D':
		dialogType := message[1]
		message = message[2:]
		idx := strings.IndexByte(message, '|')
		if idx < 0 {
			return nil, fmt.Errorf("Invalid dialog response message format: %+v", message)
		}
		callbackID := message[:idx]
		payloadData := message[idx+1:]

		switch dialogType {
		case 'O':
			var data []string
			topic = "dialog:openselected:" + callbackID
			err := json.Unmarshal([]byte(payloadData), &data)
			if err != nil {
				return nil, err
			}
			responseMessage = &parsedMessage{Topic: topic, Data: data}
		case 'S':
			topic = "dialog:saveselected:" + callbackID
			responseMessage = &parsedMessage{Topic: topic, Data: payloadData}
		case 'M':
			topic = "dialog:messageselected:" + callbackID
			responseMessage = &parsedMessage{Topic: topic, Data: payloadData}
		}

	default:
		return nil, fmt.Errorf("Invalid message to dialogMessageParser()")
	}

	return responseMessage, nil
}
