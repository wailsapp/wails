package message

import "fmt"

// Parse
type parsedMessage struct {
	Topic    string
	ClientID string
	Data     interface{}
}

// Map of different message parsers based on the header byte of the message
var messageParsers = map[byte]func(string) (*parsedMessage, error){
	'L': logMessageParser,
	'R': runtimeMessageParser,
	'E': eventMessageParser,
	'e': eventMessageParser,
	'C': callMessageParser,
	'W': windowMessageParser,
}

// Parse will attempt to parse the given message
func Parse(message string) (*parsedMessage, error) {

	parseMethod := messageParsers[message[0]]
	if parseMethod == nil {
		return nil, fmt.Errorf("message type '%b' invalid", message[0])
	}

	return parseMethod(message)
}
