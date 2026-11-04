package commands

import (
	"context"
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/applesign"
	"github.com/wailsapp/wails/v3/internal/codesignature"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func portableAppleSigning(options *flags.Sign) bool {
	return options.Certificate != "" || options.Identity == "-"
}

func appleBundleInput(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".app", ".framework", ".xpc", ".bundle", ".appex":
		return true
	default:
		return false
	}
}

// Detect the input independently of the host so Windows defaults cannot
// select a certificate for an Apple artifact.
func appleSigningInput(path string) bool {
	if appleBundleInput(path) || strings.EqualFold(filepath.Ext(path), ".dmg") {
		return true
	}
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	var header [4]byte
	if _, err := io.ReadFull(file, header[:]); err != nil {
		return false
	}
	for _, magic := range []uint32{binary.LittleEndian.Uint32(header[:]), binary.BigEndian.Uint32(header[:])} {
		if magic == macho.Magic32 || magic == macho.Magic64 || magic == macho.MagicFat || magic == 0xcafebabf {
			return true
		}
	}
	return false
}

func appleSigningOptions(options *flags.Sign) (applesign.Options, error) {
	result := applesign.Options{}
	if options.Certificate != "" {
		data, err := readSigningFile(options.Certificate, 32<<20)
		if err != nil {
			return result, err
		}
		defer clear(data)
		identity, err := codesignature.DecodePFX(data, options.Password)
		if err != nil {
			return result, fmt.Errorf("Apple signing identity: %w", err)
		}
		result, err = applesign.SigningOptions(identity, codesignature.Options{TimestampURL: options.Timestamp})
		if err != nil {
			return result, err
		}
	} else if options.Identity != "-" {
		return result, fmt.Errorf("portable Apple signing requires an exported certificate or identity '-'")
	}
	result.Identifier = filepath.Base(options.Input)
	result.HardenedRuntime = options.HardenedRuntime || options.Notarize
	if options.Entitlements != "" {
		data, err := readSigningFile(options.Entitlements, 1<<20)
		if err != nil {
			return result, err
		}
		result.Entitlements = data
	}
	return result, nil
}

func signAppleArtifact(ctx context.Context, input, output string, options applesign.Options) error {
	switch strings.ToLower(filepath.Ext(input)) {
	case ".app", ".framework", ".xpc", ".bundle", ".appex":
		return applesign.SignBundle(ctx, input, output, options)
	case ".dmg":
		// Runtime and entitlements belong to the application inside the image.
		options.HardenedRuntime, options.Entitlements = false, nil
		_, err := applesign.SignDMG(ctx, input, output, options)
		return err
	default:
		_, err := applesign.SignMachO(ctx, input, output, options)
		return err
	}
}

func signPortableApple(ctx context.Context, options *flags.Sign) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if options.Notarize {
		ext := strings.ToLower(filepath.Ext(options.Input))
		if ext != ".app" && ext != ".dmg" {
			return fmt.Errorf("notarization and stapling require an app bundle or DMG")
		}
	}
	if options.Notarize && options.Certificate == "" {
		return fmt.Errorf("notarization requires a Developer ID certificate; ad-hoc signing is insufficient")
	}
	if options.Notarize && options.Timestamp == "" {
		return fmt.Errorf("notarized portable signing requires an RFC 3161 timestamp server")
	}
	signing, err := appleSigningOptions(options)
	if err != nil {
		return err
	}
	output := options.Output
	if output == "" {
		output = options.Input
	}
	return produceNamedPathTransactional(output, ".apple-sign-output-*", filepath.Base(options.Input), func(staged string) error {
		if err := signAppleArtifact(ctx, options.Input, staged, signing); err != nil {
			return err
		}
		if !options.Notarize {
			return nil
		}
		request := *options
		request.Input, request.Output = staged, ""
		if usesNotaryAPI(&request) {
			return notarizeMacOSWithAPI(ctx, &request)
		}
		if strings.EqualFold(filepath.Ext(staged), ".app") {
			return notarizeMacOSApp(ctx, &request)
		}
		return notarizeMacOSDiskImage(ctx, &request)
	})
}
