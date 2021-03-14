package subsystem

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Runtime is the Runtime subsystem. It handles messages with topics starting
// with "runtime:"
type Runtime struct {
	runtimeChannel <-chan *servicebus.Message

	// The hooks channel allows us to hook into frontend startup
	hooksChannel     <-chan *servicebus.Message
	startupCallback  func(*runtime.Runtime)
	shutdownCallback func()

	// quit flag
	shouldQuit bool

	logger logger.CustomLogger

	// Runtime library
	runtime *runtime.Runtime

	//ctx
	ctx context.Context

	// Startup Hook
	startupOnce sync.Once

	// Service bus
	bus *servicebus.ServiceBus
}

// NewRuntime creates a new runtime subsystem
func NewRuntime(ctx context.Context, bus *servicebus.ServiceBus, logger *logger.Logger, startupCallback func(*runtime.Runtime)) (*Runtime, error) {

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
		runtime:         runtime.New(bus),
		startupCallback: startupCallback,
		ctx:             ctx,
		bus:             bus,
	}

	return result, nil
}

// Start the subsystem
func (r *Runtime) Start() error {

	// Spin off a go routine
	go func() {
		defer r.logger.Trace("Shutdown")
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
								r.startupCallback(r.runtime)

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
			case runtimeMessage := <-r.runtimeChannel:
				r.logger.Trace(fmt.Sprintf("Received message: %+v", runtimeMessage))
				// Topics have the format: "runtime:category:call"
				messageSlice := strings.Split(runtimeMessage.Topic(), ":")
				if len(messageSlice) != 3 {
					r.logger.Error("Invalid runtime message: %#v\n", runtimeMessage)
					continue
				}

				category := messageSlice[1]
				method := messageSlice[2]
				var err error
				switch category {
				case "browser":
					err = r.processBrowserMessage(method, runtimeMessage.Data())
				default:
					err = fmt.Errorf("unknown runtime message: %+v",
						runtimeMessage)
				}

				// If we had an error, log it
				if err != nil {
					r.logger.Error(err.Error())
				}
			case <-r.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// GoRuntime returns the Go Runtime object
func (r *Runtime) GoRuntime() *runtime.Runtime {
	return r.runtime
}

func (r *Runtime) processBrowserMessage(method string, data interface{}) error {
	switch method {
	case "open":
		target, ok := data.(string)
		if !ok {
			return fmt.Errorf("expected 1 string parameter for runtime:browser:open")
		}
		go func() {
			err := r.runtime.Browser.Open(target)
			if err != nil {
				r.logger.Error(err.Error())
			}
		}()
	default:
		return fmt.Errorf("unknown method runtime:browser:%s", method)
	}
	return nil
}
