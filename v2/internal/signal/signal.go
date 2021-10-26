package signal

import (
	"context"
	"os"
	gosignal "os/signal"
	"sync"
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

	// ctx
	ctx    context.Context
	cancel context.CancelFunc

	// Parent waitgroup
	wg *sync.WaitGroup
}

// NewManager creates a new signal manager
func NewManager(ctx context.Context, cancel context.CancelFunc, bus *servicebus.ServiceBus, logger *logger.Logger) (*Manager, error) {

	result := &Manager{
		bus:           bus,
		logger:        logger.CustomLogger("Event Manager"),
		signalchannel: make(chan os.Signal, 2),
		ctx:           ctx,
		cancel:        cancel,
		wg:            ctx.Value("waitgroup").(*sync.WaitGroup),
	}

	return result, nil
}

// Start the Signal Manager
func (m *Manager) Start() {

	// Hook into interrupts
	gosignal.Notify(m.signalchannel, os.Interrupt, syscall.SIGTERM)

	m.wg.Add(1)

	// Spin off signal listener and wait for either a cancellation
	// or signal
	go func() {
		select {
		case <-m.signalchannel:
			println()
			m.logger.Trace("Ctrl+C detected. Shutting down...")
			m.bus.Publish("quit", "ctrl-c pressed")

			// Start shutdown of Wails
			m.cancel()

		case <-m.ctx.Done():
		}
		m.logger.Trace("Shutdown")
		m.wg.Done()
	}()
}
