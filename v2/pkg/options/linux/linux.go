package linux

// WebviewGpuPolicy values used for determining the webview's hardware acceleration policy.
type WebviewGpuPolicy int

const (
	// WebviewGpuPolicyAlways Hardware acceleration is always enabled.
	WebviewGpuPolicyAlways WebviewGpuPolicy = iota
	// WebviewGpuPolicyOnDemand Hardware acceleration is enabled/disabled as request by web contents.
	WebviewGpuPolicyOnDemand
	// WebviewGpuPolicyNever Hardware acceleration is always disabled.
	WebviewGpuPolicyNever
)

// Options specific to Linux builds
type Options struct {
	// Icon Sets up the icon representing the window. This icon is used when the window is minimized
	// (also known as iconified).
	Icon []byte

	// WindowIsTranslucent sets the window's background to transparent when enabled.
	WindowIsTranslucent bool

	// Messages are messages that can be customised
	Messages *Messages

	// WebviewGpuPolicy used for determining the hardware acceleration policy for the webview.
	//   - WebviewGpuPolicyAlways
	//   - WebviewGpuPolicyOnDemand
	//   - WebviewGpuPolicyNever
	//
	// Due to https://github.com/wailsapp/wails/issues/2977, if options.Linux is nil
	// in the call to wails.Run(), WebviewGpuPolicy is set by default to WebviewGpuPolicyNever.
	// Client code may override this behavior by passing a non-nil Options and set
	// WebviewGpuPolicy as needed.
	WebviewGpuPolicy WebviewGpuPolicy

	// ProgramName is used to set the program's name for the window manager via GTK's g_set_prgname().
	//This name should not be localized. [see the docs]
	//
	//When a .desktop file is created this value helps with window grouping and desktop icons when the .desktop file's Name
	//property differs form the executable's filename.
	//
	//[see the docs]: https://docs.gtk.org/glib/func.set_prgname.html
	ProgramName string
}

type Messages struct {
	WebKit2GTKMinRequired string
}

func DefaultMessages() *Messages {
	return &Messages{
		WebKit2GTKMinRequired: "This application requires at least WebKit2GTK %s to be installed.",
	}
}
