//go:build !darwin && !linux

package commands

import (
	"context"
	"os"
	"syscall"

	"github.com/atterpac/refresh/engine"
)

// shutdownSignals stop the watcher and tear down its processes.
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// trapPauseSignal is a no-op: there is no suspend signal to repurpose here.
func trapPauseSignal(context.Context, *engine.Engine) {}
