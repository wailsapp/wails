package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) packageDMG(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	options, err := h.dmgOptions(s)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	outputDirectory := filepath.Join(h.packageWorkspace(s), "output")
	if err := removeGenerated(outputDirectory); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(outputDirectory, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	appBundle, err := existingArtifactPathInsideProject(h.root, s.Binary)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("DMG application bundle: %w", err)
	}
	stagedApp := filepath.Join(outputDirectory, s.Project.BinaryName+".app")
	if err := copyManifestPath(appBundle, stagedApp); err != nil {
		return pipeline.RunResult{}, fmt.Errorf("stage DMG application bundle: %w", err)
	}
	if s.BundleSigning != nil {
		signOptions, err := h.appleSignFlags(pipeline.SignSpec{Config: *s.BundleSigning}, stagedApp)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		if portableAppleSigning(signOptions) {
			err = signPortableApple(ctx, signOptions)
		} else {
			err = manifestSign(signOptions)
		}
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("sign staged DMG application: %w", err)
		}
	}
	options.Out = outputDirectory
	if err := manifestToolPackage(options); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, findAndMovePackage(outputDirectory, s.Project.BinaryName, "dmg", filepath.Join(h.root, s.Output))
}
