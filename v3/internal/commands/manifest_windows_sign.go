package commands

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) signWindowsPackage(ctx context.Context, spec pipeline.SignSpec) (pipeline.RunResult, error) {
	input, err := existingArtifactPathInsideProject(h.root, spec.Input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	certificate, err := signingconfig.ResolvePath(h.root, spec.Config, "certificate", spec.Config.Certificate)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("windows.signing.certificate: %w", err)
	}
	password, err := manifestSigningPassword(spec.Config, "windows", h.environment)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	output := filepath.Join(h.root, spec.Input+".signed")
	options := &flags.Sign{Input: input, Output: output, Certificate: certificate, Password: password, Thumbprint: spec.Config.Thumbprint, Timestamp: spec.Config.TimestampServer, Identity: spec.Config.Identity}
	err = manifestWindowsSign(ctx, options)
	return pipeline.RunResult{}, err
}
