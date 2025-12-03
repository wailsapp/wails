//go:build !darwin

package signing

import (
	"context"
	"fmt"
)

// MacOSNativeSigner uses codesign for signing macOS binaries
// This is a stub for non-darwin platforms
type MacOSNativeSigner struct{}

// NewMacOSNativeSigner creates a new macOS native signer
func NewMacOSNativeSigner() *MacOSNativeSigner {
	return &MacOSNativeSigner{}
}

// Backend returns the signer backend type
func (s *MacOSNativeSigner) Backend() SignerBackend {
	return BackendNative
}

// Available returns false on non-darwin platforms
func (s *MacOSNativeSigner) Available() bool {
	return false
}

// Sign is not available on non-darwin platforms
func (s *MacOSNativeSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	return nil, fmt.Errorf("native macOS signing requires macOS (codesign not available)")
}

// Verify is not available on non-darwin platforms
func (s *MacOSNativeSigner) Verify(ctx context.Context, path string) error {
	return fmt.Errorf("native macOS signature verification requires macOS")
}
