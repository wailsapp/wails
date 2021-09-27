package subsystem

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"strings"
	"sync"
)

// Runtime is the Runtime subsystem. It handles messages with topics starting
// with "runtime:"
type Runtime struct {
	runtimeChannel <-chan *servicebus.Message

	// The hooks channel allows us to hook into frontend startup
	hooksChannel     <-chan *servicebus.Message
	startupCallback  func(ctx context.Context)
	shutdownCallback func()

	// quit flag
	shouldQuit bool

	logger logger.CustomLogger

	//ctx
	ctx context.Context

	// OnStartup Hook
	startupOnce sync.Once

	// Service bus
	bus *servicebus.ServiceBus
}

// NewRuntime creates a new runtime subsystem
func NewRuntime(ctx context.Context, bus *servicebus.ServiceBus, logger *logger.Logger, startupCallback func(context.Context)) (*Runtime, error) {

	// Subscribe to log messages
	runtimeChannel, err := bus.Subscribe("runtime:")
	if err != nil {
		return nil, err
	}

	// Subscribe to log messages
	hooksChannel, err := bus.Subscribe("hooks:")
	if err != nil {
		return nil, err
	}

	result := &Runtime{
		runtimeChannel:  runtimeChannel,
		hooksChannel:    hooksChannel,
		logger:          logger.CustomLogger("Runtime Subsystem"),
		startupCallback: startupCallback,
		bus:             bus,
	}
	result.ctx = context.WithValue(ctx, "bus", bus)

	return result, nil
}

// Start the subsystem
func (r *Runtime) Start() error {

	// Spin off a go routine
	go func() {
		defer r.logger.Trace("OnShutdown")
		for {
			select {
			case hooksMessage := <-r.hooksChannel:
				r.logger.Trace(fmt.Sprintf("Received hooksmessage: %+v", hooksMessage))
				messageSlice := strings.Split(hooksMessage.Topic(), ":")
				hook := messageSlice[1]
				switch hook {
				case "startup":
					if r.startupCallback != nil {
						r.startupOnce.Do(func() {
							go func() {
								r.startupCallback(r.ctx)
								// If we got a url, publish it now startup completed
								url, ok := hooksMessage.Data().(string)
								if ok && len(url) > 0 {
									r.bus.Publish("url:handler", url)
								}
							}()
						})
					} else {
						r.logger.Warning("no startup callback registered!")
					}
				default:
					r.logger.Error("unknown hook message: %+v", hooksMessage)
					continue
				}
			case <-r.ctx.Done():
				return
			}
		}
	}()

	return nil
}
