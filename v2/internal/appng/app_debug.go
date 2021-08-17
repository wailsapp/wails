//go:build debug
// +build debug

package appng

func (a *App) SetupFlags() {
	a.debug = true
}
