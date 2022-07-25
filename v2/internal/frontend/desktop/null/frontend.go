package null

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// Frontend implements an empty Frontend that simply waits until the context is done.
type Frontend struct {
	// Context
	ctx  context.Context
	done bool
}

// NewFrontend returns an initialized Frontend
func NewFrontend(ctx context.Context) *Frontend {
	return &Frontend{
		ctx: ctx,
	}
}

// Show does nothing
func (f *Frontend) Show() {

}

// Hide does nothing
func (f *Frontend) Hide() {

}

// ScreenGetAll returns an empty slice
func (f *Frontend) ScreenGetAll() ([]frontend.Screen, error) {
	return []frontend.Screen{}, nil
}

// WindowSetBackgroundColour does nothing
func (f *Frontend) WindowSetBackgroundColour(col *options.RGBA) {

}

// WindowReload does nothing
func (f *Frontend) WindowReload() {

}

// WindowReloadApp does nothing
func (f *Frontend) WindowReloadApp() {

}

// WindowSetAlwaysOnTop does nothing
func (f *Frontend) WindowSetAlwaysOnTop(b bool) {

}

// WindowSetSystemDefaultTheme does nothing
func (f *Frontend) WindowSetSystemDefaultTheme() {

}

// WindowSetLightTheme does nothing
func (f *Frontend) WindowSetLightTheme() {

}

// WindowSetDarkTheme does nothing
func (f *Frontend) WindowSetDarkTheme() {

}

// Run waits until the context is done and then exits
func (f *Frontend) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			time.Sleep(1 * time.Millisecond)
		}
		if f.done {
			break
		}
	}
	return nil
}

// WindowCenter does nothing
func (f *Frontend) WindowCenter() {

}

// WindowSetPosition does nothing
func (f *Frontend) WindowSetPosition(x, y int) {

}

// WindowGetPosition does nothing
func (f *Frontend) WindowGetPosition() (int, int) {
	return 0, 0
}

// WindowSetSize does nothing
func (f *Frontend) WindowSetSize(width, height int) {

}

// WindowGetSize does nothing
func (f *Frontend) WindowGetSize() (int, int) {
	return 0, 0
}

// WindowSetTitle does nothing
func (f *Frontend) WindowSetTitle(title string) {

}

// WindowFullscreen does nothing
func (f *Frontend) WindowFullscreen() {

}

// WindowUnfullscreen does nothing
func (f *Frontend) WindowUnfullscreen() {

}

// WindowShow does nothing
func (f *Frontend) WindowShow() {

}

// WindowHide does nothing
func (f *Frontend) WindowHide() {

}

// WindowMaximize does nothing
func (f *Frontend) WindowMaximise() {

}

// WindowToggleMaximise does nothing
func (f *Frontend) WindowToggleMaximise() {

}

// WindowUnmaximise does nothing
func (f *Frontend) WindowUnmaximise() {
}

// WindowMinimise does nothing
func (f *Frontend) WindowMinimise() {

}

// WindowUnminimise does nothing
func (f *Frontend) WindowUnminimise() {

}

// WindowSetMinSize does nothing
func (f *Frontend) WindowSetMinSize(width int, height int) {

}

// WindowSetMaxSize does nothing
func (f *Frontend) WindowSetMaxSize(width int, height int) {

}

// WindowSetRGBA does nothing
func (f *Frontend) WindowSetRGBA(col *options.RGBA) {
}

// Quit does nothing
func (f *Frontend) Quit() {
	f.done = true
}

// Notify does nothing
func (f *Frontend) Notify(name string, data ...interface{}) {

}

// Callback does nothing
func (f *Frontend) Callback(message string) {

}

// startDrag does nothing
func (f *Frontend) startDrag() {

}

// ExecJS does nothing
func (f *Frontend) ExecJS(js string) {

}
