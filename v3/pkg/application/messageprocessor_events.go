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

func (m *MessageProcessor) processEventsMethod(method int, rw http.ResponseWriter, _ *http.Request, window Window, params QueryParams) {

	var event CustomEvent

	switch method {
	case EventsEmit:
		err := params.ToStruct(&event)
		if err != nil {
			m.httpError(rw, "Error parsing event: %s", err.Error())
			return
		}
		if event.Name == "" {
			m.httpError(rw, "Event name must be specified")
			return
		}
		event.Sender = window.Name()
		globalApplication.customEventProcessor.Emit(&event)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown event method: %d", method)
		return
	}

	m.Info("Runtime Call:", "method", "Events."+eventsMethodNames[method], "name", event.Name, "sender", event.Sender, "data", event.Data, "cancelled", event.IsCancelled())

}
