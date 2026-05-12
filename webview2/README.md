# wails/webview2

Generated Go bindings for the Microsoft Edge WebView2 SDK and the tooling
that produces them. This module is consumed by Wails on Windows.

## Package layout

| Path | Purpose |
|------|---------|
| `pkg/webview2/` | **Generated** Go COM bindings for ICoreWebView2 (v1 through v27 at the latest SDK). Every `.go` file except `doc.go` is regenerated from `WebView2.idl`. |
| `pkg/edge/` | **Legacy** hand-maintained bindings for ICoreWebView2 v1–v4. Still used by `v3/pkg/application` via `chromium.go`; kept until the runtime is migrated. |
| `pkg/combridge/` | Generic COM bridge for implementing Go-side COM objects (event handlers, callbacks). |
| `webviewloader/` | Runtime DLL finder and version comparator (`CompareBrowserVersions`). |
| `internal/w32/` | Win32 syscall stubs used by the loader. |
| `scripts/` | The IDL cache (`WebView2.*.idl`), the `webview2gen` CLI, the parser, the emitter, and the internal helper packages. `module updater`. |

## webview2gen — the binding generator

`webview2gen` is the single entry point for refreshing `pkg/webview2`
against a new SDK release.

```sh
cd webview2/scripts

# Download a specific SDK version (or no flag for the latest cached version).
go run ./cmd/webview2gen download --version 1.0.2903.40

# Parse the cached IDL and emit pkg/webview2/*.go.
go run ./cmd/webview2gen generate

# Build pkg/webview2/capabilities.go from the SDK release notes.
go run ./cmd/webview2gen capabilities          # fetches notes from GitHub
go run ./cmd/webview2gen capabilities --source test.md   # use a local copy

# Regenerate everything and fail if the working tree differs from the
# committed output — wire this into CI so hand-edits cannot land.
go run ./cmd/webview2gen verify

# Run download → generate → capabilities → verify in sequence.
go run ./cmd/webview2gen full
```

Run `go generate ./...` from `webview2/` to invoke `full` via the
`//go:generate` directive in `pkg/webview2/doc.go`.

### Subcommand flags

Every subcommand accepts `--help`. Most useful defaults:

- `download --dir .` — where to store cached IDLs (relative to `scripts/`).
- `generate --version <v>` — pin to a specific cached IDL (default: newest).
- `generate --out ../pkg/webview2` — where to write generated files.
- `capabilities --source <path>` — skip the network fetch.
- `capabilities --json <path>` — also dump the interface→version map as JSON.

## Capabilities — runtime feature detection

The WebView2 SDK adds interfaces (ICoreWebView2_1 → ICoreWebView2_27 and
counting). A given runtime supports a subset; calling a method on an
interface the runtime doesn't know about crashes the process. To gate
features safely:

```go
import "github.com/wailsapp/wails/webview2/pkg/webview2"

runtime := webview2runtime.GetAvailableCoreWebView2BrowserVersionString() // e.g. "121.0.2277.83"

ok, required, err := webview2.SupportsInterface(runtime, "ICoreWebView2_22")
if err != nil { ... }
if !ok {
    log.Printf("feature disabled — needs WebView2 SDK %s, runtime reports %s", required, runtime)
    return
}
// safe to call ICoreWebView2_22 methods
```

For a higher-level gate, declare a `Capability` once and check it everywhere:

```go
var screenCapture = webview2.Capability{
    Name:        "screen-capture",
    Description: "ScreenCaptureStarting event support",
    Interface:   "ICoreWebView2_22",
}
if ok, _ := webview2.HasCapability(runtime, screenCapture); ok { ... }
```

`AllCapabilities` enumerates every interface as a named gate, suitable for
boot-time logging or diagnostics.

## Bugs fixed by the v2 generator

| # | Severity | Pattern | Fix |
|---|----------|---------|-----|
| 1 | Critical | `[out] LPWSTR* name` produced `uintptr(unsafe.Pointer(_name))` — COM wrote the string pointer into the caller's stack | Now produces `uintptr(unsafe.Pointer(&_name))` so the local `*uint16` receives the string. Affects ~109 methods. |
| 2 | Moderate | `[in] LPWSTR*` was treated as a single `string` | Maps to `[]string`; the new `inputStringArraySetup.tmpl` marshals each element to `*uint16`. |
| 3 | Moderate | Generated interfaces lacked `Release()`, leaking refs | `Release() uint32` is emitted alongside `AddRef()` on every interface. |
| 4 | Minor | `QueryInterface` swallowed `HRESULT`, returning `nil` on failure | Returns `(*T, error)`; non-zero HRESULTs are surfaced. |
| 5 | Minor | `VARIANT` was a `uintptr` (8 bytes) instead of the 16-byte Windows struct | `type VARIANT struct { VT uint16; … Val [8]byte }`. |
| 6 | Minor | `AddRef/Release` returned `uintptr` instead of the COM `ULONG` (uint32) | Both return `uint32`; callback wrappers cast to `uintptr` at the syscall boundary. |

## Testing

```sh
cd webview2/scripts
go test ./...                   # generator + internal pkgs (Linux-fine)
GOOS=windows go vet ../pkg/...  # cross-vet the generated bindings
```

CI must gate on `webview2gen verify` and `go test ./...`. The generator
contains regression tests for each of the 17 typed IDL patterns and for
all six known bug classes.

Windows VM integration tests (COM smoke, property getters, event handlers)
are tracked separately on WAI-297 — they require an actual WebView2 runtime.

## See also

- `ARCHITECTURE.md` for the design rationale (programmatic emitter,
  file-per-interface, capability sources, etc.).
- Microsoft SDK release notes: <https://learn.microsoft.com/en-us/microsoft-edge/webview2/release-notes>
- WebView2 IDL on NuGet: <https://www.nuget.org/packages/Microsoft.Web.WebView2>
