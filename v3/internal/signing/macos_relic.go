package signing

import (
	"context"
	"fmt"
)

// MacOSRelicSigner uses the relic library for cross-platform macOS signing
// Note: Relic's macOS/Mach-O support is preliminary and has significant limitations.
// For full cross-platform macOS signing, consider using rcodesign instead.
type MacOSRelicSigner struct{}

// NewMacOSRelicSigner creates a new relic-based macOS signer
func NewMacOSRelicSigner() *MacOSRelicSigner {
	return &MacOSRelicSigner{}
}

// Backend returns the signer backend type
func (s *MacOSRelicSigner) Backend() SignerBackend {
	return BackendRelic
}

// Available returns true - relic is available as a Go library
// However, its macOS signing support is limited
func (s *MacOSRelicSigner) Available() bool {
	// Return true to allow fallback, but Sign will return a helpful error
	return true
}

// Sign signs a macOS binary or app bundle using the relic library
// Note: Relic's Mach-O support is preliminary. For production use,
// consider using native codesign on macOS or rcodesign for cross-platform signing.
func (s *MacOSRelicSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	// Relic's Mach-O signing has significant limitations:
	// - No fat/universal binary support
	// - No ticket stapling
	// - Preliminary implementation
	//
	// Instead of providing a broken implementation, we return a helpful error
	// directing users to better alternatives.

	return nil, fmt.Errorf(`cross-platform macOS signing via relic is not fully supported

Relic's Mach-O signing has limitations:
- No fat/universal binary support
- No notarization ticket stapling
- Preliminary implementation

Recommended alternatives:

1. Native codesign (macOS only):
   wails3 sign --app MyApp.app

2. rcodesign (cross-platform, recommended):
   cargo install apple-codesign
   rcodesign sign --p12-file cert.p12 --code-signature-flags runtime MyApp.app

3. Use CI/CD with macOS runners for production signing

For more information, see: https://wails.io/docs/guides/signing`)
}

// Verify verifies the signature on a macOS binary
func (s *MacOSRelicSigner) Verify(ctx context.Context, path string) error {
	return fmt.Errorf(`cross-platform macOS signature verification is not supported

Use native codesign on macOS:
  codesign --verify --verbose=2 %s

Or spctl for notarization verification:
  spctl --assess --verbose=2 %s`, path, path)
}
