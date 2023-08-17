package application

import (
	"net/http"
)

func (m *MessageProcessor) processEventsMethod(method string, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {

	var event WailsEvent

	switch method {
	case "Emit":
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

	m.Info("Runtime:", "method", "Events."+method, "event", event)

}
