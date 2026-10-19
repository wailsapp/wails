package application

import "testing"

// messageRecorderWindow records the messages handed to HandleMessage.
// Embedding the Window interface gives nil-bodied placeholders for every
// other method; handleWindowMessage only calls HandleMessage.
type messageRecorderWindow struct {
	Window
	handled []string
}

func (w *messageRecorderWindow) HandleMessage(message string) {
	w.handled = append(w.handled, message)
}

type rawMessage struct {
	message    string
	originInfo *OriginInfo
}

func newWindowMessageTestApp() (*App, *messageRecorderWindow, *[]rawMessage) {
	window := &messageRecorderWindow{}
	raw := &[]rawMessage{}
	app := &App{
		windows: map[uint]Window{1: window},
		options: Options{
			RawMessageHandler: func(_ Window, message string, originInfo *OriginInfo) {
				*raw = append(*raw, rawMessage{message: message, originInfo: originInfo})
			},
		},
	}
	return app, window, raw
}

func TestHandleWindowMessageIgnoresReservedMessagesFromChildFrames(t *testing.T) {
	app, window, raw := newWindowMessageTestApp()

	for _, message := range []string{"wails:runtime:ready", "wails:drag", "wails:resize:se"} {
		app.handleWindowMessage(&windowMessage{
			windowId:   1,
			message:    message,
			originInfo: &OriginInfo{IsMainFrame: false},
			childFrame: true,
		})
	}

	if len(window.handled) != 0 {
		t.Fatalf("reserved messages from a child frame reached the window: %v", window.handled)
	}
	if len(*raw) != 0 {
		t.Fatalf("reserved messages from a child frame reached RawMessageHandler: %v", *raw)
	}
}

func TestHandleWindowMessageForwardsChildFrameRawMessages(t *testing.T) {
	app, window, raw := newWindowMessageTestApp()

	origin := &OriginInfo{Origin: "https://example.com", IsMainFrame: false}
	app.handleWindowMessage(&windowMessage{
		windowId:   1,
		message:    "hello",
		originInfo: origin,
		childFrame: true,
	})

	if len(window.handled) != 0 {
		t.Fatalf("raw message reached the window: %v", window.handled)
	}
	if len(*raw) != 1 || (*raw)[0].message != "hello" || (*raw)[0].originInfo != origin {
		t.Fatalf("expected raw message with origin info to reach RawMessageHandler, got %v", *raw)
	}
}

func TestHandleWindowMessageDispatchesReservedMessagesFromMainFrame(t *testing.T) {
	app, window, raw := newWindowMessageTestApp()

	// Main frame on a platform that reports frame identity
	app.handleWindowMessage(&windowMessage{
		windowId:   1,
		message:    "wails:runtime:ready",
		originInfo: &OriginInfo{IsMainFrame: true},
	})
	// Platform that cannot identify the posting frame
	app.handleWindowMessage(&windowMessage{
		windowId: 1,
		message:  "wails:drag",
	})

	if len(window.handled) != 2 || window.handled[0] != "wails:runtime:ready" || window.handled[1] != "wails:drag" {
		t.Fatalf("expected reserved messages to reach the window, got %v", window.handled)
	}
	if len(*raw) != 0 {
		t.Fatalf("reserved messages reached RawMessageHandler: %v", *raw)
	}
}
