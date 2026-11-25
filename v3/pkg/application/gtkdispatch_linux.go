//go:build linux && !gtk3 && !android && !openharmony && !server

package application

func gtkDispatch(fn func()) {
	go func() {
		defer handlePanic()
		fn()
	}()
}
