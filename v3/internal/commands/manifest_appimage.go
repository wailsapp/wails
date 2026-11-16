package commands

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/appimage"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) appImageOptions(spec pipeline.PackageSpec) (appimage.Options, error) {
	options := appimage.Options{
		Root: h.root, Name: spec.Project.BinaryName, Arch: spec.TargetArch,
		Output:   filepath.Join(h.root, spec.Output),
		CacheDir: filepath.Join(h.root, ".wails", "tools", "appimage"),
		WebKit:   packageUsesWebView(spec), GTKVersion: spec.GTKVersion,
		UpdateInfo: spec.Config.UpdateInfo,
		Binaries:   []appimage.Binary{{Name: spec.Project.BinaryName, Path: filepath.Join(h.root, spec.Binary)}},
	}
	if spec.Config.GTKVersion != 0 {
		options.GTKVersion = spec.Config.GTKVersion
	}
	for _, companion := range spec.Companions {
		options.Binaries = append(options.Binaries, appimage.Binary{Name: companion.Name, Path: filepath.Join(h.root, companion.Path)})
	}
	for _, field := range []struct {
		name, value string
		destination *string
	}{
		{"runtime_root", spec.Config.RuntimeRoot, &options.RuntimeRoot},
		{"runtime_path", spec.Config.RuntimePath, &options.RuntimePath},
	} {
		if field.value == "" {
			continue
		}
		path, err := manifest.ResolveProjectPath(h.root, "packages.appimage."+field.name, field.value, true)
		if err != nil {
			return appimage.Options{}, err
		}
		*field.destination = path
	}
	return options, nil
}
