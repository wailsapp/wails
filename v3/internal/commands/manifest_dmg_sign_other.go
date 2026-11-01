//go:build !darwin

package commands

import (
	"context"
	"fmt"
)

func prepareManifestDMGForSigning(context.Context, string, string, string, bool) error {
	return fmt.Errorf("signing an application inside a DMG requires macOS")
}

func prepareManifestDMGWithSigner(context.Context, string, func(string) error) error {
	return fmt.Errorf("signing an application inside a DMG after package hooks requires macOS")
}
