package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type contextKey string

const (
	CallBinding              = 0
	WindowNameKey contextKey = "WindowName"
	WindowIDKey   contextKey = "WindowID"
)

func (m *MessageProcessor) callErrorCallback(window Window, message string, callID *string, err error) {
	errorMsg := fmt.Sprintf(message, err)
	m.Error(errorMsg)
	window.CallError(*callID, errorMsg)
}

func (m *MessageProcessor) callCallback(window Window, callID *string, result string, isJSON bool) {
	window.CallResponse(*callID, result)
}

func (m *MessageProcessor) processCallCancelMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}
	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.Error("call-id is required")
		return
	}

	m.l.Lock()
	cancel := m.runningCalls[*callID]
	m.l.Unlock()

	if cancel != nil {
		cancel()
	}
	m.ok(rw)
}

func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}
	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.Error("call-id is required")
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.callErrorCallback(window, "Error parsing call options: %s", callID, err)
			return
		}
		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method '%s' not found", options.MethodName))
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
		}
		if boundMethod == nil {
			m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method ID %d not found", options.MethodID))
			return
		}

		ctx, cancel := context.WithCancel(context.WithoutCancel(r.Context()))

		ambiguousID := false
		m.l.Lock()
		if m.runningCalls[*callID] != nil {
			ambiguousID = true
		} else {
			m.runningCalls[*callID] = cancel
		}
		m.l.Unlock()

		if ambiguousID {
			cancel()
			m.callErrorCallback(window, "Error calling method: %s, a method call with the same id is already running", callID, err)
			return
		}

		// Set the context values for the window
		if window != nil {
			ctx = context.WithValue(ctx, WindowNameKey, window.Name())
			ctx = context.WithValue(ctx, WindowIDKey, window.ID())
		}

		go func() {
			defer func() {
				cancel()

				m.l.Lock()
				delete(m.runningCalls, *callID)
				m.l.Unlock()
			}()

			result, err := boundMethod.Call(ctx, options.Args)
			if err != nil {
				m.callErrorCallback(window, "Error calling method: %s", callID, err)
				return
			}
			var jsonResult = []byte("{}")
			if result != nil {
				// convert result to json
				jsonResult, err = json.Marshal(result)
				if err != nil {
					m.callErrorCallback(window, "Error converting result to json: %s", callID, err)
					return
				}
			}
			m.callCallback(window, callID, string(jsonResult), true)

			var jsonArgs struct {
				Args json.RawMessage `json:"args"`
			}
			err = params.ToStruct(&jsonArgs)
			if err != nil {
				m.callErrorCallback(window, "Error parsing arguments: %s", callID, err)
				return
			}
			m.Info("Call Binding:", "method", boundMethod, "args", string(jsonArgs.Args), "result", result)
		}()
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown call method: %d", method)
	}

}
