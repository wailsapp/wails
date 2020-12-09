package build

import (
	"github.com/wailsapp/wails/v2/internal/project"
)

// HybridBuilder builds applications as a server
type HybridBuilder struct {
	*BaseBuilder
	desktop *DesktopBuilder
	server  *ServerBuilder
}

func newHybridBuilder() Builder {
	result := &HybridBuilder{
		BaseBuilder: NewBaseBuilder(),
		desktop:     newDesktopBuilder(),
		server:      newServerBuilder(),
	}
	return result
}

// BuildAssets builds the assets for the desktop application
func (b *HybridBuilder) BuildAssets(options *Options) error {
	var err error

	// Build base assets (HTML/JS/CSS/etc)
	err = b.BuildBaseAssets(options)
	if err != nil {
		return err
	}
	// Build static assets
	err = b.buildStaticAssets(b.projectData)
	if err != nil {
		return err
	}

	return nil
}

// BuildAssets builds the assets for the desktop application
func (b *HybridBuilder) BuildBaseAssets(options *Options) error {

	assets, err := b.BaseBuilder.ExtractAssets()
	if err != nil {
		return err
	}

	err = b.desktop.BuildBaseAssets(assets, options)
	if err != nil {
		return err
	}

	err = b.server.BuildBaseAssets(assets)
	if err != nil {
		return err
	}

	// Build desktop static assets
	err = b.desktop.buildStaticAssets(b.projectData)
	if err != nil {
		return err
	}

	// Build server static assets
	err = b.server.buildStaticAssets(b.projectData)
	if err != nil {
		return err
	}

	return nil
}

func (b *HybridBuilder) BuildRuntime(options *Options) error {
	err := b.desktop.BuildRuntime(options)
	if err != nil {
		return err
	}

	err = b.server.BuildRuntime(options)
	if err != nil {
		return err
	}

	return nil
}

func (b *HybridBuilder) SetProjectData(projectData *project.Project) {
	b.BaseBuilder.SetProjectData(projectData)
	b.desktop.SetProjectData(projectData)
	b.server.SetProjectData(projectData)
}

func (b *HybridBuilder) CompileProject(options *Options) error {
	return b.BaseBuilder.CompileProject(options)
}

func (b *HybridBuilder) CleanUp() {
	b.desktop.CleanUp()
	b.server.CleanUp()
}
