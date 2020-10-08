package subsystem

import (
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/internal/runtime/goruntime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Binding is the Binding subsystem. It manages all service bus messages
// starting with "binding".
type Binding struct {
	quitChannel    <-chan *servicebus.Message
	bindingChannel <-chan *servicebus.Message
	running        bool

	// Binding db
	bindings *binding.Bindings

	// logger
	logger logger.CustomLogger

	// runtime
	runtime *goruntime.Runtime
}

// NewBinding creates a new binding subsystem. Uses the given bindings db for reference.
func NewBinding(bus *servicebus.ServiceBus, logger *logger.Logger, bindings *binding.Bindings, runtime *goruntime.Runtime) (*Binding, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	// Subscribe to event messages
	bindingChannel, err := bus.Subscribe("binding")
	if err != nil {
		return nil, err
	}

	result := &Binding{
		quitChannel:    quitChannel,
		bindingChannel: bindingChannel,
		logger:         logger.CustomLogger("Binding Subsystem"),
		bindings:       bindings,
		runtime:        runtime,
	}

	// Call WailsInit methods once the frontend is loaded
	// TODO: Double check that this is actually being emitted
	// when we want it to be
	runtime.Events.On("wails:loaded", func(...interface{}) {
		result.logger.Trace("Calling WailsInit() methods")
		result.CallWailsInit()
	})

	return result, nil
}

// Start the subsystem
func (b *Binding) Start() error {

	b.running = true

	b.logger.Trace("Starting")

	// Spin off a go routine
	go func() {
		for b.running {
			select {
			case <-b.quitChannel:
				b.running = false
			case bindingMessage := <-b.bindingChannel:
				b.logger.Trace("Got binding message: %+v", bindingMessage)
			}
		}

		// Call shutdown
		b.shutdown()
	}()

	return nil
}

// CallWailsInit will callback to the registered WailsInit
// methods with the runtime object
func (b *Binding) CallWailsInit() error {
	for _, wailsinit := range b.bindings.DB().WailsInitMethods() {
		_, err := wailsinit.Call([]interface{}{b.runtime})
		if err != nil {
			return err
		}
	}
	return nil
}

// CallWailsShutdown will callback to the registered WailsShutdown
// methods with the runtime object
func (b *Binding) CallWailsShutdown() error {
	for _, wailsshutdown := range b.bindings.DB().WailsShutdownMethods() {
		_, err := wailsshutdown.Call([]interface{}{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Binding) shutdown() {
	b.CallWailsShutdown()
	b.logger.Trace("Shutdown")
}
