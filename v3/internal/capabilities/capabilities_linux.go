//go:build linux

package capabilities

func NewCapabilities() Capabilities {
	c := Capabilities{}
	// For now, assume Linux has native drag support
	// TODO: Implement proper WebKit version detection
	c.HasNativeDrag = true
	return c
}
