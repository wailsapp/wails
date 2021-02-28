package message

import "fmt"

// urlMessageParser does what it says on the tin!
func urlMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: URL messages must be at least 2 bytes
	if len(message) < 2 {
		return nil, fmt.Errorf("log message was an invalid length")
	}

	// Switch on the log type
	switch message[1] {
	case 'C':
		return &parsedMessage{Topic: "url:handler", Data: message[2:]}, nil
	default:
		return nil, fmt.Errorf("url message type '%c' invalid", message[1])
	}
}
