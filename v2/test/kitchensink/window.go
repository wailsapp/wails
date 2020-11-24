package main

import (
	wails "github.com/wailsapp/wails/v2"
)

// Window struct
type Window struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (w *Window) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	w.runtime = runtime
	return nil
}

func (w *Window) SetTitle(title string) {
	w.runtime.Window.SetTitle(title)
}

func (w *Window) Fullscreen() {
	w.runtime.Window.Fullscreen()
}

func (w *Window) UnFullscreen() {
	w.runtime.Window.UnFullscreen()
}

func (w *Window) Maximise() {
	w.runtime.Window.Maximise()
}
func (w *Window) Unmaximise() {
	w.runtime.Window.Unmaximise()
}
func (w *Window) Minimise() {
	w.runtime.Window.Minimise()
}
func (w *Window) Unminimise() {
	w.runtime.Window.Unminimise()
}
func (w *Window) Center() {
	w.runtime.Window.Center()
}
func (w *Window) Show() {
	w.runtime.Window.Show()
}
func (w *Window) Hide() {
	w.runtime.Window.Hide()
}
func (w *Window) SetSize(width int, height int) {
	w.runtime.Window.SetSize(width, height)
}
func (w *Window) SetPosition(x int, y int) {
	w.runtime.Window.SetPosition(x, y)
}
func (w *Window) Close() {
	w.runtime.Window.Close()
}
