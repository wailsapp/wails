//go:build linux

package capabilities

import "github.com/wailsapp/wails/v3/internal/operatingsystem"

func NewCapabilities() Capabilities {
	c := Capabilities{}

	webkitVersion := operatingsystem.GetWebkitVersion()
	c.HasNativeDrag = webkitVersion.IsAtLeast(2, 36, 0)
	return c
}
