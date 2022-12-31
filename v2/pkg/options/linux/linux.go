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
	WebviewGpuPolicy WebviewGpuPolicy
}

type Messages struct {
	WebKit2GTKMinRequired string
}

func DefaultMessages() *Messages {
	return &Messages{
		WebKit2GTKMinRequired: "This application requires at least WebKit2GTK %s to be installed.",
	}
}
