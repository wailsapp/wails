package application

import (
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
		err := req.Args.ToStruct(&event)
		if err != nil {
			return nil, errs.WrapInvalidEventsCallErrorf(err, "error parsing event")
		}
		if event.Name == "" {
			return nil, errs.NewInvalidEventsCallErrorf("missing event name")
		}

		event.Sender = window.Name()
		globalApplication.Event.EmitEvent(&event)

		return unit, nil
	default:
		return nil, errs.NewInvalidEventsCallErrorf("unknown method: %d", req.Method)
	}
}
