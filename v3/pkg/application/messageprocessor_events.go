package application

import (
	"fmt"
	"encoding/json"
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
		var options struct {
			Name *string         `json:"name"`
			Data json.RawMessage `json:"data"`
		}

		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Invalid events call:", fmt.Errorf("error parsing event: %w", err))
			return
		}
		if options.Name == nil {
			m.httpError(rw, "Invalid events call:", errors.New("missing event name"))
			return
		}

		data, err := decodeEventData(*options.Name, options.Data)
		if err != nil {
			m.httpError(rw, "Events.Emit failed: ", fmt.Errorf("error parsing event data: %w", err))
			return
		}

		event.Name = *options.Name
		event.Data = data
		event.Sender = window.Name()
		globalApplication.emitEvent(&event)

		m.ok(rw)
		m.Info("Runtime call:", "method", "Events."+eventsMethodNames[method], "name", event.Name, "sender", event.Sender, "data", event.Data, "cancelled", event.IsCancelled())
	default:
		m.httpError(rw, "Invalid events call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
