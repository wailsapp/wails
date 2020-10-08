package signal

import (
	"os"
	gosignal "os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Manager manages signals such as CTRL-C
type Manager struct {
	// Service Bus
	bus *servicebus.ServiceBus

	// logger
	logger logger.CustomLogger

	// signalChannel
	signalchannel chan os.Signal

	// Quit channel
	quitChannel <-chan *servicebus.Message
}

// NewManager creates a new signal manager
func NewManager(bus *servicebus.ServiceBus, logger *logger.Logger) (*Manager, error) {

	// Register quit channel
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	result := &Manager{
		bus:           bus,
		logger:        logger.CustomLogger("Event Manager"),
		signalchannel: make(chan os.Signal, 2),
		quitChannel:   quitChannel,
	}

	return result, nil
}

// Start the Signal Manager
func (m *Manager) Start() {

	// Hook into interrupts
	gosignal.Notify(m.signalchannel, os.Interrupt, syscall.SIGTERM)

	// Spin off signal listener
	go func() {
		running := true
		for running {
			select {
			case <-m.signalchannel:
				println()
				m.logger.Trace("Ctrl+C detected. Shutting down...")
				m.bus.Publish("quit", "ctrl-c pressed")
			case <-m.quitChannel:
				running = false
				break
			}
		}
		m.logger.Trace("Shutdown")
	}()
}
