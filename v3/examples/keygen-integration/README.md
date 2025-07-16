# Keygen Integration Example

This example demonstrates how to integrate Keygen licensing and update services into a Wails v3 application.

## Features

This demo application showcases:

- **License Management**
  - License key validation
  - Online and offline license validation
  - License information display
  - License caching for offline use

- **Machine Activation**
  - Machine fingerprinting and activation
  - Machine deactivation
  - Machine-specific license binding

- **Feature Entitlements**
  - Check if specific features are enabled for the license
  - Feature flag management

- **Software Updates**
  - Check for available updates
  - Multiple update channels (stable, beta, alpha, dev)
  - Download updates with progress tracking
  - Install updates

- **Event System**
  - Real-time event notifications
  - License status events
  - Update availability events
  - Download progress events

## Prerequisites

1. A Keygen account (sign up at https://keygen.sh)
2. A product created in your Keygen account
3. At least one license key for testing

## Configuration

Before running the example, update the following in `main.go`:

```go
keygen.New(keygen.ServiceOptions{
    AccountID:      "your-account-id",      // Your Keygen account ID
    ProductID:      "your-product-id",      // Your product ID
    LicenseKey:     "",                     // Will be set via UI
    PublicKey:      "your-public-key",      // Your Ed25519 public key (optional)
    CurrentVersion: "1.0.0",                // Current app version
    AutoCheck:      true,                   // Enable auto-update checks
    UpdateChannel:  "stable",               // Default update channel
})
```

## Running the Example

1. Make sure you have Wails v3 installed
2. Navigate to this directory
3. Run the application:

```bash
go run .
```

## Usage

### License Validation

1. Enter your license key in the input field
2. Click "Validate License"
3. View the license status and details
4. Save the license for offline use if needed

### Machine Activation

After validating a license:
1. Click "Activate This Machine" to bind the license to this device
2. The machine fingerprint will be displayed
3. You can deactivate the machine when needed

### Feature Entitlements

1. Enter a feature name (e.g., "premium", "api-access")
2. Click "Check Feature" to see if it's enabled for your license

### Software Updates

1. Select an update channel (stable, beta, alpha, dev)
2. Click "Check for Updates"
3. If an update is available:
   - Review the release notes
   - Click "Download Update" to start downloading
   - Monitor the download progress
   - Click "Install Update" when download completes

### Event Log

The event log at the bottom shows all Keygen-related events in real-time, including:
- License validation results
- Machine activation/deactivation
- Update availability notifications
- Download progress updates
- Errors and debugging information

## API Reference

The example uses the following Keygen service methods:

- `SetLicenseKey(key)` - Set the license key for validation
- `ValidateLicense()` - Validate the current license
- `ActivateMachine()` - Activate the current machine
- `DeactivateMachine()` - Deactivate the current machine
- `GetLicenseInfo()` - Get detailed license information
- `CheckEntitlement(feature)` - Check if a feature is enabled
- `CheckForUpdates()` - Check for available updates
- `DownloadUpdate(releaseID)` - Download an update
- `InstallUpdate()` - Install a downloaded update
- `SetUpdateChannel(channel)` - Set the update channel
- `SaveOfflineLicense()` - Save license for offline use
- `LoadOfflineLicense()` - Load offline license
- `ClearLicenseCache()` - Clear cached license data

## Events

The application listens for these events:

- `keygen:license-status` - License status changes
- `keygen:update-available` - New update available
- `keygen:download-progress` - Update download progress
- `keygen:update-installed` - Update installation complete

## Demo Mode

This example includes a demo mode configuration. To use real licenses:

1. Replace the demo account credentials with your actual Keygen account details
2. Create licenses in your Keygen dashboard
3. Optionally create releases for testing the update functionality

## Security Notes

- Never hardcode license keys in your application
- Store the Ed25519 public key securely in your app for signature verification
- Use machine activation to prevent license sharing
- Implement proper error handling for network failures

## Further Resources

- [Keygen Documentation](https://keygen.sh/docs)
- [Wails v3 Documentation](https://v3.wails.app)
- [Keygen Go SDK](https://github.com/keygen-sh/keygen-go)