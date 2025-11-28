package application

import (
	"encoding/json"
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
		// Debug: log the raw args
		if args := params["args"]; len(args) > 0 {
			m.Info("Events.Emit raw args:", "args", args[0])
		}

		var event CustomEvent
		var options struct {
			Name *string         `json:"name"`
			Data json.RawMessage `json:"data"`
		}

		err := params.ToStruct(&options)
		if err != nil {
			// Fallback: if args is a JSON string (e.g., "\"frontend-test\""), treat it as event.Name
			if raw := params["args"]; len(raw) == 1 {
				var name string
				if uerr := json.Unmarshal([]byte(raw[0]), &name); uerr == nil && name != "" {
					event.Name = name
				} else {
					m.httpError(rw, "Invalid events call:", fmt.Errorf("error parsing event: %w", err))
					return
				}
			} else {
				m.httpError(rw, "Invalid events call:", fmt.Errorf("error parsing event: %w", err))
				return
			}
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
		globalApplication.Event.EmitEvent(&event)

		m.json(rw, event.IsCancelled())
		m.Info("Runtime call:", "method", "Events."+eventsMethodNames[method], "name", event.Name, "sender", event.Sender, "data", event.Data, "cancelled", event.IsCancelled())
	default:
		m.httpError(rw, "Invalid events call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
