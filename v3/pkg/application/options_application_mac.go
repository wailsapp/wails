package application

// ActivationPolicy is the activation policy for the application.
type ActivationPolicy int

const (
	// ActivationPolicyRegular is used for applications that have a user interface,
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory is used for applications that do not have a main window,
	// such as system tray applications or background applications.
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

// MacOptions contains options for macOS applications.
type MacOptions struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
	// If set to true, the application will terminate when the last window is closed.
	ApplicationShouldTerminateAfterLastWindowClosed bool
}
