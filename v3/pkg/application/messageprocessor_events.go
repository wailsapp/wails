package application

import (
	"fmt"
	"net/http"
)

func (m *MessageProcessor) processEventsMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "Emit":
		var event CustomEvent
		err := params.ToStruct(&event)
		if err != nil {
			m.httpError(rw, "Error parsing event: %s", err)
			return
		}
		if event.Name == "" {
			m.httpError(rw, "Event name must be specified")
			return
		}
		event.Sender = fmt.Sprintf("%d", window.id)
		globalApplication.Events.Emit(&event)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown event method: %s", method)
	}

}
