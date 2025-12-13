# Wails v3 Distribution Features Implementation Plan

This document outlines the implementation plan for adding Electrobun-inspired distribution features to Wails v3. The goal is to provide a seamless, production-ready experience for code signing, cross-platform builds, and auto-updates.

## Executive Summary

### Current Wails v3 Status vs Electrobun

| Feature | Wails v3 Current | Electrobun | Gap |
|---------|-----------------|------------|-----|
| **Code Signing - macOS** | Ad-hoc only (development) | Full signing + notarization | Missing certificate identity, notarization automation |
| **Code Signing - Windows** | MSIX only | Full signing | EXE/NSIS signing missing |
| **Code Signing - Linux** | None | None | N/A |
| **Cross-Platform Builds** | Good (Taskfile-based) | CLI flag-based (`--targets`) | UX could be simpler |
| **Auto-Updates (App)** | Not implemented | Built-in Updater API with bsdiff | Complete feature gap |
| **Binary Delta Updates** | None | 14KB patches via bsdiff | Complete feature gap |
| **Build Compression** | Standard | ZSTD compression | Nice to have |

---

## Phase 1: Enhanced Code Signing

### 1.1 macOS Code Signing & Notarization

**Priority:** High
**Estimated Complexity:** Medium
**Dependencies:** None

#### Current State
- File: `v3/internal/commands/build_assets/darwin/Taskfile.yml`
- Only ad-hoc signing: `codesign --force --deep --sign - {{.BIN_DIR}}/{{.APP_NAME}}.app`
- Documentation exists for manual process: `docs/src/content/docs/guides/signing.mdx`

#### Implementation Tasks

- [ ] **1.1.1 Add signing configuration to build config**
  - Add to `config.yml` schema:
    ```yaml
    signing:
      macos:
        identity: "Developer ID Application: Company Name (TEAMID)"
        entitlements: "build/entitlements.plist"
        hardened_runtime: true
      notarization:
        enabled: true
        apple_id: "${APPLE_ID}"  # Environment variable reference
        team_id: "${APPLE_TEAM_ID}"
        keychain_profile: "notarytool-profile"  # Or use app-specific password
    ```

- [ ] **1.1.2 Implement signing identity detection**
  - File: Create `v3/internal/signing/darwin.go`
  - Auto-detect available signing identities via `security find-identity -v -p codesigning`
  - Provide helpful error messages when no valid identity found
  - Support identity selection via CLI flag: `--sign-identity="..."`

- [ ] **1.1.3 Implement proper code signing**
  - Replace ad-hoc signing with proper identity-based signing
  - Add hardened runtime support (`--options runtime`)
  - Add entitlements file support
  - Sign all binaries in the bundle (not just the app)
  - Sign frameworks and dylibs individually before signing app

- [ ] **1.1.4 Implement notarization automation**
  - File: Create `v3/internal/notarization/darwin.go`
  - Create zip archive for notarization
  - Submit via `xcrun notarytool submit`
  - Wait for completion with progress feedback
  - Staple ticket via `xcrun stapler staple`
  - Handle errors gracefully with clear messages

- [ ] **1.1.5 Add CLI commands**
  ```bash
  wails3 sign                    # Sign the built application
  wails3 notarize               # Notarize and staple
  wails3 package --sign --notarize  # All-in-one
  ```

- [ ] **1.1.6 Update Taskfile templates**
  - File: `v3/internal/commands/build_assets/darwin/Taskfile.yml`
  - Add signing task with configurable identity
  - Add notarization task
  - Make signing optional but default for production builds

#### Reference Implementation
Inspired by Electrobun's approach where code signing "is generally very fast" and notarization "takes about 1-2 minutes" with automatic stapling.

---

### 1.2 Windows Code Signing Enhancement

**Priority:** High
**Estimated Complexity:** Medium
**Dependencies:** None

#### Current State
- MSIX signing exists: `v3/internal/commands/msix.go`
- No signing for standalone EXE or NSIS installers
- MSIX command not exposed in CLI (registration issue)

#### Implementation Tasks

- [ ] **1.2.1 Fix MSIX CLI registration**
  - File: `v3/cmd/wails3/main.go`
  - Register `ToolMSIX` command properly
  - Ensure `wails3 tool msix` is accessible

- [ ] **1.2.2 Add EXE signing support**
  - File: Create `v3/internal/signing/windows.go`
  - Support PFX certificate files
  - Support Windows Certificate Store (thumbprint)
  - Support Azure Key Vault signing (for CI/CD)
  - Timestamp signing (required for production)

- [ ] **1.2.3 Add signing configuration**
  ```yaml
  signing:
    windows:
      certificate_path: "${WINDOWS_CERT_PATH}"
      certificate_password: "${WINDOWS_CERT_PASSWORD}"
      # Or use certificate store
      certificate_thumbprint: "${WINDOWS_CERT_THUMBPRINT}"
      timestamp_server: "http://timestamp.digicert.com"
  ```

- [ ] **1.2.4 Integrate with NSIS installer**
  - File: `v3/internal/commands/build_assets/windows/nsis/project.nsi.tmpl`
  - Sign the EXE before packaging
  - Sign the installer after creation
  - Add signing verification step

- [ ] **1.2.5 Add CLI commands**
  ```bash
  wails3 sign --platform windows
  wails3 package --sign
  ```

---

## Phase 2: Application Auto-Updates

### 2.1 Core Updater Runtime API

**Priority:** High
**Estimated Complexity:** High
**Dependencies:** None

#### Current State
- Only CLI update exists: `v3/internal/commands/update_cli.go`
- No runtime API for applications
- Documentation shows manual implementation: `docs/src/content/docs/guides/auto-updates.mdx`

#### Implementation Tasks

- [ ] **2.1.1 Define Updater API interface**
  - File: Create `v3/pkg/application/updater.go`
  ```go
  type Updater interface {
      // CheckForUpdate checks if a new version is available
      CheckForUpdate() (*UpdateInfo, error)

      // DownloadUpdate downloads the update (full or patch)
      DownloadUpdate(info *UpdateInfo, progress func(downloaded, total int64)) error

      // ApplyUpdate applies the downloaded update
      ApplyUpdate() error

      // GetCurrentVersion returns the current application version
      GetCurrentVersion() string

      // SetUpdateURL sets the base URL for update checks
      SetUpdateURL(url string)
  }

  type UpdateInfo struct {
      Version      string    `json:"version"`
      ReleaseDate  time.Time `json:"release_date"`
      ReleaseNotes string    `json:"release_notes"`
      DownloadURL  string    `json:"download_url"`
      PatchURL     string    `json:"patch_url,omitempty"`  // For delta updates
      Checksum     string    `json:"checksum"`
      PatchFrom    string    `json:"patch_from,omitempty"` // Version this patch applies to
      Size         int64     `json:"size"`
      PatchSize    int64     `json:"patch_size,omitempty"`
      Mandatory    bool      `json:"mandatory"`
  }
  ```

- [ ] **2.1.2 Implement update manifest format**
  - File: Create `v3/internal/updater/manifest.go`
  - JSON-based manifest (similar to Electrobun's `update.json`)
  ```json
  {
    "version": "1.2.0",
    "release_date": "2025-01-15T00:00:00Z",
    "release_notes": "Bug fixes and improvements",
    "platforms": {
      "darwin-arm64": {
        "url": "https://updates.example.com/app-1.2.0-darwin-arm64.tar.gz",
        "checksum": "sha256:...",
        "size": 12582912,
        "patches": [
          {
            "from": "1.1.0",
            "url": "https://updates.example.com/patches/1.1.0-to-1.2.0-darwin-arm64.patch",
            "checksum": "sha256:...",
            "size": 14336
          }
        ]
      }
    }
  }
  ```

- [ ] **2.1.3 Implement update checker**
  - File: Create `v3/internal/updater/checker.go`
  - HTTP client for fetching update manifest
  - Version comparison logic
  - Platform detection
  - Caching to avoid excessive checks
  - Rate limiting support

- [ ] **2.1.4 Implement update downloader**
  - File: Create `v3/internal/updater/downloader.go`
  - Resume support for interrupted downloads
  - Progress callbacks
  - Checksum verification
  - Temporary file management
  - Prefer patch files when available

- [ ] **2.1.5 Implement platform-specific update application**
  - Files:
    - `v3/internal/updater/apply_darwin.go`
    - `v3/internal/updater/apply_windows.go`
    - `v3/internal/updater/apply_linux.go`
  - macOS: Replace app bundle, handle code signing
  - Windows: Use update installer or in-place replacement
  - Linux: AppImage replacement or package manager integration

- [ ] **2.1.6 JavaScript/TypeScript bindings**
  - File: Create `v3/pkg/application/updater_bindings.go`
  - Expose Updater API to frontend
  ```typescript
  interface Updater {
      checkForUpdate(): Promise<UpdateInfo | null>;
      downloadUpdate(onProgress?: (downloaded: number, total: number) => void): Promise<void>;
      applyUpdate(): Promise<void>;
      getCurrentVersion(): string;
  }
  ```

- [ ] **2.1.7 Add configuration options**
  ```yaml
  updater:
    enabled: true
    url: "https://updates.example.com/myapp/"
    check_interval: "24h"  # Auto-check interval (0 to disable)
    allow_prerelease: false
    signature_verification: true
    public_key: "-----BEGIN PUBLIC KEY-----..."
  ```

---

### 2.2 Binary Delta Updates (bsdiff)

**Priority:** Medium
**Estimated Complexity:** High
**Dependencies:** Phase 2.1

#### Current State
- Not implemented
- Electrobun achieves ~14KB patches using SIMD-optimized bsdiff in Zig

#### Implementation Tasks

- [ ] **2.2.1 Integrate bsdiff library**
  - Options:
    1. Pure Go implementation: `github.com/gabstv/go-bsdiff`
    2. CGO wrapper for native bsdiff
    3. Ship pre-built bspatch binary (like Electrobun)
  - Recommendation: Use pure Go for simplicity, benchmark performance

- [ ] **2.2.2 Implement patch generation during build**
  - File: Create `v3/internal/patcher/generate.go`
  - Compare previous version binary with new version
  - Generate `.patch` file using bsdiff
  - Store patch metadata in manifest

- [ ] **2.2.3 Add build command for patch generation**
  ```bash
  wails3 build --generate-patches --previous-version=1.1.0
  ```
  - Requires access to previous version binaries
  - Generate patches for all configured versions (e.g., last 5 versions)

- [ ] **2.2.4 Implement patch application**
  - File: Create `v3/internal/patcher/apply.go`
  - Verify patch checksum before applying
  - Apply bspatch to create new binary
  - Verify result checksum
  - Fallback to full download on failure

- [ ] **2.2.5 Add patch validation**
  - Verify patch applies cleanly
  - Verify resulting binary checksum matches expected
  - Signature verification of patched binary

- [ ] **2.2.6 CLI tools for patch management**
  ```bash
  wails3 tool patch create --old=app-1.0.0 --new=app-1.1.0 --output=patch.bsdiff
  wails3 tool patch apply --binary=app-1.0.0 --patch=patch.bsdiff --output=app-1.1.0
  wails3 tool patch verify --patch=patch.bsdiff --checksum=sha256:...
  ```

---

### 2.3 Update Artifact Management

**Priority:** Medium
**Estimated Complexity:** Medium
**Dependencies:** Phase 2.1, 2.2

#### Implementation Tasks

- [ ] **2.3.1 Generate update artifacts during build**
  - `update.json` manifest file
  - Compressed application archive (`.tar.gz`, `.zip`)
  - Patch files for previous versions
  - Checksums file

- [ ] **2.3.2 Implement artifact structure**
  ```
  artifacts/
  ├── update.json                           # Update manifest
  ├── myapp-1.2.0-darwin-arm64.tar.gz      # Full update
  ├── myapp-1.2.0-darwin-amd64.tar.gz
  ├── myapp-1.2.0-windows-amd64.zip
  ├── myapp-1.2.0-linux-amd64.tar.gz
  ├── patches/
  │   ├── 1.1.0-to-1.2.0-darwin-arm64.patch
  │   ├── 1.0.0-to-1.2.0-darwin-arm64.patch
  │   └── ...
  └── checksums.txt
  ```

- [ ] **2.3.3 Add compression options**
  - Support ZSTD compression (like Electrobun)
  - Fallback to gzip for compatibility
  - Configuration option for compression level

- [ ] **2.3.4 Add release command**
  ```bash
  wails3 release --version=1.2.0 --output=./artifacts
  ```
  - Builds all configured platforms
  - Generates update manifest
  - Creates patches from previous versions
  - Signs all artifacts

---

### 2.4 Update Signature Verification

**Priority:** High
**Estimated Complexity:** Medium
**Dependencies:** Phase 2.1

#### Implementation Tasks

- [ ] **2.4.1 Implement Ed25519 signature generation**
  - File: Create `v3/internal/signing/updates.go`
  - Generate key pair: `wails3 tool keys generate`
  - Sign update artifacts during build
  - Store signatures in manifest

- [ ] **2.4.2 Implement signature verification**
  - Verify signature before downloading
  - Verify signature after download
  - Reject unsigned updates (configurable)

- [ ] **2.4.3 Key management**
  ```bash
  wails3 tool keys generate --output=keys/
  wails3 tool keys sign --key=private.key --file=update.json
  wails3 tool keys verify --key=public.key --file=update.json
  ```

---

## Phase 3: Cross-Platform Build Improvements

### 3.1 Simplified Cross-Compilation

**Priority:** Medium
**Estimated Complexity:** Low
**Dependencies:** None

#### Current State
- Cross-platform builds work via Taskfile
- Requires manual configuration
- No simple `--targets` flag like Electrobun

#### Implementation Tasks

- [ ] **3.1.1 Add `--targets` flag to build command**
  ```bash
  wails3 build --targets=darwin-arm64,darwin-amd64,windows-amd64,linux-amd64
  wails3 build --targets=all
  ```

- [ ] **3.1.2 Implement target resolution**
  - File: Create `v3/internal/build/targets.go`
  - Parse target strings (os-arch format)
  - Validate CGO requirements
  - Detect available cross-compilers

- [ ] **3.1.3 Add target presets**
  ```yaml
  build:
    presets:
      desktop:
        - darwin-arm64
        - darwin-amd64
        - windows-amd64
        - linux-amd64
      macos:
        - darwin-arm64
        - darwin-amd64
  ```

- [ ] **3.1.4 Improve cross-compilation documentation**
  - CGO cross-compilation setup guides
  - Docker-based build environments
  - GitHub Actions examples for all platforms

---

### 3.2 Universal Binary Support (macOS)

**Priority:** Medium
**Estimated Complexity:** Low
**Dependencies:** None

#### Current State
- Universal binary creation exists via `lipo` in Taskfile
- Works but could be more integrated

#### Implementation Tasks

- [ ] **3.2.1 Add universal binary flag**
  ```bash
  wails3 build --targets=darwin-universal
  ```

- [ ] **3.2.2 Automate lipo process**
  - Build arm64 and amd64 separately
  - Combine with lipo
  - Single output artifact

---

### 3.3 Build Artifact Organization

**Priority:** Low
**Estimated Complexity:** Low
**Dependencies:** None

#### Implementation Tasks

- [ ] **3.3.1 Implement artifact naming convention**
  - Format: `{app}-{version}-{os}-{arch}.{ext}`
  - Example: `myapp-1.2.0-darwin-arm64.dmg`

- [ ] **3.3.2 Organize build output**
  ```
  dist/
  ├── darwin-arm64/
  │   ├── MyApp.app/
  │   └── MyApp-1.2.0-darwin-arm64.dmg
  ├── darwin-amd64/
  ├── windows-amd64/
  │   ├── MyApp.exe
  │   ├── MyApp-1.2.0-windows-amd64-installer.exe
  │   └── MyApp-1.2.0-windows-amd64.msix
  └── linux-amd64/
      ├── myapp
      ├── myapp-1.2.0-linux-amd64.AppImage
      ├── myapp-1.2.0-linux-amd64.deb
      └── myapp-1.2.0-linux-amd64.rpm
  ```

---

## Phase 4: Developer Experience Improvements

### 4.1 Update Server Reference Implementation

**Priority:** Low
**Estimated Complexity:** Medium
**Dependencies:** Phase 2

#### Implementation Tasks

- [ ] **4.1.1 Create reference update server**
  - Simple HTTP server for update manifests
  - Support for S3/GCS/Azure Blob storage
  - CDN-friendly design

- [ ] **4.1.2 Add update hosting documentation**
  - Static file hosting setup (S3, Cloudflare R2, etc.)
  - CDN configuration
  - CORS requirements

---

### 4.2 CI/CD Templates

**Priority:** Medium
**Estimated Complexity:** Low
**Dependencies:** Phase 1, 2, 3

#### Implementation Tasks

- [ ] **4.2.1 GitHub Actions workflow templates**
  - Multi-platform build matrix
  - Code signing secrets management
  - Automatic release creation
  - Update artifact upload

- [ ] **4.2.2 GitLab CI templates**
  - Similar to GitHub Actions

- [ ] **4.2.3 Signing in CI documentation**
  - macOS: Keychain import from base64 certificate
  - Windows: PFX from secrets, Azure Key Vault
  - Certificate renewal process

---

## Implementation Order

### Recommended Sequence

1. **Phase 1.1** - macOS Code Signing (most requested)
2. **Phase 2.1** - Core Updater API (high value)
3. **Phase 1.2** - Windows Code Signing
4. **Phase 2.4** - Update Signature Verification (security critical)
5. **Phase 2.2** - Binary Delta Updates (bandwidth savings)
6. **Phase 3.1** - Simplified Cross-Compilation
7. **Phase 2.3** - Update Artifact Management
8. **Phase 3.2** - Universal Binary Support
9. **Phase 4.1** - Update Server Reference
10. **Phase 4.2** - CI/CD Templates

---

## File Structure

New files to create:

```
v3/
├── internal/
│   ├── signing/
│   │   ├── darwin.go       # macOS code signing
│   │   ├── windows.go      # Windows code signing
│   │   └── updates.go      # Update signature generation
│   ├── notarization/
│   │   └── darwin.go       # macOS notarization
│   ├── updater/
│   │   ├── manifest.go     # Update manifest handling
│   │   ├── checker.go      # Update checking
│   │   ├── downloader.go   # Update downloading
│   │   ├── apply_darwin.go # macOS update application
│   │   ├── apply_windows.go
│   │   └── apply_linux.go
│   ├── patcher/
│   │   ├── generate.go     # Patch generation
│   │   └── apply.go        # Patch application
│   └── build/
│       └── targets.go      # Cross-platform target handling
├── pkg/
│   └── application/
│       ├── updater.go      # Updater API
│       └── updater_bindings.go  # JS/TS bindings
└── cmd/
    └── wails3/
        └── main.go         # Register new commands
```

---

## Testing Strategy

### Unit Tests
- Signing identity detection
- Version comparison
- Manifest parsing
- Patch generation/application

### Integration Tests
- End-to-end signing workflow
- Update download and install
- Cross-platform build verification

### Manual Testing
- Real certificate signing
- App Store notarization
- Update flow on all platforms

---

## References

- [Electrobun Documentation](https://blackboard.sh/electrobun/docs/)
- [Electrobun GitHub](https://github.com/blackboardsh/electrobun)
- [bsdiff Algorithm](https://en.wikipedia.org/wiki/Bsdiff)
- [Apple Notarization](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Windows Code Signing](https://docs.microsoft.com/en-us/windows/win32/seccrypto/signtool)
