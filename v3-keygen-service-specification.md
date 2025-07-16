# Wails v3 Keygen.sh Service Integration Specification

## Overview

This specification outlines the integration of Keygen.sh licensing and auto-update functionality into Wails v3 as a native service. The integration leverages the official Keygen Go SDK to provide seamless license validation, update checking, and artifact distribution for Wails applications.

## Goals

1. **Simplified Auto-Updates**: Enable Wails developers to easily implement automatic updates through Keygen.sh's distribution API
2. **License Management**: Provide built-in license validation and enforcement
3. **Offline Support**: Support offline license validation and caching
4. **Cross-Platform**: Maintain Wails' cross-platform capabilities (Windows, macOS, Linux)
5. **Developer-Friendly**: Simple API that follows Wails v3 service conventions

## Architecture

### Service Structure

The Keygen service will follow the established Wails v3 service pattern:

```
v3/pkg/services/keygen/
├── keygen.go           # Main service implementation
├── keygen_darwin.go    # macOS-specific implementation
├── keygen_windows.go   # Windows-specific implementation
├── keygen_linux.go     # Linux-specific implementation
├── types.go            # Type definitions
└── README.md           # Service documentation
```

### Core Components

#### 1. Keygen Service

```go
package keygen

import (
    "github.com/keygen-sh/keygen-go/v3"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
    impl     platformKeygen
    client   *keygen.Client
    config   Config
    state    *LicenseState
}

type Config struct {
    Account      string
    Product      string
    LicenseKey   string
    PublicKey    string
    CacheDir     string
    AutoCheck    bool
    CheckInterval time.Duration
}

type LicenseState struct {
    Valid        bool
    License      *keygen.License
    LastChecked  time.Time
    OfflineMode  bool
    Entitlements map[string]interface{}
}
```

#### 2. Platform Interface

```go
type platformKeygen interface {
    Startup(ctx context.Context, options application.ServiceOptions) error
    Shutdown() error
    InstallUpdate(release *keygen.Release) error
    GetMachineFingerprint() (string, error)
    GetInstallPath() string
}
```

#### 3. Service Methods (Frontend-Exposed)

```go
// License Management
func (s *Service) ValidateLicense() (*LicenseValidationResult, error)
func (s *Service) ActivateMachine() (*MachineActivationResult, error)
func (s *Service) DeactivateMachine() error
func (s *Service) GetLicenseInfo() (*LicenseInfo, error)
func (s *Service) CheckEntitlement(feature string) (bool, error)

// Update Management
func (s *Service) CheckForUpdates() (*UpdateInfo, error)
func (s *Service) DownloadUpdate(releaseID string) (*DownloadProgress, error)
func (s *Service) InstallUpdate() error
func (s *Service) GetCurrentVersion() string
func (s *Service) SetUpdateChannel(channel string) error

// Offline Support
func (s *Service) SaveOfflineLicense() error
func (s *Service) LoadOfflineLicense() error
func (s *Service) ClearLicenseCache() error
```

### Update Flow

1. **Check for Updates**
   - Query Keygen API for latest release
   - Compare with current version
   - Respect update channels (stable, beta, etc.)
   - Verify license is valid for update

2. **Download Update**
   - Stream download with progress reporting
   - Verify cryptographic signatures
   - Store in temporary location

3. **Install Update**
   - Platform-specific installation
   - macOS: Replace .app bundle
   - Windows: Run installer/replace exe
   - Linux: Replace AppImage/binary

### Event System

The service will emit events for frontend notification:

```go
type UpdateAvailableEvent struct {
    Version     string
    ReleaseNotes string
    Critical    bool
}

type LicenseStatusEvent struct {
    Valid       bool
    ExpiresAt   *time.Time
    Message     string
}

type DownloadProgressEvent struct {
    BytesDownloaded int64
    TotalBytes      int64
    Percentage      float64
}
```

## Implementation Details

### Service Initialization

```go
// In application setup
keygenService := keygen.NewService(&keygen.Config{
    Account:    "your-account-id",
    Product:    "your-product-id",
    PublicKey:  "your-ed25519-public-key",
    AutoCheck:  true,
    CheckInterval: 6 * time.Hour,
})

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(keygenService),
    },
})
```

### Frontend Usage

```javascript
// Check license validity
const result = await KeygenService.ValidateLicense();
if (!result.valid) {
    showLicenseDialog(result.message);
}

// Check for updates
const updateInfo = await KeygenService.CheckForUpdates();
if (updateInfo.available) {
    const confirmed = await showUpdateDialog(updateInfo);
    if (confirmed) {
        const progress = await KeygenService.DownloadUpdate(updateInfo.releaseId);
        // Monitor progress...
        await KeygenService.InstallUpdate();
    }
}
```

### Configuration Options

The service supports configuration through:

1. **Service Options**: Passed during service creation
2. **Environment Variables**: 
   - `KEYGEN_ACCOUNT`
   - `KEYGEN_PRODUCT`
   - `KEYGEN_LICENSE_KEY`
3. **Configuration File**: `keygen.json` in app data directory

### Security Considerations

1. **Public Key Storage**: Ed25519 public key embedded in binary
2. **License Key Protection**: Stored in platform-specific secure storage
3. **Update Verification**: All updates verified with cryptographic signatures
4. **HTTPS Only**: All API communication over HTTPS
5. **Machine Fingerprinting**: Uses hardware identifiers for device-specific licenses

### Error Handling

The service provides detailed error types:

```go
type KeygenError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

// Common error codes
const (
    ErrLicenseInvalid     = "LICENSE_INVALID"
    ErrLicenseExpired     = "LICENSE_EXPIRED"
    ErrMachineLimitReached = "MACHINE_LIMIT_REACHED"
    ErrUpdateFailed       = "UPDATE_FAILED"
    ErrNetworkError       = "NETWORK_ERROR"
)
```

## Platform-Specific Considerations

### macOS
- Use `NSBundle` for version detection
- Implement update via app bundle replacement
- Store licenses in Keychain
- Handle app translocation and notarization

### Windows
- Use Windows Registry for persistent storage
- Implement update via MSI installer or in-place replacement
- Handle UAC elevation for updates
- Support both portable and installed apps

### Linux
- Support AppImage, Flatpak, and binary distributions
- Use XDG directories for storage
- Implement update based on distribution type
- Handle different permission models

## Migration Path

For existing Wails applications:

1. **Add Service**: Register the Keygen service
2. **Configure**: Set account and product details
3. **Implement UI**: Add license validation and update UI
4. **Test**: Verify update flow on all platforms

## Future Enhancements

1. **Update Rollback**: Ability to revert to previous version
2. **Delta Updates**: Only download changed files
3. **Background Updates**: Silent update downloads
4. **A/B Testing**: Support for gradual rollouts
5. **Analytics Integration**: Usage and crash reporting

## Example Implementation

A complete example will be provided in `v3/examples/keygen-integration/` demonstrating:
- License validation on startup
- Update checking and installation
- Offline license handling
- Error handling and user feedback

## Dependencies

- `github.com/keygen-sh/keygen-go/v3`: Official Keygen Go SDK
- Platform-specific dependencies for update installation

## Testing

Comprehensive tests will cover:
- License validation flows
- Update detection and installation
- Offline scenarios
- Error conditions
- Platform-specific behaviors

## Documentation

Complete documentation will be added to the Wails v3 docs including:
- Quick start guide
- API reference
- Best practices
- Troubleshooting guide
- Example applications