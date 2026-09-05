//go:build linux && gtk3

package capabilities

func NewCapabilities() Capabilities {
	return Capabilities{
		HasNativeDrag: true,
		GTKVersion:    3,
		WebKitVersion: "4.1",
	}
}
