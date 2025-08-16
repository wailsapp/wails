package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type contextKey string

const (
	CallBinding            = 0
	WindowKey   contextKey = "Window"
)

func (m *MessageProcessor) callErrorCallback(window Window, message string, callID *string, err error) {
	m.Error(message, "id", *callID, "error", err)
	if cerr := (*CallError)(nil); errors.As(err, &cerr) {
		if data, jsonErr := json.Marshal(cerr); jsonErr == nil {
			window.CallError(*callID, string(data), true)
			return
		} else {
			m.Error("Unable to convert data to JSON. Please report this to the Wails team!", "id", *callID, "error", jsonErr)
		}
	}

	window.CallError(*callID, err.Error(), false)
}

func (m *MessageProcessor) callCallback(window Window, callID *string, result string) {
	window.CallResponse(*callID, result)
}

func (m *MessageProcessor) processCallCancelMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	var cancel func()
	func() {
		m.l.Lock()
		defer m.l.Unlock()
		cancel = m.runningCalls[*callID]
	}()

	if cancel != nil {
		cancel()
		m.Info("Binding call canceled:", "id", *callID)
	}
	m.ok(rw)
}

func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("error parsing call options: %w", err))
			return
		}

		// Set CORS headers for the response
		m.setCORSHeaders(rw, r)

		// Create context with timeout from configuration
		timeout := globalApplication.getBindingTimeout()
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Track the call for cancellation support
		func() {
			m.l.Lock()
			defer m.l.Unlock()
			
			if m.runningCalls[*callID] != nil {
				m.handleBindingError(rw, fmt.Errorf("ambiguous call id: %s", *callID))
				return
			}
			m.runningCalls[*callID] = cancel
		}()

		// Clean up running call tracking
		defer func() {
			m.l.Lock()
			defer m.l.Unlock()
			delete(m.runningCalls, *callID)
		}()

		// Log call
		var methodRef any = options.MethodName
		if options.MethodName == "" {
			methodRef = options.MethodID
		}
		m.Info("Binding call started:", "id", *callID, "method", methodRef)

		// Find the bound method
		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.handleBindingError(rw, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
				})
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				m.handleBindingError(rw, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method id %d", options.MethodID),
				})
				return
			}
		}

		// Prepare args for logging
		jsonArgs, _ := json.Marshal(options.Args)

		// Set the context values for the window
		if window != nil {
			ctx = context.WithValue(ctx, WindowKey, window)
		}

		// Execute the binding method synchronously
		result, err := boundMethod.Call(ctx, options.Args)
		
		// Log the completion
		var jsonResult []byte
		if result != nil {
			jsonResult, _ = json.Marshal(result)
		}
		m.Info("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))

		// Handle errors
		if err != nil {
			if cerr := (*CallError)(nil); errors.As(err, &cerr) {
				m.handleBindingError(rw, cerr)
			} else {
				m.handleBindingError(rw, err)
			}
			return
		}

		// Return successful result directly via HTTP
		m.json(rw, result)

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
