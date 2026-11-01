package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/packagesign"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) signArch(ctx context.Context, spec pipeline.SignSpec) (pipeline.RunResult, error) {
	keyPath := spec.Config.Certificate
	if !filepath.IsAbs(keyPath) {
		keyPath = filepath.Join(h.root, keyPath)
	}
	password := environmentValues(mergeEnvironment(os.Environ(), h.environment))[spec.Config.Credential]
	if spec.Config.Credential != "" && password == "" {
		return pipeline.RunResult{}, fmt.Errorf("Arch signing credential %s is empty", spec.Config.Credential)
	}
	input, err := existingArtifactPathInsideProject(h.root, spec.Input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	err = signArchPackage(ctx, input, keyPath, password)
	return pipeline.RunResult{Output: spec.Input + ".sig"}, err
}

func signArchPackage(ctx context.Context, input, keyPath, password string) error {
	return producePathTransactional(input+".sig", ".arch-sign-*", func(staged string) error {
		return packagesign.Sign(ctx, input, staged, packagesign.Options{Format: "archlinux", KeyFile: keyPath, Password: password})
	})
}
