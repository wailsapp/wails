//go:build darwin || linux

package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/atterpac/refresh/engine"
)

// shutdownSignals stop the watcher and tear down its process groups. SIGHUP is
// what a closing terminal sends; without it the watcher dies on the default
// action and leaves the app and frontend dev server running (#6328).
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP}

// trapPauseSignal mirrors the engine's own enable_pause handling, which Run
// does not install: Ctrl+Z (SIGTSTP) toggles pause instead of suspending.
func trapPauseSignal(ctx context.Context, watcherEngine *engine.Engine) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTSTP)
	go func() {
		defer signal.Stop(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				if watcherEngine.Paused() {
					watcherEngine.Resume()
				} else {
					watcherEngine.Pause()
				}
			}
		}
	}()
}
