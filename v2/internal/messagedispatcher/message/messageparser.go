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
	'C': callMessageParser,
	'W': windowMessageParser,
	'D': dialogMessageParser,
	'S': systemMessageParser,
	'M': menuMessageParser,
	'T': trayMessageParser,
	'X': contextMenusMessageParser,
	'U': urlMessageParser,
}

// Parse will attempt to parse the given message
func Parse(message string) (*parsedMessage, error) {

	if len(message) == 0 {
		return nil, fmt.Errorf("MessageParser received blank message")
	}

	parseMethod := messageParsers[message[0]]
	if parseMethod == nil {
		return nil, fmt.Errorf("message type '%c' invalid", message[0])
	}

	return parseMethod(message)
}
