package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/buildmeta"
	"github.com/wailsapp/wails/v3/internal/voidpackage"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) packageVoid(ctx context.Context, spec pipeline.PackageSpec) (pipeline.RunResult, error) {
	version, err := buildmeta.Parse(spec.Project.Version)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	options := voidpackage.Options{Metadata: voidpackage.Metadata{Name: spec.Project.BinaryName, Version: strings.ReplaceAll(version.Version, "-", "~"), SourceVersion: version.Version, Arch: spec.TargetArch, Description: spec.Project.Description, Maintainer: spec.Config.Maintainer, Depends: spec.Config.Dependencies}, Contents: []voidpackage.Entry{{Source: filepath.Join(h.root, spec.Binary), Destination: "/usr/bin/" + spec.Project.BinaryName, Mode: 0755}}}
	for _, binary := range spec.Companions {
		options.Contents = append(options.Contents, voidpackage.Entry{Source: filepath.Join(h.root, binary.Path), Destination: "/usr/bin/" + binary.Name, Mode: 0755})
	}
	for _, item := range []struct{ source, destination string }{
		{filepath.Join(h.root, spec.Assets, "appicon.png"), "/usr/share/icons/hicolor/256x256/apps/" + spec.Project.BinaryName + ".png"},
		{filepath.Join(h.root, spec.Assets, "linux", spec.Project.BinaryName+".desktop"), "/usr/share/applications/" + spec.Project.BinaryName + ".desktop"},
	} {
		if _, err := os.Stat(item.source); err == nil {
			options.Contents = append(options.Contents, voidpackage.Entry{Source: item.source, Destination: item.destination, Mode: 0644})
		}
	}
	options.Install, err = h.voidLifecycleScript(spec.Config.PreInstall, spec.Config.PostInstall)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	options.Remove, err = h.voidLifecycleScript(spec.Config.PreRemove, spec.Config.PostRemove)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	err = producePathTransactional(filepath.Join(h.root, spec.Output), ".void-output-*", func(staged string) error {
		output, err := os.Create(staged)
		if err != nil {
			return err
		}
		writeErr := voidpackage.Write(ctx, output, options)
		closeErr := output.Close()
		if writeErr != nil {
			return writeErr
		}
		return closeErr
	})
	return pipeline.RunResult{}, err
}

func (h *manifestHandler) voidLifecycleScript(before, after string) ([]byte, error) {
	if before == "" && after == "" {
		return nil, nil
	}
	var script strings.Builder
	script.WriteString("#!/bin/sh\nset -e\ncase \"$1\" in\n")
	for _, phase := range []struct{ name, path string }{{"pre", before}, {"post", after}} {
		if phase.path == "" {
			continue
		}
		path, err := existingPathInsideProject(h.root, phase.path)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&script, "%s)\n(\n%s\n)\n;;\n", phase.name, data)
	}
	script.WriteString("esac\n")
	return []byte(script.String()), nil
}
