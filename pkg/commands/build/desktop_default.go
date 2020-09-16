// +build !linux

package build

// This is used when there is no compilation to be done for the asset
func (d *DesktopBuilder) compileIcon(assetDir string, iconFile string) error {
	return nil
}
