package subsystem

import (
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Binding is the Binding subsystem. It manages all service bus messages
// starting with "binding".
type Binding struct {
	bindingChannel <-chan *servicebus.Message

	running bool

	// Binding db
	bindings *binding.Bindings

	// logger
	logger logger.CustomLogger
}

// NewBinding creates a new binding subsystem. Uses the given bindings db for reference.
func NewBinding(bus *servicebus.ServiceBus, logger *logger.Logger, bindings *binding.Bindings) (*Binding, error) {

	// Subscribe to event messages
	bindingChannel, err := bus.Subscribe("binding")
	if err != nil {
		return nil, err
	}

	result := &Binding{
		bindingChannel: bindingChannel,
		logger:         logger.CustomLogger("Binding Subsystem"),
		bindings:       bindings,
	}

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
			case bindingMessage := <-b.bindingChannel:
				b.logger.Trace("Got binding message: %+v", bindingMessage)
			}
		}
		b.logger.Trace("Shutdown")
	}()

	return nil
}

func (b *Binding) Close() {
	b.running = false
}
