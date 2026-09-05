package updater

import (
	_ "embed"
	"strings"
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

// BYOWindow wraps a caller-owned WindowHandle (typically an
// *application.WebviewWindow created outside the updater) so it can be
// passed via Config.Window. The Updater drives Show / Close / EmitEvent on
// the wrapped handle instead of creating its own window.
//
// The WindowOption interface has an unexported marker method that prevents
// arbitrary types from being assigned to Config.Window directly; this
// constructor is the bridge — call it once at Init time:
//
//	myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{...})
//	app.Updater.Init(updater.Config{
//	    Window: updater.BYOWindow(myWin),
//	    ...,
//	})
func BYOWindow(w WindowHandle) WindowOption {
	return &byoWindow{handle: w}
}

type byoWindow struct{ handle WindowHandle }

func (*byoWindow) isWindowOption() {}

// defaultBuiltinOptions returns the chrome the framework uses when the user
// doesn't override Options.
//
// The dimensions match the compact Checking / Up-to-Date / Error states so
// the window opens at its smallest natural size. If the Updater finds an
// available release (or downloads / installs / verifies one), transition()
// grows the window to availableWidth × availableHeight via WindowSizer.
// Opening small and *growing* feels like the window adapting to fit richer
// content; opening big and *shrinking* when the answer is "nothing to do"
// reads as a janky deflate after the window's already been shown — which
// is the artifact Lea flagged during interactive testing.
func defaultBuiltinOptions() WindowOptions {
	return WindowOptions{
		Title:         "Software Update",
		Width:         upToDateWidth,
		Height:        upToDateHeight,
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
		injected := "<style>" + bw.CSS + "</style>\n</head>"
		if updated := strings.Replace(html, "</head>", injected, 1); updated != html {
			html = updated
		} else {
			// No </head> (custom HTML?) — fall back to the document end.
			html += "\n<style>" + bw.CSS + "</style>"
		}
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
	case *byoWindow:
		return windowModeBYO, nil, v.handle
	}
	return windowModeBuiltin, nil, nil
}

type windowMode int

const (
	windowModeBuiltin windowMode = iota
	windowModeBYO
	windowModeNone
)
