//go:build windows && !production

package application

import (
	"log/slog"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/wailsapp/wails/v3/internal/tint"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(tint.NewHandler(colorable.NewColorable(os.Stderr), &tint.Options{
		TimeFormat: time.StampMilli,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		Level:      level,
	}))
}
