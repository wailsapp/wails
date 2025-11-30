# Start At Login Example

This example demonstrates how to use the Start At Login feature in Wails v3 applications.

## Features

- Configure the application to start automatically when the user logs in
- Toggle the start at login setting at runtime
- Check the current start at login status

## Platform-Specific Requirements

### macOS
- The application must be properly bundled (built with `wails3 build`)
- The system may prompt for accessibility permissions when first enabling start at login
- Consider adding `NSAppleEventsUsageDescription` to your Info.plist for better user experience:
  ```xml
  <key>NSAppleEventsUsageDescription</key>
  <string>This app needs to access System Events to manage login items.</string>
  ```

### Windows
- Uses the Windows Registry under `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run`
- No special permissions required

### Linux  
- Uses XDG autostart specification (creates `.desktop` files in `~/.config/autostart/`)
- Compatible with most Linux desktop environments (GNOME, KDE, XFCE, etc.)

## Usage

### Enable Start At Login during application initialization:

```go
app := application.New(application.Options{
    Name: "My App",
    StartAtLogin: true, // Enable start at login when app first runs
    // ... other options
})
```

### Toggle Start At Login at runtime:

```go
// Check current status
enabled, err := app.StartsAtLogin()
if err != nil {
    log.Printf("Error checking start at login: %v", err)
}

// Enable start at login
if err := app.SetStartAtLogin(true); err != nil {
    log.Printf("Error enabling start at login: %v", err)
}

// Disable start at login  
if err := app.SetStartAtLogin(false); err != nil {
    log.Printf("Error disabling start at login: %v", err)
}
```

## Building and Running

1. Build the application:
   ```bash
   wails3 build
   ```

2. Run the application:
   ```bash
   ./build/bin/start-at-login-demo
   ```

3. Use the interface to toggle the start at login setting

4. Log out and log back in to test that the application starts automatically (if enabled)

## Security Considerations

- The implementation validates executable paths to prevent injection attacks
- On macOS, AppleScript injection protection is implemented
- On Windows, restrictive registry permissions are used
- On Linux, proper file permissions are set for .desktop files

## Troubleshooting

### macOS
- If you get permission errors, check System Preferences > Security & Privacy > Privacy > Automation
- Ensure your app is properly code-signed for distribution
- For Mac App Store distribution, consider using `SMAppService` API (available in macOS 13+)

### Windows
- If registry access fails, ensure the user has write permissions to HKEY_CURRENT_USER
- Antivirus software may sometimes block registry modifications

### Linux
- Ensure `~/.config/autostart/` directory exists and is writable
- Check that your desktop environment supports XDG autostart specification
- Some desktop environments may require manual enabling of autostart functionality