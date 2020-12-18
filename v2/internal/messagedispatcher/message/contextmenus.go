package message

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

// ContextMenusOnMessage is used to emit listener registration requests
// on the service bus
type ContextMenusOnMessage struct {
	// MenuID is the id of the menu item we are interested in
	MenuID string
	// Callback is called when the menu is clicked
	Callback func(*menu.MenuItem, string)
}

// contextMenusMessageParser does what it says on the tin!
func contextMenusMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Menu messages must be at least 2 bytes
	if len(message) < 3 {
		return nil, fmt.Errorf("context menus message was an invalid length")
	}

	var topic string
	var data interface{}

	// Switch the message type
	switch message[1] {
	case 'C':
		contextMenuData := message[2:]
		topic = "contextmenus:clicked"
		data = contextMenuData
	default:
		return nil, fmt.Errorf("invalid menu message: %s", message)
	}

	// Create a new parsed message struct
	parsedMessage := &parsedMessage{Topic: topic, Data: data}

	return parsedMessage, nil
}
