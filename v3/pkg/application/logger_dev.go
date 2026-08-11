//go:build !windows && !production && !ios && !android

package application

import (
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/internal/tint"
	"github.com/mattn/go-isatty"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Kitchen,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		Level:      level,
	}))
}
