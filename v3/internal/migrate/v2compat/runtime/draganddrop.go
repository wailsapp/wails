package runtime

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	fileDropLock          sync.Mutex
	fileDropUnsubscribers []func()
	// fileDropGeneration invalidates listeners registered by earlier
	// OnFileDrop calls once OnFileDropOff has run, including the window
	// creation hooks which cannot be unregistered.
	fileDropGeneration int
)

// OnFileDrop mirrors the v2 runtime.OnFileDrop function. The callback is
// registered on the current window and on any window created afterwards.
//
// The drop coordinates are not exposed to the application side in v3, so the
// callback always receives x=0, y=0. Note that file-drop events only fire for
// windows created with EnableFileDrop: true in their
// application.WebviewWindowOptions.
// v3 equivalent: window.OnWindowEvent(events.Common.WindowFilesDropped, ...).
func OnFileDrop(_ context.Context, callback func(x, y int, paths []string)) {
	a := app()
	if a == nil || callback == nil {
		return
	}

	fileDropLock.Lock()
	generation := fileDropGeneration
	fileDropLock.Unlock()

	register := func(window application.Window) {
		unsubscribe := window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
			callback(0, 0, event.Context().DroppedFiles())
		})
		fileDropLock.Lock()
		if fileDropGeneration != generation {
			// OnFileDropOff was called in the meantime.
			fileDropLock.Unlock()
			unsubscribe()
			return
		}
		fileDropUnsubscribers = append(fileDropUnsubscribers, unsubscribe)
		fileDropLock.Unlock()
	}

	if w := currentWindow(); w != nil {
		register(w)
	}
	a.Window.OnCreate(func(window application.Window) {
		fileDropLock.Lock()
		stale := fileDropGeneration != generation
		fileDropLock.Unlock()
		if stale {
			return
		}
		register(window)
	})
}

// OnFileDropOff mirrors the v2 runtime.OnFileDropOff function. It removes all
// file-drop listeners registered via OnFileDrop.
// v3 equivalent: calling the unsubscribe function returned by window.OnWindowEvent.
func OnFileDropOff(_ context.Context) {
	fileDropLock.Lock()
	fileDropGeneration++
	unsubscribers := fileDropUnsubscribers
	fileDropUnsubscribers = nil
	fileDropLock.Unlock()

	for _, unsubscribe := range unsubscribers {
		unsubscribe()
	}
}
