package dispatcher

import (
	"encoding/json"
	"errors"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

type EventMessage struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

func (d *Dispatcher) processEventMessage(message string, sender frontend.Frontend) (string, error) {
	if len(message) < 3 {
		return "", errors.New("Invalid Event Message: " + message)
	}

	switch message[1] {
	case 'E':
		var eventMessage EventMessage
		err := json.Unmarshal([]byte(message[2:]), &eventMessage)
		if err != nil {
			return "", err
		}
		go d.events.Notify(sender, eventMessage.Name, eventMessage.Data...)
	case 'X':
		eventName := message[2:]
		go d.events.Off(eventName)
	}

	return "", nil
}
