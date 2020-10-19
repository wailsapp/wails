package message

import "fmt"

// runtimeMessageParser does what it says on the tin!
func runtimeMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Log messages must be at least 2 bytes
	if len(message) < 3 {
		return nil, fmt.Errorf("runtime message was an invalid length")
	}

	// Switch on the runtime module type
	module := message[1]
	switch module {
	case 'B':
		return processBrowserMessage(message)
	}

	return nil, fmt.Errorf("unknown message: %s", message)
}

// processBrowserMessage expects messages of the following format:
// RB<METHOD><DATA>
// O = Open
func processBrowserMessage(message string) (*parsedMessage, error) {
	method := message[2]
	switch method {
	case 'O':
		// Open URL
		target := message[3:]
		return &parsedMessage{Topic: "runtime:browser:open", Data: target}, nil
	}

	return nil, fmt.Errorf("unknown browser message: %s", message)

}
