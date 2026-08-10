//go:build !windows && !production && !ios && !android

package application

import (
	"log/slog"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/wailsapp/wails/v3/internal/tint"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Kitchen,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		Level:      level,
	}))
}
