//go:build windows && !production

package application

import (
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"log/slog"
	"os"
	"time"
)

func DefaultLogger(level slog.Level) *slog.Logger {
	return slog.New(tint.NewHandler(colorable.NewColorable(os.Stderr), &tint.Options{
		TimeFormat: time.StampMilli,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		Level:      level,
	}))
}
