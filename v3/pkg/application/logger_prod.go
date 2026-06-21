//go:build production && !ios && !android

package application

import (
	"io"
	"log/slog"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
