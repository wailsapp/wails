package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wailsapp/wails/v3/pkg/errs"
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

func (m *MessageProcessor) processCallCancelMethod(req *RuntimeRequest) (any, error) {
	args, err := req.Params.Args()
	if err != nil {
		return nil, errs.WrapInvalidBindingCallErrorf(err, "unable to parse arguments")
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		return nil, errs.NewInvalidBindingCallErrorf("missing argument 'call-id'")
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
	return unit, nil
}

func (m *MessageProcessor) processCallMethod(ctx context.Context, req *RuntimeRequest, window Window) (any, error) {
	args, err := req.Params.Args()
	if err != nil {
		return nil, errs.NewInvalidBindingCallErrorf("unable to parse arguments")
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		return nil, errs.NewInvalidBindingCallErrorf("missing argument 'call-id'")
	}

	switch req.Method {
	case CallBinding:
		var options CallOptions
		err := req.Params.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidBindingCallErrorf(err, "error parsing call options")
		}

		ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))

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
			return nil, errs.NewInvalidBindingCallErrorf("ambiguous call id: %s", *callID)
		}

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
				return nil, errs.NewBindingCallFailedErrorf("unknown bound method name '%s'", options.MethodName)
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				return nil, errs.NewBindingCallFailedErrorf("unknown bound method id %d", options.MethodID)
			}
		}

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
		return result, nil

	default:
		return nil, errs.NewInvalidBindingCallErrorf("unknown method: %d", req.Method)
	}
}
