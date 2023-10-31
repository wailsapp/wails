//go:build linux && !production

package application

func init() {
	showDevTools = func(wv pointer) {
		windowToggleDevTools(wv)
	}
}
