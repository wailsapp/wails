package commands

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"os"
	"path/filepath"
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
		entitlements, err := existingOptionalArtifactPathInsideProject(h.root, s.BundleSigning.Entitlements)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		if err := manifestSign(&flags.Sign{Input: stagedApp, Identity: s.BundleSigning.Identity, Entitlements: entitlements, HardenedRuntime: s.BundleSigning.Notarize}); err != nil {
			return pipeline.RunResult{}, fmt.Errorf("sign staged DMG application: %w", err)
		}
	}
	options.Out = outputDirectory
	if err := manifestToolPackage(options); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, findAndMovePackage(outputDirectory, s.Project.BinaryName, "dmg", filepath.Join(h.root, s.Output))
}
