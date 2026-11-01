package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) appleSignFlags(signing pipeline.SignSpec, input string) (*flags.Sign, error) {
	config := signing.Config
	certificate, err := signingconfig.ResolvePath(h.root, config, "certificate", config.Certificate)
	if err != nil {
		return nil, fmt.Errorf("darwin.signing.certificate: %w", err)
	}
	entitlements, err := existingOptionalArtifactPathInsideProject(h.root, config.Entitlements)
	if err != nil {
		return nil, fmt.Errorf("darwin.signing.entitlements: %w", err)
	}
	password := ""
	if config.Credential != "" {
		password = environmentValues(mergeEnvironment(os.Environ(), h.environment))[config.Credential]
		if password == "" {
			return nil, fmt.Errorf("Apple signing password environment variable %s is empty", config.Credential)
		}
	}
	return &flags.Sign{Input: input, Certificate: certificate, Password: password, Identity: config.Identity,
		Entitlements: entitlements, HardenedRuntime: config.Notarize, Timestamp: config.TimestampServer}, nil
}

func (h *manifestHandler) signApplePackage(ctx context.Context, spec pipeline.SignSpec) (pipeline.RunResult, error) {
	input, err := existingArtifactPathInsideProject(h.root, spec.Input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	options, err := h.appleSignFlags(spec, input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	notaryKey, err := signingconfig.ResolvePath(h.root, spec.Config, "notarization_key", spec.Config.NotarizationKey)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	options.Notarize, options.KeychainProfile = spec.Config.Notarize, spec.Config.NotarizationCredential
	options.NotaryKey, options.NotaryKeyID, options.NotaryIssuerID = notaryKey, spec.Config.NotarizationKeyID, spec.Config.NotarizationIssuerID
	output := input + ".signed"
	if spec.Format != "dmg" || spec.BundleSigned {
		options.Output = output
		return pipeline.RunResult{}, signPortableApple(ctx, options)
	}
	// Hooks may edit the completed DMG. Preserve their edits and sign the
	// contained application on a private mounted copy before signing the image.
	err = produceNamedPathTransactional(output, ".apple-dmg-sign-*", filepath.Base(input), func(staged string) error {
		if err := copyManifestPath(input, staged); err != nil {
			return err
		}
		if err := prepareManifestDMGWithSigner(ctx, staged, func(app string) error {
			application := *options
			application.Input, application.Output, application.Notarize = app, "", false
			return signPortableApple(ctx, &application)
		}); err != nil {
			return err
		}
		image := *options
		image.Input, image.Output = staged, ""
		return signPortableApple(ctx, &image)
	})
	return pipeline.RunResult{}, err
}
