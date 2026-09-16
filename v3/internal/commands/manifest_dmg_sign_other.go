//go:build !darwin

package commands

import (
	"context"
	"fmt"
)

func prepareManifestDMGForSigning(context.Context, string, string, string, bool) error {
	return fmt.Errorf("signing an application inside a DMG requires macOS")
}
