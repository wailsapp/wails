//go:build !production && wails_native

package application

func (a *App) preRun() error { return nil }

func (a *App) postQuit() {}

func (a *App) enableDevTools() {}
