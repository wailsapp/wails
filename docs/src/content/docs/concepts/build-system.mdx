---
title: Build System
description: Understanding how Wails builds and packages your application
sidebar:
  order: 4
---

import { Tabs, TabItem } from "@astrojs/starlight/components";

## Unified Build System

Wails provides a **unified build system** that compiles Go code, bundles frontend assets, embeds everything into a single executable, and handles platform-specific builds—all with one command.

```bash
wails3 build
```

**Output:** Native executable with everything embedded.

## Build Process Overview

{/*
TODO: Fix D2 diagram generation or embed as image.
The previous D2 code block was causing MDX parsing errors in the build pipeline.
*/}

**[Build Process Diagram Placeholder]**


## Build Phases

### 1. Analysis Phase

Wails scans your Go code to understand your services:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**What Wails extracts:**
- Service name: `GreetService`
- Method name: `Greet`
- Parameter types: `string`
- Return types: `string`

**Used for:** Generating TypeScript bindings

### 2. Generation Phase

#### TypeScript Bindings

Wails generates type-safe bindings:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Benefits:**
- Full type safety
- IDE autocomplete
- Compile-time errors
- JSDoc comments

#### Frontend Build

Your frontend bundler runs (Vite, webpack, etc.):

```bash
# Vite example
vite build --outDir dist
```

**What happens:**
- JavaScript/TypeScript compiled
- CSS processed and minified
- Assets optimised
- Source maps generated (dev only)
- Output to `frontend/dist/`

### 3. Compilation Phase

#### Go Compilation

Go code is compiled with optimisations:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Flags:**
- `-s`: Strip symbol table
- `-w`: Strip DWARF debugging info
- Result: Smaller binary (~30% reduction)

**Platform-specific:**
- Windows: `.exe` with icon embedded
- macOS: `.app` bundle structure
- Linux: ELF binary

#### Asset Embedding

Frontend assets are embedded into the Go binary:

```go
//go:embed frontend/dist
var assets embed.FS
```

**Result:** Single executable with everything inside.

### 4. Output

**Single native binary:**
- Windows: `myapp.exe` (~15MB)
- macOS: `myapp.app` (~15MB)
- Linux: `myapp` (~15MB)

**No dependencies** (except system WebView).

## Development vs Production

<Tabs syncKey="mode">
  <TabItem label="Development (wails3 dev)" icon="seti:config">
    **Optimised for speed:**
    
    ```bash
    wails3 dev
    ```
    
    **What happens:**
    1. Starts frontend dev server (Vite on port 9245 by default)
    2. Compiles Go without optimisations
    3. Launches app pointing to dev server
    4. Enables hot reload
    5. Includes source maps
    
    **Characteristics:**
    - **Fast rebuilds** (&lt;1s for frontend changes)
    - **No asset embedding** (served from dev server)
    - **Debug symbols** included
    - **Source maps** enabled
    - **Verbose logging**
    
    **File size:** Larger (~50MB with debug symbols)
  </TabItem>

  <TabItem label="Production (wails3 build)" icon="rocket">
    **Optimised for size and performance:**
    
    ```bash
    wails3 build
    ```
    
    **What happens:**
    1. Builds frontend for production (minified)
    2. Compiles Go with optimisations
    3. Strips debug symbols
    4. Embeds assets
    5. Creates single binary
    
    **Characteristics:**
    - **Optimised code** (minified, tree-shaken)
    - **Assets embedded** (no external files)
    - **Debug symbols stripped**
    - **No source maps**
    - **Minimal logging**
    
    **File size:** Smaller (~15MB)
  </TabItem>
</Tabs>

## Build Commands

### Basic Build

```bash
wails3 build
```

**Output:** `bin/<APP_NAME>` (or `bin/<APP_NAME>.exe` on Windows). The `bin/` directory sits at the project root.

`wails3 build` is a thin wrapper around `wails3 task build`. The only build-time flag it forwards is `--tags`, which becomes the `EXTRA_TAGS` Taskfile variable:

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

There are no `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags`, or `-package` flags on `wails3 build`. Cross-compilation, output paths, icons, and packaging are controlled through the project's Taskfile (`Taskfile.yml` + `build/config.yml`).

### Cross-platform and platform-specific builds

Platform builds are exposed as Taskfile tasks under the `darwin:` / `windows:` / `linux:` namespaces (defined in `build/Taskfile.<platform>.yml`). For example:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

To see every available task in the current project:

```bash
wails3 task --list
```

### Icons and packaging

Generate platform icons (`build/icons.icns`, `build/icon.ico`, etc.) from a source PNG:

```bash
wails3 generate icons -input appicon.png
```

Build platform-specific installers/packages:

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Build Configuration

### Taskfile.yml

Wails 3 projects use [Taskfile](https://taskfile.dev/) as their build orchestrator. The root `Taskfile.yml` includes per-platform task files from `build/`:

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

Run tasks with `wails3 task <name>` or `task <name>`:

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Project configuration: `build/config.yml`

Project metadata (name, identifier, version, info-plist values, NSIS settings, `.desktop` fields, custom protocols, etc.) lives in `build/config.yml`. The Taskfile reads this file when generating icons, manifests, installers, and the like. There is **no** `build/build.json` file in Wails 3.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Run `wails3 generate build-assets` (or `wails3 update build-assets`) to refresh the platform-specific build assets from this config.

## Asset Embedding

### How It Works

Wails uses Go's `embed` package:

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**At build time:**
1. Frontend built to `frontend/dist/`
2. `//go:embed` directive includes files
3. Files compiled into binary
4. Binary contains everything

**At runtime:**
1. App starts
2. Assets served from memory
3. No disk I/O for assets
4. Fast loading

### Custom Assets

Embed additional files:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Build Optimisations

### Frontend Optimisations

**Vite (default):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**Results:**
- JavaScript minified (~70% reduction)
- CSS minified (~60% reduction)
- Images optimised
- Tree-shaking applied

### Go Optimisations

**Compiler flags:**

```bash
-ldflags="-s -w"
```

- `-s`: Strip symbol table (~10% reduction)
- `-w`: Strip DWARF debug info (~20% reduction)

**Additional optimisations:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: Set variable values at build time
- Useful for version numbers, build dates

### Binary Compression

**UPX (optional):**

```bash
# After building
upx --best bin/myapp.exe
```

**Results:**
- ~50% size reduction
- Slightly slower startup (~100ms)
- Not recommended for macOS (code signing issues)

## Platform-Specific Builds

### Windows

**Output:** `myapp.exe`

**Includes:**
- Application icon
- Version information
- Manifest (UAC settings)

**Icon:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

The Windows `tool package` step then embeds the generated `.ico` into the executable.

**Manifest:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**Output:** `myapp.app` (application bundle)

**Structure:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**Universal Binary:**

The macOS Taskfile ships a `darwin:build:universal` (and `darwin:package:universal`) task that builds both architectures and combines them via `wails3 tool lipo`:

```bash
wails3 task darwin:build:universal
```

### Linux

**Output:** `myapp` (ELF binary)

**Dependencies:**
- GTK3
- WebKitGTK

**Desktop file:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Installation:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Build Performance

### Typical Build Times

| Phase | Time | Notes |
|-------|------|-------|
| Analysis | &lt;1s | Go code scanning |
| Binding Generation | &lt;1s | TypeScript generation |
| Frontend Build | 5-30s | Depends on project size |
| Go Compilation | 2-10s | Depends on code size |
| Asset Embedding | &lt;1s | Embedding frontend |
| **Total** | **10-45s** | First build |
| **Incremental** | **5-15s** | Subsequent builds |

### Speeding Up Builds

**1. Use build cache:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Run only what you need:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Parallel builds (multi-machine / CI):**

Cross-compilation between Linux/Windows/macOS in v3 generally happens in a Docker `wails-cross` container or on dedicated runners per platform — `wails3 build` itself targets the host OS. See [Cross-platform Builds](/guides/build/cross-platform) for the supported workflows.

**4. Use faster tools:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Troubleshooting

### Build Fails

**Symptom:** `wails3 build` exits with error

**Common causes:**

1. **Go compilation error**
   ```bash
   # Check Go code compiles
   go build
   ```

2. **Frontend build error**
   ```bash
   # Check frontend builds
   cd frontend
   npm run build
   ```

3. **Missing dependencies**
   ```bash
   # Install dependencies
   npm install
   go mod download
   ```

### Binary Too Large

**Symptom:** Binary is >50MB

**Solutions:**

1. **Strip debug symbols** (the shipped Taskfile already passes `-ldflags="-s -w"` to `go build`).

2. **Check embedded assets**
   ```bash
   # Remove unnecessary files from frontend/dist/
   # Check for large images, videos, etc.
   ```

3. **Use UPX compression**
   ```bash
   upx --best bin/myapp.exe
   ```

### Slow Builds

**Symptom:** Builds take >1 minute

**Solutions:**

1. **Use build cache**
   - Go cache is automatic
   - Frontend cache (Vite) is automatic

2. **Run only the task you need**
   ```bash
   wails3 task common:build:frontend
   wails3 task windows:build
   ```

3. **Optimise frontend build**
   ```javascript
   // vite.config.js
   export default {
     build: {
       minify: 'esbuild',  // Faster than terser
     },
   }
   ```

## Best Practices

### ✅ Do

- **Use `wails3 dev` during development** - Fast iteration
- **Use `wails3 build` for releases** - Optimised output
- **Version your builds** - Use `-ldflags` to embed version
- **Test builds on target platforms** - Cross-compilation isn't perfect
- **Keep frontend builds fast** - Optimise bundler config
- **Use build cache** - Speeds up subsequent builds

### ❌ Don't

- **Don't commit `build/` directory** - Add to `.gitignore`
- **Don't skip testing builds** - Always test before release
- **Don't embed unnecessary assets** - Keep binaries small
- **Don't use debug builds for production** - Use optimised builds
- **Don't forget code signing** - Required for distribution

## Next Steps

**Building Applications** - Detailed guide to building and packaging  
[Learn More →](/guides/building)

**Cross-Platform Builds** - Build for all platforms from one machine  
[Learn More →](/guides/cross-platform)

**Creating Installers** - Create installers for end users  
[Learn More →](/guides/installers)

---

**Questions about building?** Ask in [Discord](https://discord.gg/JDdSxwjhGf) or check the [build examples](https://github.com/wailsapp/wails/tree/master/v3/examples/build).
