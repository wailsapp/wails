//go:build !windows

package application

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func DefaultLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Kitchen,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
	}))
}
