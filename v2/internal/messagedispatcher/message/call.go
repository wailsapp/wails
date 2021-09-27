package message

import (
	"encoding/json"
	"fmt"
)

type CallMessage struct {
	Name       string            `json:"name"`
	Args       []json.RawMessage `json:"args"`
	CallbackID string            `json:"callbackID,omitempty"`
}

// callMessageParser does what it says on the tin!
func callMessageParser(message string) (*parsedMessage, error) {

	// Sanity check: Call messages must be at least 3 bytes `C{}``
	if len(message) < 3 {
		return nil, fmt.Errorf("call message was an invalid length")
	}

	callMessage := new(CallMessage)

	m := message[1:]

	err := json.Unmarshal([]byte(m), callMessage)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	topic := "call:invoke"

	// Create a new parsed message struct
	parsedMessage := &parsedMessage{Topic: topic, Data: callMessage}

	return parsedMessage, nil
}
