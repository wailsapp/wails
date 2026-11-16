//go:build !windows && !production && ((!ios && !android) || server)

package application

import (
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/internal/tint"
	"golang.org/x/term"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Kitchen,
		NoColor:    !term.IsTerminal(int(os.Stderr.Fd())),
		Level:      level,
	}))
}
