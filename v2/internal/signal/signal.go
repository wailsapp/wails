package signal

import (
	"os"
	gosignal "os/signal"
	"sync"
	"syscall"
)

var signalChannel = make(chan os.Signal, 2)

var (
	callbacks []func()
	lock      sync.Mutex
)

func OnShutdown(callback func()) {
	lock.Lock()
	defer lock.Unlock()
	callbacks = append(callbacks, callback)
}

// Start the Signal Manager
func Start() {
	// Hook into interrupts
	gosignal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Spin off signal listener and wait for either a cancellation
	// or signal
	go func() {
		<-signalChannel
		println("")
		println("Ctrl+C detected. Shutting down...")
		for _, callback := range callbacks {
			callback()
		}
	}()
}
