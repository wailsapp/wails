# Mac PKG Installer Builder

This package provides comprehensive Mac PKG installer creation and notarization tools for Wails v3 applications.

## Features

- **PKG Creation**: Creates professional macOS installer packages using `pkgbuild` and `productbuild`
- **Code Signing**: Supports signing with Developer ID certificates
- **Notarization**: Full Apple notarization workflow with automatic retry and status monitoring
- **Configuration-Driven**: YAML-based configuration with environment variable support
- **Template Generation**: Generates sample configuration files
- **Validation**: Comprehensive configuration and dependency validation

## Quick Start

### 1. Generate a Configuration Template

```bash
wails3 tool package --format pkg --generate-template --config my-pkg-config.yaml
```

### 2. Edit the Configuration

Edit the generated `my-pkg-config.yaml` file:

```yaml
# Application information
app_name: "MyApp"
app_path: "./build/MyApp.app"
bundle_id: "com.mycompany.myapp"
version: "1.0.0"

# Code signing (required for notarization)
signing_identity: "Developer ID Installer: Your Name (TEAM_ID)"

# Installation
install_location: "/Applications"

# Distribution package appearance
title: "MyApp Installer"

# Notarization (optional but recommended)
apple_id: "${APPLE_ID}"
app_password: "${APP_PASSWORD}"
team_id: "${TEAM_ID}"

# Output
output_path: "./dist/MyApp-installer.pkg"
```

### 3. Build the PKG Installer

```bash
wails3 tool package --format pkg --config my-pkg-config.yaml
```

## Configuration Options

### Required Fields

| Field | Description | Example |
|-------|-------------|---------|
| `app_name` | Application name | `"MyApp"` |
| `app_path` | Path to .app bundle | `"./build/MyApp.app"` |
| `bundle_id` | Unique bundle identifier | `"com.mycompany.myapp"` |
| `version` | Version number | `"1.0.0"` |
| `output_path` | Output PKG file path | `"./dist/MyApp.pkg"` |

### Optional Fields

| Field | Description | Default |
|-------|-------------|---------|
| `signing_identity` | Code signing identity | None |
| `install_location` | Installation directory | `/Applications` |
| `title` | Installer title | Same as `app_name` |
| `background` | Background image (PNG) | None |
| `welcome_file` | Welcome RTF file | None |
| `readme_file` | Readme RTF file | None |
| `license_file` | License RTF file | None |

### Notarization Fields

| Field | Description | Source |
|-------|-------------|--------|
| `apple_id` | Apple ID for notarization | Apple Developer Account |
| `app_password` | App-specific password | [App-Specific Passwords](https://support.apple.com/en-us/HT204397) |
| `team_id` | Developer team ID | Apple Developer Account |

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to configuration file | Required |
| `--generate-template` | Generate sample config | `false` |
| `--skip-notarization` | Skip notarization step | `false` |
| `--validate-only` | Only validate configuration | `false` |

## Environment Variables

Configuration values support environment variable expansion using `${VAR_NAME}` syntax:

```yaml
apple_id: "${APPLE_ID}"
app_password: "${APP_PASSWORD}"
team_id: "${TEAM_ID}"
```

Set environment variables:

```bash
export APPLE_ID="your.email@example.com"
export APP_PASSWORD="abcd-efgh-ijkl-mnop"
export TEAM_ID="ABC123DEF4"
```

## Notarization Setup

### 1. Get Required Credentials

1. **Apple ID**: Your Apple Developer account email
2. **App-Specific Password**: Generate at [appleid.apple.com](https://appleid.apple.com/account/manage)
3. **Team ID**: Found in Apple Developer account under "Membership"

### 2. Code Signing Certificate

Install a "Developer ID Installer" certificate from Apple Developer:

```bash
# List available signing identities
security find-identity -p basic -v

# Use the full identity name in your config
signing_identity: "Developer ID Installer: Your Name (TEAM_ID)"
```

## Examples

### Basic PKG (No Notarization)

```yaml
app_name: "SimpleApp"
app_path: "./build/SimpleApp.app"
bundle_id: "com.example.simpleapp"
version: "1.0.0"
output_path: "./SimpleApp-installer.pkg"
```

### Full-Featured PKG with Notarization

```yaml
app_name: "MyApp"
app_path: "./build/MyApp.app"
bundle_id: "com.mycompany.myapp"
version: "2.1.0"
signing_identity: "Developer ID Installer: My Company (ABC123DEF4)"
install_location: "/Applications"
title: "MyApp v2.1 Installer"
background: "./assets/installer-bg.png"
welcome_file: "./docs/welcome.rtf"
license_file: "./docs/license.rtf"
apple_id: "${APPLE_ID}"
app_password: "${APP_PASSWORD}"
team_id: "${TEAM_ID}"
output_path: "./dist/MyApp-v2.1-installer.pkg"
```

### Custom Installation Location

```yaml
app_name: "DevTool"
app_path: "./build/DevTool.app"
bundle_id: "com.dev.tool"
version: "1.0.0"
install_location: "/usr/local/bin"
output_path: "./DevTool-installer.pkg"
```

## Workflow Integration

### CI/CD Pipeline

```bash
#!/bin/bash
# Build and package for macOS

# Build the app
wails3 build -platform darwin

# Create PKG installer
wails3 tool package --format pkg --config build/pkg-config.yaml

# The signed and notarized PKG is ready for distribution
```

### Multiple Formats

```bash
# Create both DMG and PKG
wails3 tool package --format dmg --config dmg-config.yaml
wails3 tool package --format pkg --config pkg-config.yaml
```

## Troubleshooting

### Common Issues

1. **"pkgbuild not found"**
   - Install Xcode Command Line Tools: `xcode-select --install`

2. **"Notarization failed"**
   - Verify Apple ID credentials
   - Ensure app is properly signed
   - Check Apple Developer account status

3. **"Package signature validation failed"**
   - Verify signing identity is correct
   - Ensure certificate is valid and not expired

### Validation

```bash
# Validate configuration without building
wails3 tool package --format pkg --config my-config.yaml --validate-only

# Skip notarization for testing
wails3 tool package --format pkg --config my-config.yaml --skip-notarization

# Check package signature
spctl --assess --verbose --type install MyApp-installer.pkg
```

## Dependencies

- **macOS**: Required for PKG building
- **Xcode Command Line Tools**: Required for `pkgbuild`, `productbuild`, `notarytool`, `stapler`
- **Developer ID Certificate**: Required for signing and notarization
- **Apple Developer Account**: Required for notarization

## Security Considerations

- Store sensitive credentials (Apple ID, app password) in environment variables
- Use app-specific passwords instead of main Apple ID password
- Keep signing certificates secure and backed up
- Regularly rotate app-specific passwords

## Related Commands

- `wails3 build`: Build the application
- `wails3 tool package --format dmg`: Create DMG files
- `wails3 tool package --format deb`: Create Linux packages

## Support

For issues related to Mac PKG building:

1. Check system dependencies with built-in validation
2. Verify Apple Developer account status
3. Review notarization logs for detailed error information
4. Consult Apple's notarization documentation