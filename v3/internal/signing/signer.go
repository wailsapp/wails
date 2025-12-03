package signing

import (
	"context"
	"fmt"
	"io"
	"runtime"
)

// SignerBackend indicates which signing backend is being used
type SignerBackend string

const (
	// BackendNative uses platform-native tools (codesign, signtool)
	BackendNative SignerBackend = "native"
	// BackendRelic uses the relic library for cross-platform signing
	BackendRelic SignerBackend = "relic"
)

// Platform represents a target platform for signing
type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformMacOS   Platform = "darwin"
	PlatformLinux   Platform = "linux"
)

// SignRequest contains all information needed to sign a binary
type SignRequest struct {
	// InputPath is the path to the binary to sign
	InputPath string
	// OutputPath is where to write the signed binary (if empty, signs in place)
	OutputPath string
	// Certificate is the signing certificate configuration
	Certificate CertificateConfig
	// Timestamp enables timestamping with the specified server
	TimestampServer string
	// Description is the application description (for Authenticode)
	Description string
	// URL is the application URL (for Authenticode)
	URL string
	// Platform is the target platform (auto-detected from binary if empty)
	Platform Platform

	// macOS-specific options
	Entitlements    string
	HardenedRuntime bool
	BundleSign      bool // Sign as app bundle vs single binary

	// Verbose enables verbose logging
	Verbose bool
}

// CertificateConfig defines how to load the signing certificate
type CertificateConfig struct {
	// PKCS12Path is the path to a .pfx/.p12 certificate file
	PKCS12Path string
	// PKCS12Password is the password for the PKCS12 file
	PKCS12Password string

	// For macOS native signing: identity name or hash
	Identity string

	// For Windows native signing: certificate thumbprint in Windows cert store
	Thumbprint string

	// For Linux package signing: path to PGP private key file
	PGPKeyPath string
	// PGPKeyPassword is the password for the PGP key (if encrypted)
	PGPKeyPassword string
}

// SignResult contains the result of a signing operation
type SignResult struct {
	// Backend indicates which backend was used
	Backend SignerBackend
	// OutputPath is the path to the signed binary
	OutputPath string
}

// Signer is the interface for code signing implementations
type Signer interface {
	// Sign signs a binary with the given configuration
	Sign(ctx context.Context, req SignRequest) (*SignResult, error)

	// Verify verifies the signature on a binary
	Verify(ctx context.Context, path string) error

	// Backend returns which backend this signer uses
	Backend() SignerBackend

	// Available returns true if this signer is available on the current system
	Available() bool
}

// SignerRegistry holds available signers for each platform
type SignerRegistry struct {
	signers map[Platform][]Signer
}

// NewSignerRegistry creates a new signer registry with auto-detected signers
func NewSignerRegistry() *SignerRegistry {
	r := &SignerRegistry{
		signers: make(map[Platform][]Signer),
	}

	// Register platform-specific signers
	// Order matters: native signers first, then fallbacks
	r.registerSigners()

	return r
}

// registerSigners registers all available signers
func (r *SignerRegistry) registerSigners() {
	// Windows signers
	r.signers[PlatformWindows] = []Signer{
		NewWindowsNativeSigner(),
		NewWindowsRelicSigner(),
	}

	// macOS signers
	r.signers[PlatformMacOS] = []Signer{
		NewMacOSNativeSigner(),
		NewMacOSRelicSigner(),
	}

	// Linux signers (for DEB/RPM packages)
	r.signers[PlatformLinux] = []Signer{
		NewLinuxRelicSigner(),
	}
}

// GetSigner returns the best available signer for the given platform
func (r *SignerRegistry) GetSigner(platform Platform) (Signer, error) {
	signers, ok := r.signers[platform]
	if !ok {
		return nil, fmt.Errorf("no signers registered for platform: %s", platform)
	}

	for _, s := range signers {
		if s.Available() {
			return s, nil
		}
	}

	return nil, fmt.Errorf("no available signer for platform %s (current OS: %s)", platform, runtime.GOOS)
}

// GetSignerWithBackend returns a signer with a specific backend
func (r *SignerRegistry) GetSignerWithBackend(platform Platform, backend SignerBackend) (Signer, error) {
	signers, ok := r.signers[platform]
	if !ok {
		return nil, fmt.Errorf("no signers registered for platform: %s", platform)
	}

	for _, s := range signers {
		if s.Backend() == backend && s.Available() {
			return s, nil
		}
	}

	return nil, fmt.Errorf("signer with backend %s not available for platform %s", backend, platform)
}

// ListAvailableSigners returns all available signers for a platform
func (r *SignerRegistry) ListAvailableSigners(platform Platform) []Signer {
	var available []Signer
	for _, s := range r.signers[platform] {
		if s.Available() {
			available = append(available, s)
		}
	}
	return available
}

// DefaultRegistry is the default signer registry
var DefaultRegistry = NewSignerRegistry()

// Sign signs a binary using the best available signer
func SignBinary(ctx context.Context, req SignRequest) (*SignResult, error) {
	platform := req.Platform
	if platform == "" {
		// Auto-detect from file extension or current platform
		platform = detectPlatform(req.InputPath)
	}

	signer, err := DefaultRegistry.GetSigner(platform)
	if err != nil {
		return nil, err
	}

	return signer.Sign(ctx, req)
}

// detectPlatform detects the target platform from a file path
func detectPlatform(path string) Platform {
	// Simple detection based on extension
	switch {
	case hasExtension(path, ".exe", ".dll", ".msi", ".msix", ".appx", ".ps1"):
		return PlatformWindows
	case hasExtension(path, ".app", ".dmg", ".pkg"):
		return PlatformMacOS
	case hasExtension(path, ".deb", ".rpm"):
		return PlatformLinux
	default:
		// Fall back to current platform
		return Platform(runtime.GOOS)
	}
}

func hasExtension(path string, exts ...string) bool {
	for _, ext := range exts {
		if len(path) > len(ext) && path[len(path)-len(ext):] == ext {
			return true
		}
	}
	return false
}

// ProgressCallback is called during long operations to report progress
type ProgressCallback func(current, total int64, message string)

// SignWithProgress signs a binary and reports progress
func SignWithProgress(ctx context.Context, req SignRequest, progress ProgressCallback) (*SignResult, error) {
	// For now, just call Sign - progress can be added later
	return SignBinary(ctx, req)
}

// ReadSeekCloser combines io.Reader, io.Seeker, and io.Closer
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
