package commands

import (
	"context"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/keychain"
	"github.com/wailsapp/wails/v3/internal/packagesign"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func signExportedPackage(options *flags.Sign, format, key string) error {
	password := options.PGPPassword
	if password == "" {
		password, _ = keychain.Get(keychain.KeyPGPPassword)
	}
	output := options.Output
	if output == "" {
		output = options.Input
	}
	return produceNamedPathTransactional(output, ".package-sign-*", filepath.Base(output), func(staged string) error {
		return packagesign.Sign(context.Background(), options.Input, staged, packagesign.Options{Format: format, KeyFile: key, Password: password, Role: options.Role})
	})
}

func (h *manifestHandler) signLinuxPackage(ctx context.Context, spec pipeline.SignSpec, key string) (pipeline.RunResult, error) {
	password, err := manifestSigningPassword(spec.Config, "linux", h.environment)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	input, err := existingArtifactPathInsideProject(h.root, spec.Input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	output := filepath.Join(h.root, spec.Input+".signed")
	err = produceNamedPathTransactional(output, ".linux-sign-output-*", filepath.Base(input), func(staged string) error {
		return packagesign.Sign(ctx, input, staged, packagesign.Options{Format: spec.Format, KeyFile: key, Password: password, Role: spec.Config.Identity})
	})
	return pipeline.RunResult{}, err
}
