package application

import (
	"encoding/json"

	"github.com/wailsapp/wails/v3/pkg/application/monitor"
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	EventsEmit = 0
)

var eventsMethodNames = map[int]string{
	EventsEmit: "Emit",
}

func (m *MessageProcessor) processEventsMethod(req *RuntimeRequest, window Window) (any, error) {
	switch req.Method {
	case EventsEmit:
		var event CustomEvent
		var options struct {
			Name *string         `json:"name"`
			Data json.RawMessage `json:"data"`
		}

		err := req.Args.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidEventsCallErrorf(err, "error parsing event")
		}
		if options.Name == nil {
			return nil, errs.NewInvalidEventsCallErrorf("missing event name")
		}

		data, err := decodeEventData(*options.Name, options.Data)
		if err != nil {
			return nil, errs.WrapInvalidEventsCallErrorf(err, "error parsing event data")
		}

		event.Name = *options.Name
		event.Data = data
		if window != nil {
			event.Sender = window.Name()
		}

		// IPC monitor tap: inbound (JS->Go) event.
		if monitor.Enabled() {
			var monWindow string
			if window != nil {
				monWindow = window.Name()
			}
			monitor.Emit(monitor.Trace{
				Kind:       "event",
				Dir:        "in",
				Object:     eventsRequest,
				ObjectName: "Events",
				Method:     event.Name,
				Window:     monWindow,
				Args:       json.RawMessage(options.Data),
			})
		}

		globalApplication.Event.EmitEvent(&event)

		return event.IsCancelled(), nil
	default:
		return nil, errs.NewInvalidEventsCallErrorf("unknown method: %d", req.Method)
	}
}
