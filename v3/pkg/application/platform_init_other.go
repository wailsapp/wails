//go:build !(linux && purego && !gtk3 && !android && !server) && !(windows && !server)

package application

// platformInit has nothing to load lazily on this platform.
func platformInit() {}
