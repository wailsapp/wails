package application

import (
	"fmt"
	"strings"
)

type MessageProcessor struct {
	window *WebviewWindow
}

func NewMessageProcessor(w *WebviewWindow) *MessageProcessor {
	return &MessageProcessor{
		window: w,
	}
}

func (m *MessageProcessor) ProcessMessage(message string) {

	// TODO: Implement calls to other windows
	// Check for prefix "WINDOWID"
	// If prefix exists, get window ID by parsing: "WINDOWID:12:MESSAGE"

	if strings.HasPrefix(message, "WINDOWID") {
		m.Error("Window ID prefix not yet implemented")
		return
	}

	window := m.window

	if message == "" {
		m.Error("Blank message received")
		return
	}
	m.Info("Processing message: %s", message)
	switch message[0] {
	//case 'L':
	//	m.processLogMessage(message)
	//case 'E':
	//	return m.processEventMessage(message)
	//case 'C':
	//	return m.processCallMessage(message)
	//case 'c':
	//	return m.processSecureCallMessage(message)
	case 'W':
		m.processWindowMessage(message, window)
	//case 'B':
	//	return m.processBrowserMessage(message)
	case 'Q':
		globalApplication.Quit()
	case 'S':
		//globalApplication.Show()
	case 'H':
		//globalApplication.Hide()
	default:
		m.Error("Unknown message from front end:", message)
	}
}

func (m *MessageProcessor) Error(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Error: "+message, args...)
}

func (m *MessageProcessor) Info(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Info: "+message, args...)
}
