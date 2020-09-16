package message

import "fmt"

// windowMessageParser does what it says on the tin!
func windowMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Window messages must be at least 2 bytes
	if len(message) < 2 {
		return nil, fmt.Errorf("window message was an invalid length")
	}

	// Extract event type
	windowEvent := message[1]
	parsedMessage := &parsedMessage{}

	// Switch the windowEvent type
	switch windowEvent {

	// Closed window
	case 'C':
		parsedMessage.Topic = "quit"
		parsedMessage.Data = "Window Closed"

	// Center window
	case 'c':
		parsedMessage.Topic = "window:center"
		parsedMessage.Data = ""

	// Hide window
	case 'H':
		parsedMessage.Topic = "window:hide"
		parsedMessage.Data = ""

	// Show window
	case 'S':
		parsedMessage.Topic = "window:show"
		parsedMessage.Data = ""

	// Position window
	case 'p':
		parsedMessage.Topic = "window:position:" + message[3:]
		parsedMessage.Data = ""

	// Set window size
	case 's':
		parsedMessage.Topic = "window:size:" + message[3:]
		parsedMessage.Data = ""

	// Maximise window
	case 'M':
		parsedMessage.Topic = "window:maximise"
		parsedMessage.Data = ""

	// Unmaximise window
	case 'U':
		parsedMessage.Topic = "window:unmaximise"
		parsedMessage.Data = ""

	// Minimise window
	case 'm':
		parsedMessage.Topic = "window:minimise"
		parsedMessage.Data = ""

	// Unminimise window
	case 'u':
		parsedMessage.Topic = "window:unminimise"
		parsedMessage.Data = ""

	// Unknown event type
	default:
		return nil, fmt.Errorf("unknown message: %s", message)
	}

	return parsedMessage, nil
}
