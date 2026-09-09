//go:build linux && !gtk3

package capabilities

func NewCapabilities() Capabilities {
	return Capabilities{
		HasNativeDrag: true,
		GTKVersion:    4,
		WebKitVersion: "6.0",
	}
}
