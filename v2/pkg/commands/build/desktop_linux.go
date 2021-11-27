//go:build linux
// +build linux

package build

// PostCompilation is called after the compilation step, if successful
func (d *DesktopBuilder) PostCompilation(options *Options) error {
	return nil
}
