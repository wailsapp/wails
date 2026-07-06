package runtime

import (
	"errors"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// errNoApp is returned by functions that cannot silently no-op when the
// application has not been created yet.
var errNoApp = errors.New("no application instance available")

// app returns the global v3 application instance, or nil if it has not been
// created yet.
func app() *application.App {
	return application.Get()
}

// currentWindow returns the current window, falling back to the first
// window when none is focused. Returns nil when no window exists yet.
func currentWindow() application.Window {
	a := application.Get()
	if a == nil {
		return nil
	}
	if w := a.Window.Current(); w != nil {
		return w
	}
	all := a.Window.GetAll()
	if len(all) > 0 {
		return all[0]
	}
	return nil
}

// logger returns the application logger, falling back to slog's default
// logger when the application has not been created yet.
func logger() *slog.Logger {
	if a := app(); a != nil && a.Logger != nil {
		return a.Logger
	}
	return slog.Default()
}
