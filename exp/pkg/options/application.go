package options

type Application struct {
	Mac *Mac
}

// Mac contains macOS specific options

type ActivationPolicy int

const (
	ActivationPolicyRegular ActivationPolicy = iota
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

type Mac struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
}

type Window struct {
	Title         string
	Width, Height int
	AlwaysOnTop   bool
	URL           string
	DisableResize bool
	Resizable     bool
	MinWidth      int
	MinHeight     int
	MaxWidth      int
	MaxHeight     int
}
