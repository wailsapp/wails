package application

type EventIPCTransport struct {
	app *App
}

func (t *EventIPCTransport) DispatchWailsEvent(event *CustomEvent) {
	// Snapshot windows under RLock
	t.app.windowsLock.RLock()
	defer t.app.windowsLock.RUnlock()
	for _, window := range t.app.windows {
		if event.IsCancelled() {
			return
		}
		window.DispatchWailsEvent(event)
	}
}
