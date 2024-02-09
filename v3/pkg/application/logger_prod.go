//go:build production

package application

import (
	"io"
	"log/slog"
)

func DefaultLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
