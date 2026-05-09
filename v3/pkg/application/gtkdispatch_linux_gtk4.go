//go:build linux && gtk4 && !android && !server

package application

func gtkDispatch(fn func()) {
	go func() {
		defer handlePanic()
		fn()
	}()
}
