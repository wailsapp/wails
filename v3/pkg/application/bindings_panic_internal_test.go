package application

import (
	"context"
	"errors"
	"testing"
)

type panicTestService struct{}

func (panicTestService) Boom() { panic("kaboom") }

// TestBoundMethodPanicInvokesPanicHandler guards the second half of #5037:
// with a custom PanicHandler registered, a panic in a bound method must both
// invoke the handler (telemetry preserved) AND still reject the call with a
// *CallError. The pre-fix code returned (nil, nil) on this path, silently
// resolving the frontend promise with null.
//
// This is an internal test so it can swap the PanicHandler on the global
// application (application.New is a singleton, so the handler cannot be
// installed from the external test package). The previous handler is restored
// on cleanup to avoid affecting other tests.
func TestBoundMethodPanicInvokesPanicHandler(t *testing.T) {
	app := New(Options{})
	prev := app.options.PanicHandler
	t.Cleanup(func() { app.options.PanicHandler = prev })

	var handled *PanicDetails
	app.options.PanicHandler = func(d *PanicDetails) { handled = d }

	bindings := NewBindings(nil, nil)
	if err := bindings.Add(NewService(&panicTestService{})); err != nil {
		t.Fatalf("bindings.Add() error = %v", err)
	}

	var bound *BoundMethod
	for _, m := range bindings.boundMethods {
		bound = m
	}
	if bound == nil {
		t.Fatal("no bound method registered")
	}

	result, err := bound.Call(context.TODO(), nil)

	// (1) the handler is invoked with the panic details
	if handled == nil {
		t.Error("PanicHandler was not invoked")
	} else if handled.Error == nil {
		t.Error("PanicDetails.Error is nil")
	}

	// (2) the call still rejects with a RuntimeError CallError
	if result != nil {
		t.Errorf("result = %v, expected nil", result)
	}
	var cerr *CallError
	if !errors.As(err, &cerr) {
		t.Fatalf("err = %#v, expected *CallError", err)
	}
	if cerr.Kind != RuntimeError {
		t.Errorf("err.Kind = %q, expected RuntimeError", cerr.Kind)
	}
}
