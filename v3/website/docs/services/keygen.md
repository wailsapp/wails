# Keygen Service

The Keygen service integrates [Keygen.sh](https://keygen.sh) licensing and automatic updates into your Wails v3 application, providing a complete solution for software monetization and distribution.

## Quick Start

### 1. Set Up Your Keygen Account

1. Sign up at [keygen.sh](https://keygen.sh)
2. Create a product in your dashboard
3. Note your Account ID and Product ID
4. Generate an Ed25519 public key for license verification

### 2. Add the Service to Your App

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/keygen"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(keygen.New(keygen.ServiceOptions{
                AccountID:      "your-account-id",
                ProductID:      "your-product-id",
                PublicKey:      "your-ed25519-public-key",
                CurrentVersion: "1.0.0",
                AutoCheck:      true,  // Enable automatic update checks
            })),
        },
    })
    
    app.Run()
}
```

### 3. Handle License Validation

```javascript
// Frontend: Validate a license key
async function validateLicense(key) {
    try {
        const result = await window.KeygenService.ValidateLicense(key);
        if (result.Valid) {
            showSuccess("License validated successfully!");
            
            // Activate machine if required
            if (result.RequiresActivation) {
                const activation = await window.KeygenService.ActivateMachine();
                console.log("Machine activated:", activation.Fingerprint);
            }
        } else {
            showError("Invalid license: " + result.Message);
        }
    } catch (error) {
        showError("Validation failed: " + error);
    }
}
```

### 4. Implement Automatic Updates

```javascript
// Listen for update notifications
window.wails.Event.On('keygen:update-available', (event) => {
    const update = event.data;
    if (update.mandatory) {
        // Force update for critical releases
        showMandatoryUpdate(update);
    } else {
        // Show optional update notification
        showUpdateNotification(update);
    }
});

// Check for updates manually
async function checkForUpdates() {
    const info = await window.KeygenService.CheckForUpdates();
    if (info.Available) {
        showUpdateDialog(info);
    } else {
        showMessage("You're running the latest version!");
    }
}
```

## Key Features

### 🔐 Licensing
- **Online & Offline Validation**: Works with or without internet connection
- **Machine Activation**: Lock licenses to specific devices
- **Entitlements**: Feature flags for premium functionality
- **Secure Storage**: Platform-specific encrypted storage

### 🚀 Automatic Updates
- **Background Checks**: Periodic update checking
- **Progress Tracking**: Real-time download progress
- **Auto-Installation**: One-click update installation
- **Release Channels**: Stable, beta, and custom channels

### 🛡️ Security
- **Ed25519 Signatures**: Cryptographic license verification
- **Checksum Validation**: Integrity verification for updates
- **Secure Communication**: HTTPS-only API calls
- **Hardware Fingerprinting**: Unique machine identification

## Common Use Cases

### Trial License Implementation

```go
func (s *LicenseService) CheckTrialStatus() (bool, int, error) {
    info, err := s.keygen.GetLicenseInfo()
    if err != nil {
        return false, 0, err
    }
    
    // Check for trial entitlement
    isTrial, _ := s.keygen.CheckEntitlement("trial")
    if !isTrial {
        return false, 0, nil
    }
    
    // Calculate days remaining
    if info.ExpiresAt != nil {
        daysLeft := int(time.Until(*info.ExpiresAt).Hours() / 24)
        return true, daysLeft, nil
    }
    
    return true, 0, nil
}
```

### Feature Gating

```javascript
// Check if user has access to premium features
async function initializeFeatures() {
    const hasPro = await window.KeygenService.CheckEntitlement("pro-features");
    const hasTeam = await window.KeygenService.CheckEntitlement("team-features");
    
    if (hasPro) {
        enableProFeatures();
    }
    
    if (hasTeam) {
        enableTeamFeatures();
    }
}
```

### Update UI with Progress

```javascript
// Complete update UI implementation
class UpdateUI {
    constructor() {
        this.progressBar = document.getElementById('update-progress');
        this.statusText = document.getElementById('update-status');
        this.setupListeners();
    }
    
    setupListeners() {
        window.wails.Event.On('keygen:update-available', this.onUpdateAvailable.bind(this));
        window.wails.Event.On('keygen:download-progress', this.onProgress.bind(this));
    }
    
    onUpdateAvailable(event) {
        const update = event.data;
        this.showUpdateDialog({
            version: update.version,
            notes: update.notes,
            size: this.formatBytes(update.size),
            mandatory: update.mandatory
        });
    }
    
    onProgress(event) {
        const data = event.data;
        this.progressBar.style.width = data.progress + '%';
        this.statusText.textContent = `Downloading: ${data.progress.toFixed(1)}% at ${this.formatBytes(data.speed)}/s`;
        
        if (data.status === 'completed') {
            this.promptInstall();
        } else if (data.status === 'failed') {
            this.showError(data.error);
        }
    }
    
    formatBytes(bytes) {
        const sizes = ['B', 'KB', 'MB', 'GB'];
        if (bytes === 0) return '0 B';
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
    }
}
```

## Best Practices

### 1. License Validation Strategy

```go
// Implement a grace period for offline usage
func (s *Service) ValidateWithGracePeriod() error {
    // Try online validation first
    result, err := s.keygen.ValidateLicense()
    if err == nil && result.Valid {
        // Save for offline use
        s.keygen.SaveOfflineLicense()
        return nil
    }
    
    // Fall back to offline validation
    offlineErr := s.keygen.LoadOfflineLicense()
    if offlineErr != nil {
        return fmt.Errorf("license validation failed: %w", err)
    }
    
    // Check grace period (e.g., 7 days)
    info, _ := s.keygen.GetLicenseInfo()
    if info.LastValidatedAt != nil {
        daysSinceValidation := int(time.Since(*info.LastValidatedAt).Hours() / 24)
        if daysSinceValidation > 7 {
            return errors.New("grace period expired, online validation required")
        }
    }
    
    return nil
}
```

### 2. Update Channels

```javascript
// Let users choose their update channel
async function setUpdateChannel(channel) {
    const channels = ['stable', 'beta', 'alpha', 'dev'];
    
    if (!channels.includes(channel)) {
        showError('Invalid update channel');
        return;
    }
    
    try {
        await window.KeygenService.SetUpdateChannel(channel);
        localStorage.setItem('update-channel', channel);
        showSuccess(`Switched to ${channel} channel`);
        
        // Check for updates in new channel
        await window.KeygenService.CheckForUpdates();
    } catch (error) {
        showError('Failed to change channel: ' + error);
    }
}
```

### 3. Error Handling

```javascript
// Comprehensive error handling
async function handleLicenseAction(action) {
    try {
        await action();
    } catch (error) {
        // Parse Keygen error codes
        if (error.includes('MACHINE_LIMIT_REACHED')) {
            showError('Machine limit reached. Please deactivate another device.');
            showDeactivationDialog();
        } else if (error.includes('LICENSE_EXPIRED')) {
            showError('Your license has expired.');
            showRenewalOptions();
        } else if (error.includes('NETWORK_ERROR')) {
            showWarning('Cannot connect to license server. Working offline.');
        } else {
            showError('An unexpected error occurred: ' + error);
        }
    }
}
```

## Platform-Specific Notes

### macOS
- **Code Signing Required**: Your app must be signed for updates to work
- **Notarization**: Recommended for smooth installation
- **Keychain Access**: First-time license storage requires user permission

### Windows
- **UAC Prompts**: Update installation may trigger UAC
- **Antivirus**: Some AV software may flag updates - consider signing
- **Registry Access**: Service stores metadata in HKCU

### Linux
- **Package Formats**: Supports AppImage, deb, rpm, and tar.gz
- **Permissions**: May need elevated permissions for system-wide installs
- **Desktop Integration**: AppImage format recommended for simplicity

## Troubleshooting

### License Issues

**Problem**: "License validation failed"
- Check your internet connection
- Verify the public key matches your Keygen account
- Ensure the license is active in Keygen dashboard
- Check system time is correct

**Problem**: "Machine limit reached"
- Deactivate unused machines in Keygen dashboard
- Implement machine deactivation in your app
- Consider increasing machine limit for the license

### Update Issues

**Problem**: "Update installation failed"
- Ensure app has write permissions
- Check available disk space
- Verify update package is signed (macOS)
- Temporarily disable antivirus

**Problem**: "No updates found"
- Verify update channel setting
- Check if releases exist for your platform
- Ensure version comparison is working correctly

### Debug Mode

Enable detailed logging for troubleshooting:

```go
// In your service configuration
keygen.ServiceOptions{
    // ... other options
    Debug: true,  // Enable debug logging
}
```

## FAQ

**Q: Can I use Keygen service without licensing?**
A: Yes! You can use it just for automatic updates by not setting a license key.

**Q: How secure is the license validation?**
A: Very secure. Licenses are cryptographically signed and verified using Ed25519.

**Q: What happens if Keygen.sh is down?**
A: The service includes offline validation support with configurable grace periods.

**Q: Can I customize the update installation process?**
A: Yes, you can handle the events and implement custom UI/UX for updates.

**Q: Does it support gradual rollouts?**
A: Yes, through Keygen's release channels and distribution rules.

## Next Steps

1. [Create your Keygen account](https://keygen.sh)
2. [Read the full API documentation](/v3/pkg/services/keygen/README.md)
3. [View example implementations](https://github.com/wailsapp/wails/tree/v3/examples)
4. [Learn about Keygen policies](https://keygen.sh/docs/policies)

Need help? Join our [Discord community](https://discord.gg/wails) or check the [Keygen documentation](https://keygen.sh/docs).