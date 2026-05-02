package application

type EventIPCTransport struct {
	app *App
}

func (t *EventIPCTransport) DispatchWailsEvent(event *CustomEvent) {
	// Snapshot the window list under the lock, then release before dispatching.
	// DispatchWailsEvent calls ExecJS → InvokeSync which blocks until the main
	// thread executes the JS. Holding windowsLock.RLock during InvokeSync causes
	// a deadlock when the main thread (or any other goroutine) needs windowsLock
	// for write operations (NewWithOptions, Remove) — the pending writer blocks
	// new readers, and the existing readers can't complete because InvokeSync
	// needs the main thread which is waiting for the write lock.
	t.app.windowsLock.RLock()
	windows := make([]Window, 0, len(t.app.windows))
	for _, w := range t.app.windows {
		windows = append(windows, w)
	}
	t.app.windowsLock.RUnlock()

	for _, window := range windows {
		if event.IsCancelled() {
			return
		}
		window.DispatchWailsEvent(event)
	}
}
