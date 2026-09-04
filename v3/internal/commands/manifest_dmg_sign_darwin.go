//go:build darwin

package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const manifestDMGSigningHeadroom = int64(32 << 20)

// prepareManifestDMGForSigning signs the application contained by a generated
// DMG before the outer image is signed or notarized. Generated DMGs are
// compressed and read-only, so work on a private writable conversion and only
// replace the caller's already-staged image after every native operation has
// succeeded.
func prepareManifestDMGForSigning(ctx context.Context, input, identity, entitlements string, hardenedRuntime bool) error {
	if identity == "" {
		return fmt.Errorf("sign DMG application: identity is required")
	}
	workspace, err := os.MkdirTemp(filepath.Dir(input), ".dmg-content-sign-*")
	if err != nil {
		return fmt.Errorf("create DMG signing workspace: %w", err)
	}
	defer os.RemoveAll(workspace)

	writable := filepath.Join(workspace, "writable.dmg")
	if err := runDMGSigningCommand(ctx, workspace, "convert DMG to writable form", "hdiutil", "convert", input, "-format", "UDRW", "-o", writable); err != nil {
		return err
	}
	info, err := os.Stat(writable)
	if err != nil {
		return fmt.Errorf("inspect writable DMG: %w", err)
	}
	megabyte := int64(1 << 20)
	resizedMB := (info.Size() + manifestDMGSigningHeadroom + megabyte - 1) / megabyte
	if err := runDMGSigningCommand(ctx, workspace, "reserve DMG signing space", "hdiutil", "resize", "-size", fmt.Sprintf("%dm", resizedMB), writable); err != nil {
		return err
	}

	mountpoint := filepath.Join(workspace, "mount")
	if err := os.Mkdir(mountpoint, 0o700); err != nil {
		return fmt.Errorf("create DMG signing mountpoint: %w", err)
	}
	if err := runDMGSigningCommand(ctx, workspace, "mount writable DMG", "hdiutil", "attach", "-readwrite", "-nobrowse", "-mountpoint", mountpoint, writable); err != nil {
		return err
	}
	mounted := true
	defer func() {
		if !mounted {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, _ = runManifestCommand(cleanupCtx, workspace, nil, "hdiutil", "detach", mountpoint)
	}()

	base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	app := filepath.Join(mountpoint, base+".app")
	appInfo, err := os.Stat(app)
	if err != nil {
		return fmt.Errorf("find application in DMG: %w", err)
	}
	if !appInfo.IsDir() {
		return fmt.Errorf("application in DMG is not a bundle directory: %s", app)
	}
	signArgs := []string{"--force", "--deep", "--sign", identity}
	if entitlements != "" {
		signArgs = append(signArgs, "--entitlements", entitlements)
	}
	if hardenedRuntime {
		signArgs = append(signArgs, "--options", "runtime")
	}
	signArgs = append(signArgs, app)
	if err := runDMGSigningCommand(ctx, workspace, "sign DMG application", "codesign", signArgs...); err != nil {
		return err
	}
	if err := runDMGSigningCommand(ctx, workspace, "detach signed DMG", "hdiutil", "detach", mountpoint); err != nil {
		return err
	}
	mounted = false

	repacked := filepath.Join(workspace, "signed.dmg")
	if err := runDMGSigningCommand(ctx, workspace, "recompress signed DMG", "hdiutil", "convert", writable, "-format", "UDZO", "-imagekey", "zlib-level=9", "-o", repacked); err != nil {
		return err
	}
	if err := replacePathTransactional(repacked, input); err != nil {
		return fmt.Errorf("publish DMG with signed application: %w", err)
	}
	return nil
}

func runDMGSigningCommand(ctx context.Context, directory, operation, name string, arguments ...string) error {
	output, err := runManifestCommand(ctx, directory, nil, name, arguments...)
	if err == nil {
		return nil
	}
	if output == "" {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return fmt.Errorf("%s: %s: %w", operation, output, err)
}
