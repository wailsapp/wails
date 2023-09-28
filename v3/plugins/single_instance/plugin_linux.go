//go:build linux

package single_instance

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGUSR2,
	)
	go func() {
		for {
			s := <-sigc
			application.Get().Show()
		}
	}()
}

func (p *Plugin) activeInstance(pid int) error {
	syscall.Kill(pid, syscall.SIGUSR2)
	return nil
}
