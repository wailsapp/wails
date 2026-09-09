package application

import (
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func windowMethodRequest(t *testing.T, method int, args string) *RuntimeRequest {
	t.Helper()
	parsedArgs := &Args{}
	if err := json.Unmarshal([]byte(args), parsedArgs); err != nil {
		t.Fatal(err)
	}
	return &RuntimeRequest{Method: method, Args: parsedArgs}
}

func TestProcessWindowMethodSetFullscreenButtonEnabled(t *testing.T) {
	m := NewMessageProcessor(slog.New(slog.NewTextHandler(io.Discard, nil)))
	window := &WebviewWindow{}

	result, err := m.processWindowMethod(windowMethodRequest(t, WindowSetFullscreenButtonEnabled, `{"enabled":false}`), window)
	if err != nil {
		t.Fatalf("SetFullscreenButtonEnabled(false) returned error: %v", err)
	}
	if result != unit {
		t.Errorf("SetFullscreenButtonEnabled(false) result = %v, want unit", result)
	}
	if window.options.FullscreenButtonState != ButtonDisabled {
		t.Errorf("FullscreenButtonState = %v, want ButtonDisabled", window.options.FullscreenButtonState)
	}

	_, err = m.processWindowMethod(windowMethodRequest(t, WindowSetFullscreenButtonEnabled, `{"enabled":true}`), window)
	if err != nil {
		t.Fatalf("SetFullscreenButtonEnabled(true) returned error: %v", err)
	}
	if window.options.FullscreenButtonState != ButtonEnabled {
		t.Errorf("FullscreenButtonState = %v, want ButtonEnabled", window.options.FullscreenButtonState)
	}
}

func TestProcessWindowMethodSetFullscreenButtonEnabledInvalidArgs(t *testing.T) {
	m := NewMessageProcessor(slog.New(slog.NewTextHandler(io.Discard, nil)))
	window := &WebviewWindow{}

	tests := []struct {
		name string
		args string
	}{
		{name: "missing argument", args: `{}`},
		{name: "wrong argument type", args: `{"enabled":"yes"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.processWindowMethod(windowMethodRequest(t, WindowSetFullscreenButtonEnabled, tt.args), window)
			if err == nil {
				t.Fatal("expected an error")
			}
			if strings.Contains(err.Error(), "Unknown method") {
				t.Errorf("method was not dispatched: %v", err)
			}
		})
	}
}
