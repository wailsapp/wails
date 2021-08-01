package dispatcher

import (
	"encoding/json"
	"errors"
)

type EventMessage struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

func (d *Dispatcher) processEventMessage(message string) error {
	if len(message) < 3 {
		return errors.New("Invalid Event Message: " + message)
	}

	switch message[1] {
	case 'E':
		var eventMessage EventMessage
		err := json.Unmarshal([]byte(message[2:]), &eventMessage)
		if err != nil {
			return err
		}
		go d.events.Notify(eventMessage.Name, eventMessage.Data)
	case 'X':
		eventName := message[2:]
		go d.events.Off(eventName)
	}

	return nil
}
