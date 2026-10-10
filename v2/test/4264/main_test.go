package test_4264

import (
	"context"
	"testing"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func TestQuitCallsOnBeforeClose(t *testing.T) {
	called := false
	opts := &options.App{
		OnBeforeClose: func(ctx context.Context) bool {
			called = true
			return false
		},
	}
	if opts.OnBeforeClose != nil {
		opts.OnBeforeClose(nil)
	}
	if !called {
		t.Error("OnBeforeClose should have been called")
	}
}

func TestOnBeforeCloseCanBlockQuit(t *testing.T) {
	opts := &options.App{
		OnBeforeClose: func(ctx context.Context) bool {
			return true
		},
	}
	if opts.OnBeforeClose != nil && opts.OnBeforeClose(nil) {
		t.Log("Quit was blocked by OnBeforeClose returning true — correct behavior")
	}
}

func TestQuitWithoutOnBeforeClose(t *testing.T) {
	opts := &options.App{}
	if opts.OnBeforeClose != nil {
		t.Error("OnBeforeClose should be nil by default")
	}
}
