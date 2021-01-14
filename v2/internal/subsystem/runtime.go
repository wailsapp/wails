package subsystem

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Runtime is the Runtime subsystem. It handles messages with topics starting
// with "runtime:"
type Runtime struct {
	quitChannel    <-chan *servicebus.Message
	runtimeChannel <-chan *servicebus.Message

	// The hooks channel allows us to hook into frontend startup
	hooksChannel     <-chan *servicebus.Message
	startupCallback  func(*runtime.Runtime)
	shutdownCallback func()

	running bool

	logger logger.CustomLogger

	// Runtime library
	runtime *runtime.Runtime
}

// NewRuntime creates a new runtime subsystem
func NewRuntime(bus *servicebus.ServiceBus, logger *logger.Logger, startupCallback func(*runtime.Runtime), shutdownCallback func()) (*Runtime, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

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
		quitChannel:      quitChannel,
		runtimeChannel:   runtimeChannel,
		hooksChannel:     hooksChannel,
		logger:           logger.CustomLogger("Runtime Subsystem"),
		runtime:          runtime.New(bus),
		startupCallback:  startupCallback,
		shutdownCallback: shutdownCallback,
	}

	return result, nil
}

// Start the subsystem
func (r *Runtime) Start() error {

	r.running = true

	// Spin off a go routine
	go func() {
		for r.running {
			select {
			case <-r.quitChannel:
				r.running = false
				break
			case hooksMessage := <-r.hooksChannel:
				r.logger.Trace(fmt.Sprintf("Received hooksmessage: %+v", hooksMessage))
				messageSlice := strings.Split(hooksMessage.Topic(), ":")
				hook := messageSlice[1]
				switch hook {
				case "startup":
					if r.startupCallback != nil {
						go r.startupCallback(r.runtime)
					} else {
						r.logger.Error("no startup callback registered!")
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
			}
		}

		// Call shutdown
		r.shutdown()
	}()

	return nil
}

// GoRuntime returns the Go Runtime object
func (r *Runtime) GoRuntime() *runtime.Runtime {
	return r.runtime
}

func (r *Runtime) shutdown() {
	if r.shutdownCallback != nil {
		go r.shutdownCallback()
	}
	r.logger.Trace("Shutdown")
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
