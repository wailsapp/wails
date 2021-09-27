package message

import (
	"encoding/json"
	"fmt"
)

type EventMessage struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

type OnEventMessage struct {
	Name     string
	Callback func(optionalData ...interface{})
	Counter  int
}

// eventMessageParser does what it says on the tin!
func eventMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Event messages must be at least 2 bytes
	if len(message) < 3 {
		return nil, fmt.Errorf("event message was an invalid length")
	}

	eventMessage := new(EventMessage)
	direction := message[1]

	// Switch the event type (with or without data)
	switch message[0] {
	case 'E':
		m := message[2:]
		err := json.Unmarshal([]byte(m), eventMessage)
		if err != nil {
			println(err.Error())
			return nil, err
		}
	}

	topic := "event:emit:from:" + string(direction)

	// Create a new parsed message struct
	parsedMessage := &parsedMessage{Topic: topic, Data: eventMessage}

	return parsedMessage, nil
}
