package updater

import (
	_ "embed"
)

//go:embed assets/window.html
var defaultWindowHTML string

// BuiltinWindow customises the framework's default update window. Embed it
// in Config.Window when you want to override the HTML, layer extra CSS, or
// change the window chrome (frameless, size, etc.).
type BuiltinWindow struct {
	// HTML, if non-empty, replaces the default template entirely. The
	// replacement is expected to listen to `updater:*` events and emit
	// `updater:user:*` events using the standard Wails runtime — exactly
	// what the default template does.
	HTML string

	// CSS, if non-empty, is appended to the default template inside a final
	// <style> tag. Ignored when HTML overrides the template entirely.
	CSS string

	// Options overrides the chrome of the built-in window. Zero values fall
	// back to the framework's sensible defaults (small, centred, resizable).
	Options WindowOptions
}

func (*BuiltinWindow) isWindowOption() {}

// windowNoneType is the concrete singleton type behind WindowNone. Declared
// as a named type so a user-supplied value can be checked with ==.
type windowNoneType struct{}

func (windowNoneType) isWindowOption() {}

// WindowNone selects headless mode: the Updater drives the flow but never
// asks the host to open a window. Subscribe to events from your own UI.
var WindowNone WindowOption = windowNoneType{}

// defaultBuiltinOptions returns the chrome the framework uses when the user
// doesn't override Options.
func defaultBuiltinOptions() WindowOptions {
	return WindowOptions{
		Title:         "Update",
		Width:         480,
		Height:        420,
		Frameless:     false,
		AlwaysOnTop:   false,
		DisableResize: false,
	}
}

// composeHTML builds the HTML for the built-in window. Callers may supply a
// BuiltinWindow override; nil means "use defaults."
func composeHTML(bw *BuiltinWindow) string {
	if bw != nil && bw.HTML != "" {
		return bw.HTML
	}
	html := defaultWindowHTML
	if bw != nil && bw.CSS != "" {
		html += "\n<style>" + bw.CSS + "</style>"
	}
	return html
}

// composeWindowOptions returns the WindowOptions the Updater asks its host
// to open the built-in window with. Caller-supplied overrides on a
// BuiltinWindow take precedence over framework defaults, but framework
// defaults fill any zero values.
func composeWindowOptions(bw *BuiltinWindow) WindowOptions {
	opts := defaultBuiltinOptions()
	if bw == nil {
		return opts
	}
	o := bw.Options
	if o.Title != "" {
		opts.Title = o.Title
	}
	if o.Width > 0 {
		opts.Width = o.Width
	}
	if o.Height > 0 {
		opts.Height = o.Height
	}
	opts.Frameless = o.Frameless
	opts.AlwaysOnTop = o.AlwaysOnTop
	opts.DisableResize = o.DisableResize
	return opts
}

// classifyWindowOption walks a user-supplied Config.Window value and returns
// a small enum describing the runtime mode + any concrete BuiltinWindow
// configuration. Falls back to "builtin defaults" for nil input.
func classifyWindowOption(opt WindowOption) (mode windowMode, bw *BuiltinWindow, byo WindowHandle) {
	if opt == nil {
		return windowModeBuiltin, nil, nil
	}
	switch v := opt.(type) {
	case *BuiltinWindow:
		return windowModeBuiltin, v, nil
	case windowNoneType:
		return windowModeNone, nil, nil
	case WindowHandle:
		return windowModeBYO, nil, v
	}
	return windowModeBuiltin, nil, nil
}

type windowMode int

const (
	windowModeBuiltin windowMode = iota
	windowModeBYO
	windowModeNone
)
