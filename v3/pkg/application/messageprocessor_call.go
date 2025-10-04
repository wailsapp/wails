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

// jsonCallError writes a CallError to the HTTP response in JSON format
func (m *MessageProcessor) jsonCallError(rw http.ResponseWriter, err error) {
	var cerr *CallError
	if !errors.As(err, &cerr) {
		cerr = &CallError{
			Kind:    RuntimeError,
			Message: err.Error(),
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	response := map[string]any{
		"error": cerr,
	}

	if jsonErr := json.NewEncoder(rw).Encode(response); jsonErr != nil {
		m.Error("Unable to encode error response", "error", jsonErr)
	}
}

// jsonCallResult writes a successful call result to the HTTP response in JSON format
func (m *MessageProcessor) jsonCallResult(rw http.ResponseWriter, result any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	response := map[string]any{
		"result": result,
	}

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		m.Error("Unable to encode result response", "error", err)
	}
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

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// Register cancel function for explicit cancellation
		func() {
			m.l.Lock()
			defer m.l.Unlock()

			if m.runningCalls[*callID] != nil {
				m.httpError(rw, "Invalid binding call:", fmt.Errorf("ambiguous call id: %s", *callID))
				return
			}
			m.runningCalls[*callID] = cancel
		}()

		// Clean up on exit
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

		// Find bound method
		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.jsonCallError(rw, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
				})
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				m.jsonCallError(rw, &CallError{
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

		// Execute the bound method synchronously
		result, err := boundMethod.Call(ctx, options.Args)

		if cerr := (*CallError)(nil); errors.As(err, &cerr) {
			var logMessage string
			switch cerr.Kind {
			case ReferenceError, TypeError:
				logMessage = "Binding call failed:"
			case RuntimeError:
				logMessage = "Bound method returned an error:"
			}
			m.Error(logMessage, "id", *callID, "error", cerr)
			m.jsonCallError(rw, cerr)
			return
		}

		// Convert result to JSON
		var jsonResult []byte
		if result != nil {
			jsonResult, err = json.Marshal(result)
			if err != nil {
				m.Error("Binding call failed:", "id", *callID, "error", err)
				m.jsonCallError(rw, &CallError{
					Kind:    TypeError,
					Message: fmt.Sprintf("error marshaling result: %s", err),
				})
				return
			}
		}

		// Log completion
		m.Info("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))

		// Return result via HTTP response
		m.jsonCallResult(rw, result)

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
