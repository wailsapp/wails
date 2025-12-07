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

func (m *MessageProcessor) processCallCancelMethod(req *RuntimeRequest) (any, error) {
	callID := req.Args.AsMap().String("call-id")
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
		m.Info("Binding call cancelled:", "id", *callID)
	}
	return unit, nil
}

func (m *MessageProcessor) processCallMethod(ctx context.Context, req *RuntimeRequest, window Window) (any, error) {
	callID := req.Args.AsMap().String("call-id")
	if callID == nil || *callID == "" {
		return nil, errs.NewInvalidBindingCallErrorf("missing argument 'call-id'")
	}

	switch req.Method {
	case CallBinding:
		var options CallOptions
		err := req.Args.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidBindingCallErrorf(err, "error parsing call options")
		}

		// Log call
		var methodRef any = options.MethodName
		if options.MethodName == "" {
			methodRef = options.MethodID
		}
		m.Info("Binding call started:", "id", *callID, "method", methodRef)

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

		jsonArgs, _ := json.Marshal(options.Args)
		var result any
		defer func() {
			var jsonResult []byte
			jsonResult, _ = json.Marshal(result)
			m.Info("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))
		}()

		// Set the context values for the window
		if window != nil {
			ctx = context.WithValue(ctx, WindowKey, window)
		}

		result, err = boundMethod.Call(ctx, options.Args)
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

	default:
		return nil, errs.NewInvalidBindingCallErrorf("unknown method: %d", req.Method)
	}
}
