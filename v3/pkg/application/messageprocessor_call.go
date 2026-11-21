package application

import (
	"context"
	"errors"

	"encoding/json"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

type contextKey string

const (
	CallBinding            = 0
	WindowKey   contextKey = "Window"
)

type bindingCallScope struct {
	clientID   string
	windowName string
}

type bindingCallKey struct {
	scope bindingCallScope
	id    string
}

func (m *MessageProcessor) processCallCancelMethod(req *RuntimeRequest, scope bindingCallScope) (any, error) {
	callID := req.Args.AsMap().String("call-id")
	if callID == nil || *callID == "" {
		return nil, errs.NewInvalidBindingCallErrorf("missing argument 'call-id'")
	}

	var cancel func()
	func() {
		m.l.Lock()
		defer m.l.Unlock()
		cancel = m.runningCalls[bindingCallKey{scope: scope, id: *callID}]
	}()

	if cancel != nil {
		cancel()
		m.Debug("Binding call cancelled:", "id", *callID)
	}
	return unit, nil
}

func (m *MessageProcessor) processCallMethod(ctx context.Context, req *RuntimeRequest, window Window) (any, error) {
	run, err := m.prepareCallMethod(ctx, req, window, bindingCallScope{})
	if err != nil {
		return nil, err
	}
	return run()
}

// prepareCallMethod registers the binding context and returns the work that
// runs and unregisters it. The Android bridge prepares calls before dispatching
// their workers, so a later CancelCall can find even a queued binding.
func (m *MessageProcessor) prepareCallMethod(ctx context.Context, req *RuntimeRequest, window Window, scope bindingCallScope) (func() (any, error), error) {
	callID := req.Args.AsMap().String("call-id")
	if callID == nil || *callID == "" {
		return nil, errs.NewInvalidBindingCallErrorf("missing argument 'call-id'")
	}
	if req.Method != CallBinding {
		return nil, errs.NewInvalidBindingCallErrorf("unknown method: %d", req.Method)
	}

	var options CallOptions
	if err := req.Args.ToStruct(&options); err != nil {
		return nil, errs.WrapInvalidBindingCallErrorf(err, "error parsing call options")
	}

	var methodRef any = options.MethodName
	if options.MethodName == "" {
		methodRef = options.MethodID
	}
	m.Debug("Binding call started:", "id", *callID, "method", methodRef)

	// Calls may originate from windowless transports.
	var windowID uint
	if window != nil {
		windowID = window.ID()
	}
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	key := bindingCallKey{scope: scope, id: *callID}

	m.l.Lock()
	if m.runningCalls[key] != nil {
		m.l.Unlock()
		cancel()
		return nil, errs.NewInvalidBindingCallErrorf("ambiguous call id: %s", *callID)
	}
	m.runningCalls[key] = cancel
	if window != nil {
		if m.windowCalls[windowID] == nil {
			m.windowCalls[windowID] = make(map[bindingCallKey]bool)
		}
		m.windowCalls[windowID][key] = true
	}
	m.l.Unlock()

	return func() (any, error) {
		defer func() {
			m.l.Lock()
			defer m.l.Unlock()
			delete(m.runningCalls, key)
			if windowCalls, exists := m.windowCalls[windowID]; window != nil && exists {
				delete(windowCalls, key)
				if len(windowCalls) == 0 {
					delete(m.windowCalls, windowID)
				}
			}
		}()
		defer cancel()

		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				return nil, errs.NewBindingCallFailedErrorf("unknown bound method name '%s'", options.MethodName)
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				return nil, errs.NewBindingCallFailedErrorf("unknown bound method id %d", options.MethodID)
			}
		}

		jsonArgs, _ := json.Marshal(options.Args)
		var result any
		defer func() {
			var jsonResult []byte
			jsonResult, _ = json.Marshal(result)
			m.Debug("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))
		}()

		// Set the context values for the window
		if window != nil {
			ctx = context.WithValue(ctx, WindowKey, window)
		}

		result, err := boundMethod.Call(ctx, options.Args)
		if cerr := (*CallError)(nil); errors.As(err, &cerr) {
			switch cerr.Kind {
			case ReferenceError, TypeError:
				return nil, errs.WrapBindingCallFailedErrorf(cerr, "failed to call binding")
			case RuntimeError:

				return nil, errs.WrapBindingCallFailedErrorf(cerr, "Bound method returned an error")
			}
		}
		if err != nil {
			return nil, errs.WrapBindingCallFailedErrorf(err, "failed to call binding")
		}
		return result, nil
	}, nil
}
