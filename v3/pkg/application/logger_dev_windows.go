//go:build windows && !production

package application

import (
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/internal/tint"
	"golang.org/x/sys/windows"
)

func DefaultLogger(level slog.Leveler) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.StampMilli,
		NoColor:    !enableConsoleColors(os.Stderr),
		Level:      level,
	}))
}

func enableConsoleColors(file *os.File) bool {
	handle := windows.Handle(file.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		// Files, pipes and invalid handles do not accept console modes.
		return false
	}
	// Preserve existing flags and enable processed output as required for VT.
	colorMode := mode | windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	if colorMode == mode {
		return true
	}
	return windows.SetConsoleMode(handle, colorMode) == nil
}
