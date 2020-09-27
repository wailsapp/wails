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
	var data []string

	// Switch the event type (with or without data)
	switch message[0] {
	// Format of Dialog response messages: D<callbackID>|<[]string as json encoded string>
	case 'D':
		idx := strings.IndexByte(message[1:], '|')
		if idx < 0 {
			return nil, fmt.Errorf("Invalid dialog response message format")
		}
		callbackID := message[1 : idx+1]
		jsonData := message[idx+2:]
		topic = "dialog:openselected:" + callbackID

		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid message to dialogMessageParser()")
	}

	// Create a new parsed message struct
	parsedMessage := &parsedMessage{Topic: topic, Data: data}

	return parsedMessage, nil
}
