package application

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

// The Android transport: the WebView cannot deliver fetch() POST bodies to
// shouldInterceptRequest, so the JavaScript runtime routes runtime calls
// through the JavascriptInterface bridge to nativeHandleRuntimeCall.
// Keep envelope handling independent of JNI so it can be tested on the host.

var (
	androidMessageProcessor     *MessageProcessor
	androidMessageProcessorOnce sync.Once
)

func androidProcessor(app *App) *MessageProcessor {
	androidMessageProcessorOnce.Do(func() {
		androidMessageProcessor = NewMessageProcessor(app.Logger)
	})
	return androidMessageProcessor
}

type androidRuntimeCallPayload struct {
	Object     *int            `json:"object"`
	Method     *int            `json:"method"`
	WindowName string          `json:"windowName"`
	Args       json.RawMessage `json:"args"`
	ClientID   string          `json:"clientId"`
}

// handleRuntimeCallForAndroid processes a runtime call payload from the
// JavaScript Android transport and returns a response envelope:
// {"ok":true,"data":...} / {"ok":true,"text":"..."} / {"ok":false,"error":...}.
// Binding errors retain their CallError fields; other errors remain strings.
func handleRuntimeCallForAndroid(app *App, payload string) string {
	fail := func(err error) string {
		var value any = err.Error()
		var callErr *CallError
		if errors.As(err, &callErr) {
			value = callErr
		}
		b, marshalErr := json.Marshal(map[string]any{"ok": false, "error": value})
		if marshalErr != nil {
			b, _ = json.Marshal(map[string]any{"ok": false, "error": err.Error()})
		}
		return string(b)
	}

	var call androidRuntimeCallPayload
	if err := json.Unmarshal([]byte(payload), &call); err != nil {
		return fail(errors.New("unable to parse runtime call: " + err.Error()))
	}
	if call.Object == nil {
		return fail(errors.New("missing object value"))
	}
	if call.Method == nil {
		return fail(errors.New("missing method value"))
	}

	resp, err := androidProcessor(app).HandleRuntimeCallWithIDs(context.Background(), &RuntimeRequest{
		Object:            *call.Object,
		Method:            *call.Method,
		Args:              &Args{call.Args},
		WebviewWindowName: call.WindowName,
		ClientID:          call.ClientID,
	})
	if err != nil {
		return fail(err)
	}

	var envelope map[string]any
	if text, ok := resp.(string); ok {
		envelope = map[string]any{"ok": true, "text": text}
	} else {
		envelope = map[string]any{"ok": true, "data": resp}
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return fail(errors.New("unable to marshal response: " + err.Error()))
	}
	return string(b)
}
