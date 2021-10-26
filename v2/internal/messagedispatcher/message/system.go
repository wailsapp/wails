package message

import (
	"fmt"
	"strings"
)

// systemMessageParser does what it says on the tin!
func systemMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: system messages must be at least 2 bytes
	if len(message) < 2 {
		return nil, fmt.Errorf("system message was an invalid length")
	}

	var responseMessage *parsedMessage

	// Remove 'S'
	message = message[1:]

	// Switch the event type (with or without data)
	switch message[0] {
	// Format of system response messages: S<command><callbackID>|<payload>
	// DarkModeEnabled
	case 'D':
		if len(message) < 4 {
			return nil, fmt.Errorf("system message was an invalid length")
		}
		message = message[1:]
		idx := strings.IndexByte(message, '|')
		if idx < 0 {
			return nil, fmt.Errorf("Invalid system response message format")
		}
		callbackID := message[:idx]
		payloadData := message[idx+1:]

		topic := "systemresponse:" + callbackID
		responseMessage = &parsedMessage{Topic: topic, Data: payloadData == "T"}

		// This is our startup hook - the frontend is now ready
	case 'S':
		topic := "hooks:startup"
		startupURL := message[1:]
		responseMessage = &parsedMessage{Topic: topic, Data: startupURL}
	default:
		return nil, fmt.Errorf("Invalid message to systemMessageParser()")
	}

	return responseMessage, nil
}
