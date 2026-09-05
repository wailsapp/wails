---
title: Obfuscated Builds
description: Build your Wails application with Garble to protect your source code from reverse engineering
sidebar:
  order: 8
---

import { Aside, Steps, Tabs, TabItem } from '@astrojs/starlight/components';

[Garble](https://github.com/burrowers/garble) is a Go build tool that replaces `go build` to rename symbols, obfuscate constants, and strip debug information from the resulting binary. Wails v3 has first-class support for Garble through two new commands.

## Prerequisites

- **Go 1.26.2 or later** — required by Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

<Aside type="tip">
Garble's minimum Go version changes between releases. If you are on an older Go toolchain, check the [Garble releases page](https://github.com/burrowers/garble/releases) for a version that matches your toolchain before installing.
</Aside>

## Required: add JSON tags to your service types

Any struct that a bound service method returns or accepts must have explicit JSON tags on every exported field:

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

<Aside type="caution">
Garble renames exported struct fields, and Wails passes these structs to `json.Marshal` through an `interface{}` parameter that Garble cannot statically trace. An obfuscated build without JSON tags will compile successfully but your frontend will receive garbled or empty field names at runtime. Add tags before running an obfuscated build.
</Aside>

Wails's own types — `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` — are already tagged. You only need to tag your own types.

## Building with obfuscation

<Steps>

1. **Generate the stable ID file**

   Run this any time you add, rename, or remove a bound service method:

   ```bash
   wails3 generate bindings -obfuscated
   ```

   This creates `wails_obfuscated.gen.go` in your main package directory — commit this file.

2. **Build with Garble**

   ```bash
   wails3 build --obfuscated
   ```

   Builds the application using the obfuscated bindings.

</Steps>

## Passing extra flags to Garble

Use `--garbleargs` to forward options directly to `garble`:

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

See the [Garble documentation](https://github.com/burrowers/garble#flags) for the full list of supported flags.

## Advanced: writing the ID file to a different package

By default `wails_obfuscated.gen.go` is written next to your `main` package. If your project keeps services in a sub-package imported by `main`, you can write the file there instead with `-obfuscated-output`:

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

<Aside type="caution">
The destination package must be imported — directly or transitively — by your `main` package so its `init()` runs at startup. If the package is not reachable, the stable IDs are never registered and binding calls will fail (for example, with `binding not found` errors at runtime).
</Aside>

## Troubleshooting

### `garble: command not found`

Garble is not installed or `$(go env GOPATH)/bin` is not on your `PATH`.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Frontend receives wrong or empty field values

Your service return types are missing `json:"..."` tags. Check every struct returned by your bound methods and add explicit tags to each exported field.

### `binding not found` error in the browser console

The stable ID file is missing or was not compiled in. Check:

- `wails_obfuscated.gen.go` exists in your main package directory (or the directory you passed to `-obfuscated-output`)
- You ran `wails3 build --obfuscated`, which adds the `wails_obfuscated` build tag
- If you used `-obfuscated-output`, the destination package is imported by `main`

### Windows Defender flags the build as a virus

Garble-obfuscated Go binaries are heuristically flagged by Windows Defender during the build because they lack debug symbols and resemble packed executables. The build fails with:

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Add your temp directory (where Go writes intermediate build artifacts) and your project directory to Defender's exclusion list:

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

These exclusions apply only to the specified paths and do not disable Defender globally.

### Build fails with `unsupported Go version`

Garble v0.16.0 requires Go 1.26.2 or later. Either upgrade Go, or consult the [Garble releases page](https://github.com/burrowers/garble/releases) for a version compatible with your toolchain.
