//go:build darwin

package capabilities

func newCapabilities(_ string) Capabilities {
	c := Capabilities{}
	c.HasNativeDrag = false
	return c
}
