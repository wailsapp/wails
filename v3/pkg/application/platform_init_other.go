//go:build !(linux && purego && !gtk3 && !android && !server)

package application

// platformInit has nothing to load lazily on this platform.
func platformInit() {}
