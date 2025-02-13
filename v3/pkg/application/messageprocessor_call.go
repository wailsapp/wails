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

		ctx, cancel := context.WithCancel(context.WithoutCancel(r.Context()))

		// Schedule cancel in case panics happen before starting the call.
		cancelRequired := true
		defer func() {
			if cancelRequired {
				cancel()
			}
		}()

		ambiguousID := false
		func() {
			m.l.Lock()
			defer m.l.Unlock()

			if m.runningCalls[*callID] != nil {
				ambiguousID = true
			} else {
				m.runningCalls[*callID] = cancel
			}
		}()

		if ambiguousID {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("ambiguous call id: %s", *callID))
			return
		}

		m.ok(rw) // From now on, failures are reported through the error callback.

		// Log call
		var methodRef any = options.MethodName
		if options.MethodName == "" {
			methodRef = options.MethodID
		}
		m.Info("Binding call started:", "id", *callID, "method", methodRef)

		go func() {
			defer handlePanic()
			defer func() {
				m.l.Lock()
				defer m.l.Unlock()
				delete(m.runningCalls, *callID)
			}()
			defer cancel()

			var boundMethod *BoundMethod
			if options.MethodName != "" {
				boundMethod = globalApplication.bindings.Get(&options)
				if boundMethod == nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    ReferenceError,
						Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
					})
					return
				}
			} else {
				boundMethod = globalApplication.bindings.GetByID(options.MethodID)
				if boundMethod == nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    ReferenceError,
						Message: fmt.Sprintf("unknown bound method id %d", options.MethodID),
					})
					return
				}
			}

			// Prepare args for logging. This should never fail since json.Unmarshal succeeded before.
			jsonArgs, _ := json.Marshal(options.Args)
			var jsonResult []byte
			defer m.Info("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))

			// Set the context values for the window
			if window != nil {
				ctx = context.WithValue(ctx, WindowKey, window)
			}

			result, err := boundMethod.Call(ctx, options.Args)
			if cerr := (*CallError)(nil); errors.As(err, &cerr) {
				switch cerr.Kind {
				case ReferenceError, TypeError:
					m.callErrorCallback(window, "Binding call failed:", callID, cerr)
				case RuntimeError:
					m.callErrorCallback(window, "Bound method returned an error:", callID, cerr)
				}
				return
			}

			if result != nil {
				// convert result to json
				jsonResult, err = json.Marshal(result)
				if err != nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    TypeError,
						Message: fmt.Sprintf("error marshaling result: %s", err),
					})
					return
				}
			}

			m.callCallback(window, callID, string(jsonResult))
		}()

		cancelRequired = false

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
