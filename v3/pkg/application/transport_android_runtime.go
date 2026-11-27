package application

import (
	"context"
	"encoding/json"
)

type androidRuntimeCallPayload struct {
	Object     *int            `json:"object"`
	Method     *int            `json:"method"`
	WindowName string          `json:"windowName"`
	Args       json.RawMessage `json:"args"`
	ClientID   string          `json:"clientId"`
}

// prepareAndroidRuntimeCall runs on the JavascriptInterface thread. It registers
// bindings and handles cancellations in submission order before their workers
// can overtake each other. The returned function produces the response.
func (m *MessageProcessor) prepareAndroidRuntimeCall(payload string) func() string {
	fail := func(msg string) func() string {
		return func() string { return androidRuntimeError(msg) }
	}
	var call androidRuntimeCallPayload
	if err := json.Unmarshal([]byte(payload), &call); err != nil {
		return fail("unable to parse runtime call: " + err.Error())
	}
	if call.Object == nil {
		return fail("missing object value")
	}
	if call.Method == nil {
		return fail("missing method value")
	}

	req := &RuntimeRequest{
		Object:            *call.Object,
		Method:            *call.Method,
		Args:              &Args{call.Args},
		WebviewWindowName: call.WindowName,
		ClientID:          call.ClientID,
	}
	ctx := context.Background()
	if req.Object != callRequest && req.Object != cancelCallRequest {
		return func() string {
			resp, err := m.HandleRuntimeCallWithIDs(ctx, req)
			return androidRuntimeResponse(resp, err)
		}
	}

	window, err := m.resolveRuntimeWindow(req)
	if err != nil {
		return fail(err.Error())
	}
	// clientID identifies one JS runtime instance. Use the supplied window
	// name so the default target matches on both sides of the handoff.
	scope := bindingCallScope{clientID: req.ClientID, windowName: req.WebviewWindowName}
	if req.Object == cancelCallRequest {
		m.logRuntimeCall(req)
		resp, err := m.processCallCancelMethod(req, scope)
		return func() string { return androidRuntimeResponse(resp, err) }
	}
	run, err := m.prepareCallMethod(ctx, req, window, scope)
	if err != nil {
		return fail(err.Error())
	}
	return func() string {
		resp, err := run()
		return androidRuntimeResponse(resp, err)
	}
}

func androidRuntimeResponse(resp any, err error) string {
	if err != nil {
		return androidRuntimeError(err.Error())
	}
	var envelope map[string]any
	if text, ok := resp.(string); ok {
		envelope = map[string]any{"ok": true, "text": text}
	} else {
		envelope = map[string]any{"ok": true, "data": resp}
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return androidRuntimeError("unable to marshal response: " + err.Error())
	}
	return string(b)
}

func androidRuntimeError(message string) string {
	b, _ := json.Marshal(map[string]any{"ok": false, "error": message})
	return string(b)
}
