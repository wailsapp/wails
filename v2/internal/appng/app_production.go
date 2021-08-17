//go:build production
// +build production

package appng

func (a *App) SetupFlags() {
	a.debug = false
}
