package build

import (
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
)

// DesktopBuilder builds applications for the desktop
type DesktopBuilder struct {
	*BaseBuilder
}

func newDesktopBuilder(options *Options) *DesktopBuilder {
	return &DesktopBuilder{
		BaseBuilder: NewBaseBuilder(options),
	}
}

// BuildAssets builds the assets for the desktop application
func (d *DesktopBuilder) BuildAssets(options *Options) error {

	// Check assets directory exists
	if !fs.DirExists(options.ProjectData.BuildDir) {
		// Path to default assets
		err := buildassets.Install(options.ProjectData.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
