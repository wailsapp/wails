package application

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	EventsEmit = 0
)

var eventsMethodNames = map[int]string{
	EventsEmit: "Emit",
}

func (m *MessageProcessor) processEventsMethod(method int, rw http.ResponseWriter, _ *http.Request, window Window, params QueryParams) {
	switch method {
	case EventsEmit:
		var event CustomEvent
		err := params.ToStruct(&event)
		if err != nil {
			m.httpError(rw, "Invalid events call:", fmt.Errorf("error parsing event: %w", err))
			return
		}
		if event.Name == "" {
			m.httpError(rw, "Invalid events call:", errors.New("missing event name"))
			return
		}

		event.Sender = window.Name()
		globalApplication.customEventProcessor.Emit(&event)

		m.ok(rw)
		m.Info("Runtime call:", "method", "Events."+eventsMethodNames[method], "name", event.Name, "sender", event.Sender, "data", event.Data, "cancelled", event.IsCancelled())
	default:
		m.httpError(rw, "Invalid events call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
