package commands

import (
	"context"

	"github.com/wailsapp/wails/v3/internal/packagesign"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) signArch(ctx context.Context, spec pipeline.SignSpec) (pipeline.RunResult, error) {
	keyPath, err := signingconfig.ResolvePath(h.root, spec.Config, "certificate", spec.Config.Certificate)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	password, err := manifestSigningPassword(spec.Config, "linux", h.environment)
	if err != nil {
		return pipeline.RunResult{}, err
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
