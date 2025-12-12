package application

type EventIPCTransport struct {
	app *App
}

func (t *EventIPCTransport) DispatchWailsEvent(event *CustomEvent) {
	// Snapshot windows under RLock - release lock before dispatching
	// to avoid holding the lock during potentially blocking operations.
	// This prevents deadlocks when ExecJS blocks waiting for the main thread.
	t.app.windowsLock.RLock()
	windows := make([]Window, 0, len(t.app.windows))
	for _, window := range t.app.windows {
		windows = append(windows, window)
	}
	t.app.windowsLock.RUnlock()

	for _, window := range windows {
		if event.IsCancelled() {
			return
		}
		window.DispatchWailsEvent(event)
	}
}
