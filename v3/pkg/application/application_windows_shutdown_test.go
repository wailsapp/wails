//go:build windows && !server

package application

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// Guards #6368: closing the last window on Windows must run the application
// shutdown hooks (OnShutdown and the tasks registered through it), as Linux
// and macOS do, instead of just posting WM_QUIT.
func TestUnregisterLastWindowRunsShutdown(t *testing.T) {
	tests := []struct {
		name     string
		options  Options
		wantHook bool
	}{
		{
			name:     "default options run shutdown hooks",
			options:  Options{},
			wantHook: true,
		},
		{
			name: "DisableQuitOnLastWindowClosed keeps app alive",
			options: Options{
				Windows: WindowsOptions{DisableQuitOnLastWindowClosed: true},
			},
			wantHook: false,
		},
		{
			name:     "ShouldQuit false keeps app alive",
			options:  Options{ShouldQuit: func() bool { return false }},
			wantHook: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// cleanup() dispatches through InvokeSync. Reporting the current
			// thread as the main thread makes it execute inline; without the
			// lock the goroutine could migrate before that check.
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()

			previous := globalApplication
			t.Cleanup(func() {
				globalApplication = previous
			})

			app := &App{
				ctx:     context.Background(),
				cancel:  func() {},
				options: tt.options,
				Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
			}
			wa := &windowsApp{
				parent:               app,
				windowMap:            make(map[w32.HWND]*windowsWebviewWindow),
				mainThreadID:         w32.GetCurrentThreadId(),
				mainThreadWindowHWND: w32.HWND(1),
			}
			app.impl = wa
			globalApplication = app

			var hookCalls atomic.Int32
			app.OnShutdown(func() {
				hookCalls.Add(1)
			})

			fakeWindow := &windowsWebviewWindow{hwnd: w32.HWND(0xBEEF)}
			wa.windowMap[fakeWindow.hwnd] = fakeWindow

			wa.unregisterWindow(fakeWindow)

			// A second call must not run the hooks again.
			if tt.wantHook {
				wa.unregisterWindow(fakeWindow)
			}

			// destroy() leaves WM_QUIT on this thread's queue; drain it so it
			// cannot leak into another test running on the same thread.
			var msg w32.MSG
			for w32.PeekMessage(&msg, 0, 0, 0, w32.PM_REMOVE) {
			}

			got := hookCalls.Load()
			if tt.wantHook {
				if got != 1 {
					t.Fatalf("shutdown hook ran %d times, want exactly 1", got)
				}
				if !app.performingShutdown {
					t.Error("performingShutdown = false after closing the last window, want true")
				}
			} else {
				if got != 0 {
					t.Errorf("shutdown hook ran %d times, want 0", got)
				}
				if app.performingShutdown {
					t.Error("performingShutdown = true, want false")
				}
			}
		})
	}
}
