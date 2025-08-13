package dispatcher

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Dispatcher struct {
	log                  *logger.Logger
	bindings             *binding.Bindings
	events               frontend.Events
	bindingsDB           *binding.DB
	ctx                  context.Context
	errfmt               options.ErrorFormatter
	disablePanicRecovery bool
}

func NewDispatcher(ctx context.Context, log *logger.Logger, bindings *binding.Bindings, events frontend.Events, errfmt options.ErrorFormatter, disablePanicRecovery bool) *Dispatcher {
	return &Dispatcher{
		log:                  log,
		bindings:             bindings,
		events:               events,
		bindingsDB:           bindings.DB(),
		ctx:                  ctx,
		errfmt:               errfmt,
		disablePanicRecovery: disablePanicRecovery,
	}
}

func (d *Dispatcher) ProcessMessage(message string, sender frontend.Frontend) (_ string, err error) {
	if !d.disablePanicRecovery {
		defer func() {
			if e := recover(); e != nil {
				if errPanic, ok := e.(error); ok {
					err = errPanic
				} else {
					err = fmt.Errorf("%v", e)
				}
			}
			if err != nil {
				d.log.Error("process message error: %s -> %s", message, err)
			}
		}()
	}

	if message == "" {
		return "", errors.New("No message to process")
	}
	switch message[0] {
	case 'L':
		return d.processLogMessage(message)
	case 'E':
		return d.processEventMessage(message, sender)
	case 'C':
		return d.processCallMessage(message, sender)
	case 'c':
		return d.processSecureCallMessage(message, sender)
	case 'W':
		return d.processWindowMessage(message, sender)
	case 'B':
		return d.processBrowserMessage(message, sender)
	case 'D':
		return d.processDragAndDropMessage(message)
	case 'Q':
		sender.Quit()
		return "", nil
	case 'S':
		sender.Show()
		return "", nil
	case 'H':
		sender.Hide()
		return "", nil
	default:
		return "", errors.New("Unknown message from front end: " + message)
	}
}
