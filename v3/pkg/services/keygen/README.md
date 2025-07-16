# Keygen Service

The Keygen service provides comprehensive licensing and automatic update functionality for Wails v3 applications, integrating with the [Keygen.sh](https://keygen.sh) licensing platform. This service handles license validation, machine activation, update checking, downloading, and installation with platform-specific implementations.

## Features

- **License Validation**: Validate licenses online and offline with cryptographic verification
- **Machine Activation**: Activate and deactivate machines with hardware fingerprinting
- **Automatic Updates**: Check for, download, and install updates automatically
- **Offline Support**: Cache license data for offline validation
- **Platform-Specific Installation**: Native update installation for macOS, Windows, and Linux
- **Event-Driven Architecture**: Real-time events for license status and update progress
- **Secure Storage**: Platform-specific secure storage for license keys (Keychain on macOS)
- **Entitlements**: Feature flag support through license entitlements
- **Update Channels**: Support for stable, beta, alpha, and dev release channels

## Installation

The Keygen service is included in Wails v3. To use it in your application:

1. Create a Keygen account at [keygen.sh](https://keygen.sh)
2. Create a product and obtain your account ID and product ID
3. Generate an Ed25519 public key for signature verification
4. Configure the service in your application

## Configuration

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/keygen"
)

app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(keygen.New(keygen.ServiceOptions{
            AccountID:      "your-account-id",
            ProductID:      "your-product-id",
            LicenseKey:     "user-license-key", // Optional - can be set later
            PublicKey:      "your-ed25519-public-key",
            CurrentVersion: "1.0.0",
            AutoCheck:      true,
            CheckInterval:  24 * time.Hour,
            UpdateChannel:  "stable",
            Environment:    "production", // or "staging"
        })),
    },
})
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `AccountID` | `string` | Required | Your Keygen account ID |
| `ProductID` | `string` | Required | Your Keygen product ID |
| `LicenseKey` | `string` | Optional | User's license key (can be set later) |
| `PublicKey` | `string` | Required | Ed25519 public key for signature verification |
| `CurrentVersion` | `string` | Required | Current application version |
| `AutoCheck` | `bool` | `false` | Enable automatic update checking |
| `CheckInterval` | `time.Duration` | `24h` | Interval between automatic update checks |
| `UpdateChannel` | `string` | `stable` | Release channel (stable, beta, alpha, dev) |
| `Environment` | `string` | `production` | Keygen environment |
| `CacheDir` | `string` | Platform default | Directory for caching license data |

## API Reference

### License Management

#### ValidateLicense

Validates the current license against the Keygen API.

```go
result, err := service.ValidateLicense()
if err != nil {
    // Handle error
}

if result.Valid {
    // License is valid
    if result.RequiresActivation {
        // Machine activation is required
    }
}
```

Returns:
- `*LicenseValidationResult`: Validation result with status, message, and state
- `error`: Error if validation fails

#### ActivateMachine

Activates the current machine for the validated license.

```go
result, err := service.ActivateMachine()
if err != nil {
    // Handle error
}

fmt.Printf("Machine activated: %s\n", result.Fingerprint)
```

Returns:
- `*MachineActivationResult`: Activation result with machine details
- `error`: Error if activation fails

#### DeactivateMachine

Deactivates the current machine.

```go
err := service.DeactivateMachine()
if err != nil {
    // Handle error
}
```

#### GetLicenseInfo

Retrieves detailed information about the current license.

```go
info, err := service.GetLicenseInfo()
if err != nil {
    // Handle error
}

fmt.Printf("License: %s\n", info.Key)
fmt.Printf("Email: %s\n", info.Email)
fmt.Printf("Expires: %v\n", info.ExpiresAt)
```

Returns:
- `*LicenseInfo`: Detailed license information
- `error`: Error if no license is available

#### CheckEntitlement

Checks if a specific feature is enabled in the license.

```go
hasFeature, err := service.CheckEntitlement("advanced-features")
if err != nil {
    // Handle error
}

if hasFeature {
    // Enable advanced features
}
```

### Update Management

#### CheckForUpdates

Checks for available updates.

```go
updateInfo, err := service.CheckForUpdates()
if err != nil {
    // Handle error
}

if updateInfo.Available {
    fmt.Printf("Update available: %s\n", updateInfo.LatestVersion)
    if updateInfo.Critical {
        // Handle critical update
    }
}
```

Returns:
- `*UpdateInfo`: Information about available updates
- `error`: Error if check fails

#### DownloadUpdate

Downloads an available update.

```go
progress, err := service.DownloadUpdate(updateInfo.ReleaseID)
if err != nil {
    // Handle error
}

// Progress is tracked via events
```

Returns:
- `*DownloadProgress`: Initial download progress info
- `error`: Error if download cannot start

#### InstallUpdate

Installs a downloaded update.

```go
err := service.InstallUpdate()
if err != nil {
    // Handle error
}
// Application will restart automatically
```

### Utility Methods

#### GetCurrentVersion

Returns the current application version.

```go
version := service.GetCurrentVersion()
```

#### SetUpdateChannel

Changes the update channel.

```go
err := service.SetUpdateChannel("beta")
if err != nil {
    // Handle invalid channel
}
```

#### SaveOfflineLicense / LoadOfflineLicense

Manages offline license caching.

```go
// Save current license for offline use
err := service.SaveOfflineLicense()

// Load offline license (called automatically on startup)
err := service.LoadOfflineLicense()
```

#### ClearLicenseCache

Clears all cached license data.

```go
err := service.ClearLicenseCache()
```

## Event Types

The service emits real-time events for license and update status changes. All events are emitted to the frontend and can be subscribed to using the Wails event system.

### UpdateAvailableEvent

Emitted when a new update is available for download.

**Event Name**: `keygen:update-available`

```go
type UpdateAvailableEvent struct {
    Version     string `json:"version"`      // New version number
    ReleaseDate string `json:"releaseDate"`  // ISO 8601 release date
    Notes       string `json:"notes"`        // Release notes
    Mandatory   bool   `json:"mandatory"`    // Is this a critical update
    DownloadURL string `json:"downloadUrl"`  // Download URL (set during download)
    Size        int64  `json:"size"`         // Download size in bytes
}
```

### LicenseStatusEvent

Emitted when the license status changes or is checked.

**Event Name**: `keygen:license-status`

```go
type LicenseStatusEvent struct {
    Valid       bool   `json:"valid"`        // Is license valid
    Status      string `json:"status"`       // Status: "active", "inactive", "expired", etc.
    Key         string `json:"key"`          // License key
    Email       string `json:"email"`        // License holder email
    ExpiresAt   string `json:"expiresAt"`   // ISO 8601 expiration date
    EnteredAt   string `json:"enteredAt"`   // ISO 8601 activation date
    LastChecked string `json:"lastChecked"` // ISO 8601 last check time
    Error       string `json:"error"`        // Error message if any
}
```

### DownloadProgressEvent

Emitted during update downloads to track progress.

**Event Name**: `keygen:download-progress`

```go
type DownloadProgressEvent struct {
    BytesDownloaded int64   `json:"bytesDownloaded"` // Bytes downloaded so far
    TotalBytes      int64   `json:"totalBytes"`      // Total size in bytes
    Progress        float64 `json:"progress"`        // Progress percentage (0-100)
    Speed           int64   `json:"speed"`           // Download speed in bytes/sec
    TimeRemaining   int     `json:"timeRemaining"`   // ETA in seconds
    Status          string  `json:"status"`          // Status: "downloading", "completed", "failed"
    Error           string  `json:"error"`           // Error message if failed
}
```

## Frontend Integration

Listen for events in your frontend code:

```javascript
// Listen for update available events
window.wails.Event.On('keygen:update-available', (event) => {
    console.log('Update available:', event.data);
    // Show update notification to user
    showUpdateNotification(event.data.version, event.data.notes);
});

// Listen for license status events
window.wails.Event.On('keygen:license-status', (event) => {
    console.log('License status:', event.data);
    // Update UI based on license status
    updateLicenseUI(event.data.valid, event.data.status);
});

// Listen for download progress events
window.wails.Event.On('keygen:download-progress', (event) => {
    console.log('Download progress:', event.data.progress + '%');
    // Update progress bar
    updateProgressBar(event.data.progress, event.data.speed);
});

// Call service methods from frontend
async function validateLicense() {
    try {
        const result = await window.KeygenService.ValidateLicense();
        if (result.Valid) {
            console.log('License is valid');
        }
    } catch (error) {
        console.error('License validation failed:', error);
    }
}

async function checkForUpdates() {
    try {
        const updateInfo = await window.KeygenService.CheckForUpdates();
        if (updateInfo.Available) {
            // Prompt user to download
            if (confirm(`Update ${updateInfo.LatestVersion} is available. Download now?`)) {
                await window.KeygenService.DownloadUpdate(updateInfo.ReleaseID);
            }
        }
    } catch (error) {
        console.error('Update check failed:', error);
    }
}
```

## Platform-Specific Features

### macOS Implementation

The macOS implementation provides comprehensive support for secure licensing and native update installation:

#### Features
- **Secure License Storage**: Uses macOS Keychain for encrypted license key storage
- **Machine Fingerprinting**: Combines hardware UUID, serial number, and MAC addresses for unique identification
- **App Bundle Updates**: Handles `.app` bundle replacement with proper permissions and code signing
- **App Translocation**: Removes quarantine attributes to prevent Gatekeeper translocation issues
- **Automatic Restart**: Gracefully restarts the application after update installation
- **Cache Directory**: `~/Library/Caches/[appname]/keygen/` for offline license data

#### Supported Update Formats
- `.app` bundles - Direct replacement of application bundle
- `.dmg` files - Automatically mounted, extracted, and installed
- `.zip` archives - Extracted and installed with permission preservation

#### Security Requirements
- Application must be code signed for production use
- Notarization recommended for smooth updates
- Keychain access requires user authorization on first use

### Windows Implementation

The Windows implementation provides Windows-native licensing and update support:

#### Features
- **Secure License Storage**: Uses Windows Credential Manager for encrypted storage
- **Machine Fingerprinting**: Combines machine GUID, CPU ID, and network adapter info
- **Registry Integration**: Stores application metadata in Windows Registry
- **Update Installation**: Handles `.exe` and `.msi` installers
- **Cache Directory**: `%LOCALAPPDATA%\[appname]\keygen\` for offline data

#### Supported Update Formats
- `.exe` installers - Self-extracting executables
- `.msi` packages - Windows Installer packages
- `.zip` archives - Extracted and installed

### Linux Implementation

The Linux implementation provides cross-distribution support:

#### Features
- **Secure License Storage**: Uses Secret Service API (libsecret) when available
- **Machine Fingerprinting**: Combines machine ID, product UUID, and network interfaces
- **Update Installation**: Supports AppImage, deb, rpm, and tar.gz formats
- **Cache Directory**: `~/.cache/[appname]/keygen/` for offline data

#### Supported Update Formats
- `.AppImage` - Portable application format
- `.deb` packages - Debian/Ubuntu packages
- `.rpm` packages - RedHat/Fedora packages
- `.tar.gz` archives - Generic Linux archives

## Security Considerations

### License Security
1. **Cryptographic Verification**: All licenses are verified using Ed25519 signatures
2. **Secure Storage**: Platform-specific secure storage for license keys
3. **Offline Validation**: Cached licenses include cryptographic proof for offline verification
4. **Machine Binding**: Licenses can be bound to specific machines using hardware fingerprints

### Update Security
1. **Signature Verification**: All updates are verified using Ed25519 signatures
2. **Checksum Validation**: SHA-256 checksums verified before installation
3. **Secure Download**: HTTPS-only downloads with certificate validation
4. **Atomic Updates**: Updates are downloaded to temporary locations before installation

### Best Practices
1. Always use HTTPS for Keygen API communication
2. Store the Ed25519 public key securely in your application
3. Implement proper error handling for network failures
4. Use machine activation for node-locked licenses
5. Regular license validation (recommended: daily)
6. Implement grace periods for offline usage

## Troubleshooting

### Common Issues

#### License Validation Fails
- **Check Network Connection**: Ensure the application can reach api.keygen.sh
- **Verify Public Key**: Ensure the Ed25519 public key matches your Keygen account
- **Check License Status**: Verify the license is active in your Keygen dashboard
- **Clock Synchronization**: Ensure system time is correct (affects signature validation)

#### Machine Activation Issues
- **Activation Limit**: Check if machine limit has been reached
- **Deactivate Old Machines**: Use the Keygen dashboard to manage activations
- **Fingerprint Changes**: Hardware changes may require reactivation

#### Update Installation Fails
- **Permissions**: Ensure the application has write permissions
- **Code Signing**: On macOS, ensure the update is properly signed
- **Disk Space**: Verify sufficient disk space for download and installation
- **Antivirus**: Some antivirus software may interfere with updates

### Debug Mode

Enable debug logging for troubleshooting:

```go
// Enable debug logging
service := keygen.New(keygen.ServiceOptions{
    // ... other options
    Debug: true,
})
```

### Error Codes

The service uses specific error codes for different failure scenarios:

| Code | Description |
|------|-------------|
| `ErrLicenseInvalid` | License validation failed |
| `ErrLicenseExpired` | License has expired |
| `ErrMachineLimit` | Machine activation limit reached |
| `ErrNetworkError` | Network communication failed |
| `ErrUpdateNotAvailable` | No update available |
| `ErrUpdateInProgress` | Update already in progress |
| `ErrUpdateInstallFailed` | Update installation failed |

## Complete Examples

### Basic License Validation

```go
package main

import (
    "log"
    "time"
    
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/keygen"
)

func main() {
    app := application.New(application.Options{
        Name: "My Licensed App",
        Services: []application.Service{
            application.NewService(keygen.New(keygen.ServiceOptions{
                AccountID:      "demo",
                ProductID:      "prod_123",
                PublicKey:      "e8601...key...",
                CurrentVersion: "1.0.0",
            })),
        },
    })
    
    // Run the application
    app.Run()
}
```

### Update System with UI

```javascript
// Frontend code for update management
class UpdateManager {
    constructor() {
        this.setupEventListeners();
    }
    
    setupEventListeners() {
        window.wails.Event.On('keygen:update-available', (event) => {
            this.showUpdateDialog(event.data);
        });
        
        window.wails.Event.On('keygen:download-progress', (event) => {
            this.updateProgress(event.data);
        });
    }
    
    async checkForUpdates() {
        try {
            const info = await window.KeygenService.CheckForUpdates();
            if (!info.Available) {
                this.showMessage('You are running the latest version');
            }
        } catch (error) {
            this.showError('Failed to check for updates: ' + error);
        }
    }
    
    showUpdateDialog(update) {
        const dialog = document.createElement('div');
        dialog.innerHTML = `
            <h2>Update Available</h2>
            <p>Version ${update.version} is available</p>
            <p>${update.notes}</p>
            <button onclick="updateManager.downloadUpdate('${update.releaseId}')">
                Download Update
            </button>
        `;
        document.body.appendChild(dialog);
    }
    
    async downloadUpdate(releaseId) {
        try {
            await window.KeygenService.DownloadUpdate(releaseId);
        } catch (error) {
            this.showError('Download failed: ' + error);
        }
    }
    
    updateProgress(progress) {
        const progressBar = document.getElementById('progress-bar');
        progressBar.style.width = progress.progress + '%';
        
        if (progress.status === 'completed') {
            this.promptInstall();
        }
    }
    
    async promptInstall() {
        if (confirm('Update downloaded. Install now?')) {
            await window.KeygenService.InstallUpdate();
            // App will restart automatically
        }
    }
}

const updateManager = new UpdateManager();
```

### License Activation Flow

```go
// Backend service for license activation
type LicenseService struct {
    keygen *keygen.Service
}

func (s *LicenseService) ActivateLicense(licenseKey string) error {
    // First validate the license
    result, err := s.keygen.ValidateLicense()
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if !result.Valid {
        return fmt.Errorf("invalid license: %s", result.Message)
    }
    
    // Check if machine activation is required
    if result.RequiresActivation {
        activation, err := s.keygen.ActivateMachine()
        if err != nil {
            return fmt.Errorf("activation failed: %w", err)
        }
        
        log.Printf("Machine activated: %s", activation.Fingerprint)
    }
    
    // Save for offline use
    if err := s.keygen.SaveOfflineLicense(); err != nil {
        log.Printf("Warning: failed to save offline license: %v", err)
    }
    
    return nil
}
```

## Migration Guide

### From Other Update Systems

If you're migrating from another update system (Sparkle, Squirrel, etc.), follow these steps:

1. **Set up Keygen Account**
   - Create account at keygen.sh
   - Create a product
   - Generate Ed25519 keys

2. **Configure Release Channels**
   - Set up stable, beta, and other channels as needed
   - Upload your existing releases

3. **Update Your Code**
   - Replace update check calls with Keygen service
   - Update event handlers for new event format
   - Implement license validation if needed

4. **Test Migration**
   - Test update flow with test licenses
   - Verify machine activation works correctly
   - Test offline scenarios

### Code Migration Example

```go
// Old (Sparkle/Squirrel)
updater.CheckForUpdates()

// New (Keygen)
keygenService.CheckForUpdates()

// Old event handling
updater.OnUpdateAvailable(func(version string) {
    // Handle update
})

// New event handling
window.wails.Event.On('keygen:update-available', (event) => {
    // Handle update with event.data
})
```

## Additional Resources

- [Keygen Documentation](https://keygen.sh/docs)
- [Wails v3 Documentation](https://wails.io/docs)
- [Example Applications](https://github.com/wailsapp/wails/tree/v3/examples)
- [Keygen Go SDK](https://github.com/keygen-sh/keygen-go)