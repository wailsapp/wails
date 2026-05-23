//go:build linux && gtk3 && !android && !server

package application

func gtkDispatch(fn func()) {
	InvokeAsync(fn)
}
