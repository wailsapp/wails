package application

// LinuxOptions contains options for Linux applications.
type LinuxOptions struct {
	// DisableQuitOnLastWindowClosed disables the auto quit of the application if the last window has been closed.
	DisableQuitOnLastWindowClosed bool
}
