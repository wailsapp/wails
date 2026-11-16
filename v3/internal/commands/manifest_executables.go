package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func executablePackage(value string) string {
	if value == "" || value == "." {
		return "."
	}
	return "./" + strings.TrimPrefix(value, "./")
}

func (h *manifestHandler) compressExecutable(ctx context.Context, spec pipeline.CompileSpec, staged string, detail *string) error {
	if spec.UPXLevel == 0 || spec.TargetOS == "darwin" {
		return nil
	}
	output, err := runManifestCommand(ctx, h.root, h.environment, "upx", fmt.Sprintf("-%d", spec.UPXLevel), staged)
	*detail += output
	if err != nil {
		return fmt.Errorf("compress executable: %w", err)
	}
	return nil
}
