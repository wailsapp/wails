package application

import (
	"net/http"
)

const (
	EventsEmit = 0
)

var eventsMethodNames = map[int]string{
	EventsEmit: "Emit",
}

func (m *MessageProcessor) processEventsMethod(method int, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {

	var event WailsEvent

	switch method {
	case EventsEmit:
		err := params.ToStruct(&event)
		if err != nil {
			m.httpError(rw, "Error parsing event: %s", err)
			return
		}
		if event.Name == "" {
			m.httpError(rw, "Event name must be specified")
			return
		}
		event.Sender = window.Name()
		globalApplication.Events.Emit(&event)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown event method: %s", method)
		return
	}

	m.Info("Runtime:", "method", "Events."+eventsMethodNames[method], "event", event)

}
