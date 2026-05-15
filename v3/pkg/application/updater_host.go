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
	wopts := WebviewWindowOptions{
		Title:         opts.Title,
		Width:         opts.Width,
		Height:        opts.Height,
		Frameless:     opts.Frameless,
		AlwaysOnTop:   opts.AlwaysOnTop,
		DisableResize: opts.DisableResize,
	}
	win := h.app.Window.NewWithOptions(wopts)
	if opts.InitialHTML != "" {
		win.SetHTML(opts.InitialHTML)
	}
	win.Show()
	return &updaterWindowHandle{win: win}
}

// updaterWindowHandle bridges *WebviewWindow into the updater.WindowHandle
// interface. The wrapper drops the fluent return values so the smaller
// updater interface stays free of application-package types.
type updaterWindowHandle struct {
	win *WebviewWindow
}

func (h *updaterWindowHandle) SetHTML(s string) { h.win.SetHTML(s) }
func (h *updaterWindowHandle) Show()            { h.win.Show() }
func (h *updaterWindowHandle) Close()           { h.win.Close() }
func (h *updaterWindowHandle) EmitEvent(name string, data ...any) bool {
	return h.win.EmitEvent(name, data...)
}
