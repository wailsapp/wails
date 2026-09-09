package application

import (
	"github.com/wailsapp/wails/v3/pkg/updater"
)

// updaterHost adapts the application's EventManager and window factory to
// the small interface the updater package needs. Defined here (and not in
// the updater package) to keep the dependency direction clean: application
// imports updater, never the other way round.
type updaterHost struct {
	app *App
}

func newUpdaterHost(a *App) *updaterHost { return &updaterHost{app: a} }

// Emit forwards to the app's event bus.
func (h *updaterHost) Emit(name string, data ...any) bool {
	return h.app.Event.Emit(name, data...)
}

// Quit forwards to the app's normal shutdown sequence — the updater calls
// this from Restart after spawning the helper so the helper's "wait for
// parent PID to exit" step actually completes.
func (h *updaterHost) Quit() { h.app.Quit() }

// OnEvent subscribes via the event bus, unwrapping CustomEvent.Data so the
// updater's callback signature is decoupled from application's CustomEvent
// type.
func (h *updaterHost) OnEvent(name string, cb func(payload any)) func() {
	return h.app.Event.On(name, func(e *CustomEvent) {
		var d any
		if e != nil {
			d = e.Data
		}
		cb(d)
	})
}

// OpenWindow creates a real webview window with the supplied options and
// returns a handle the updater can drive. The returned handle wraps a
// *WebviewWindow so the methods the updater calls do not collide with the
// fluent-returning variants on the application's Window interface.
func (h *updaterHost) OpenWindow(opts updater.WindowOptions) updater.WindowHandle {
	// HTML must be supplied at construction time, not via a post-creation
	// SetHTML() call. On webkit2gtk (and likely WebView2 in some
	// configurations) SetHTML loads the supplied document into the
	// about:blank context, which does not have the Wails runtime
	// (window.wails, _wails.dispatchWailsEvent) injected — JS event emits
	// from the loaded page silently no-op.
	wopts := WebviewWindowOptions{
		Title:         opts.Title,
		Width:         opts.Width,
		Height:        opts.Height,
		Frameless:     opts.Frameless,
		AlwaysOnTop:   opts.AlwaysOnTop,
		DisableResize: opts.DisableResize,
		HTML:          opts.InitialHTML,
		// The updater's window has to drive Restart / Install / Skip / etc.
		// via the simple postMessage path — its HTML is loaded with no
		// asset-server origin so the modern HTTP runtime is unreachable.
		// HTML for this window is fully controlled by the framework (or by
		// a developer who opted in via BYOWindow), so the broader threat
		// model AllowSimpleEventEmit guards against doesn't apply here.
		AllowSimpleEventEmit: true,
	}
	win := h.app.Window.NewWithOptions(wopts)
	win.Show()
	return &updaterWindowHandle{win: win}
}

// updaterWindowHandle bridges *WebviewWindow into the updater.WindowHandle
// interface. The wrapper drops the fluent return values so the smaller
// updater interface stays free of application-package types.
type updaterWindowHandle struct {
	win *WebviewWindow
}

func (h *updaterWindowHandle) Show()  { h.win.Show() }
func (h *updaterWindowHandle) Close() { h.win.Close() }
func (h *updaterWindowHandle) EmitEvent(name string, data ...any) bool {
	return h.win.EmitEvent(name, data...)
}

// SetSize implements updater.WindowSizer so the Updater can shrink the
// Up-to-Date state's window to a compact card. WebviewWindow.SetSize is
// fluent (returns the window); the return value is dropped here to keep
// the updater-facing interface free of application-package types.
func (h *updaterWindowHandle) SetSize(width, height int) {
	h.win.SetSize(width, height)
}

// AsUpdaterWindow wraps the receiver as an updater.WindowHandle so it can
// be passed via updater.BYOWindow to Config.Window. Use this when you want
// the updater to drive a webview window you own rather than letting it
// create its own builtin.
//
// The owning window MUST be constructed with AllowSimpleEventEmit set —
// the updater's HTML drives the install/skip/remind/cancel/restart actions
// through the `wails:event:emit:` postMessage path, and that path is gated
// on the field. Without it the buttons silently no-op.
//
//	myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
//	    Title:                "My Updater",
//	    HTML:                 myCustomHTML,
//	    AllowSimpleEventEmit: true,
//	})
//	app.Updater.Init(updater.Config{
//	    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
//	    ...,
//	})
//
// Returned values are independently allocated; create one per Init call.
func (w *WebviewWindow) AsUpdaterWindow() updater.WindowHandle {
	return &updaterWindowHandle{win: w}
}
