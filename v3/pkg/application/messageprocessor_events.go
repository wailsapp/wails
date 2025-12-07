package application

import (
	"encoding/json"

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
		event.Sender = window.Name()
		globalApplication.Event.EmitEvent(&event)

		return event.IsCancelled(), nil
	default:
		return nil, errs.NewInvalidEventsCallErrorf("unknown method: %d", req.Method)
	}
}
