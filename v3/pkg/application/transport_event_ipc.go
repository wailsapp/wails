package application

type EventIPCTransport struct {
	app *App
}

func (t *EventIPCTransport) DispatchWailsEvent(event *CustomEvent) {
	// Snapshot windows under RLock
	t.app.windowsLock.RLock()
	for _, window := range t.app.windows {
		if event.IsCancelled() {
			t.app.windowsLock.RUnlock()
			return
		}
		window.DispatchWailsEvent(event)
	}
	t.app.windowsLock.RUnlock()
}
