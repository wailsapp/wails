package signal

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler struct {
	cleanup     func()
	ExitMessage func(sig os.Signal) string
	MaxSignal   int
	Logger      *slog.Logger
	LogLevel    slog.Level
}

func NewSignalHandler(cleanup func()) *SignalHandler {
	return &SignalHandler{
		cleanup:     cleanup,
		ExitMessage: func(sig os.Signal) string { return fmt.Sprintf("Received signal: %v. Quitting...\n", sig) },
		MaxSignal:   3,
		Logger:      slog.New(slog.NewTextHandler(os.Stderr, nil)),
		LogLevel:    slog.LevelInfo,
	}
}

func (s *SignalHandler) Start() {
	ctrlC := make(chan os.Signal, s.MaxSignal)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	go func() {
		for i := 1; i <= s.MaxSignal; i++ {
			sig := <-ctrlC

			if i == 1 {
				s.Logger.Info(s.ExitMessage(sig))
				s.cleanup()
				break
			} else if i < s.MaxSignal {
				s.Logger.Info(fmt.Sprintf("Received signal: %v. Press CTRL+C %d more times to force quit...\n", sig, s.MaxSignal-i))
				continue
			} else {
				s.Logger.Info(fmt.Sprintf("Received signal: %v. Force quitting...\n", sig))
				os.Exit(1)
			}
		}
	}()
}
