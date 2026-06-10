# wails/webview2

Generated Go bindings for the Microsoft Edge WebView2 SDK and the tooling
that produces them. This module is consumed by Wails on Windows.

## Versioning

The module path is `github.com/wailsapp/wails/webview2/v2` (tags
`webview2/v2.X.Y`). Releases follow a **v2.X.0** strategy: every routine
release — SDK regenerations from the update workflow, new tooling — bumps
the minor version, leaving the patch position free for targeted fixes to an
already-released minor. The release workflow auto-bumps the minor; pass an
explicit version to cut a patch.

`v1.x` tags (un-suffixed `github.com/wailsapp/wails/webview2`) predate the
generator overhaul and remain resolvable for existing consumers; all new
work lands on `/v2` only.

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

# Build pkg/webview2/capabilities.go from the SDK release notes and the
# IDL interface inventory (every generated interface gets an entry).
go run ./cmd/webview2gen capabilities          # fetches notes from GitHub
go run ./cmd/webview2gen capabilities --source test.md   # use a local copy

# Record an SDK bump in webview2/CHANGELOG.md (added/removed interfaces,
# new methods), diffing two cached IDLs.
go run ./cmd/webview2gen changelog --from 1.0.2903.40 --to 1.0.3967.48

# Regenerate everything (bindings AND the generated *_gen_test.go suite)
# and fail if the working tree differs from the committed output — wired
# into CI so hand-edits cannot land.
go run ./cmd/webview2gen verify

# Run download → generate → capabilities → verify in sequence.
go run ./cmd/webview2gen full --source test.md
```

`generate` is incremental: only files whose content changed are rewritten,
and generated files the IDL no longer produces are removed.

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
import "github.com/wailsapp/wails/webview2/v2/pkg/webview2"

runtime := webview2runtime.GetAvailableCoreWebView2BrowserVersionString() // e.g. "121.0.2277.83"

ok, required, err := webview2.SupportsInterface(runtime, "ICoreWebView2_22")
if err != nil { ... }
if !ok {
    log.Printf("feature disabled — needs WebView2 Runtime %s, installed %s", required, runtime)
    return
}
// safe to call ICoreWebView2_22 methods
```

The comparison is runtime-version against runtime-version: every entry in
`InterfaceSupportTable` carries the SDK release that introduced the
interface *and* that release's minimum WebView2 Runtime version, scraped
from the release notes. Interfaces that predate the release-notes archive
are baseline entries (always supported). The table covers every generated
interface — completeness is enforced at generation time and by a generated
test, so `SupportsInterface` never errors for an interface this package
defines.

## Bugs fixed by the v2 generator

| # | Severity | Pattern | Fix |
|---|----------|---------|-----|
| 1 | Critical | `[out] LPWSTR* name` produced `uintptr(unsafe.Pointer(_name))` — COM wrote the string pointer into the caller's stack | Now produces `uintptr(unsafe.Pointer(&_name))` so the local `*uint16` receives the string. Affects ~109 methods. |
| 2 | Moderate | `[in] LPWSTR*` was treated as a single `string` | Maps to `[]string`; the new `inputStringArraySetup.tmpl` marshals each element to `*uint16`. |
| 3 | Moderate | Generated interfaces lacked `Release()`, leaking refs | `Release() uint32` is emitted alongside `AddRef()` on every interface. |
| 4 | Minor | `QueryInterface` swallowed `HRESULT`, returning `nil` on failure | Returns `(*T, error)`; non-zero HRESULTs are surfaced. |
| 5 | Minor | `VARIANT` was a `uintptr` (8 bytes) instead of the 16-byte Windows struct | `type VARIANT struct { VT uint16; … Val [8]byte }`. |
| 6 | Minor | `AddRef/Release` returned `uintptr` instead of the COM `ULONG` (uint32) | Both return `uint32`; callback wrappers cast to `uintptr` at the syscall boundary. |
| 7 | Critical | Derived interface vtbls embedded only `IUnknownVtbl`, dropping every inherited slot — all 89 derived interfaces dispatched methods to the wrong vtbl entry | Vtbls embed their base interface's vtbl recursively, reproducing the flat COM layout. |
| 8 | Critical | `[in] double` emitted `uintptr(v)`, truncating `1.5` to `1` (12 methods incl. `PutZoomFactor`) | IEEE-754 bits via `math.Float64bits`; two stack words on 386. |
| 9 | Critical | 8-byte by-value aggregates (`EventRegistrationToken` in all 65 `Remove*` methods, `POINT`) passed by address — event handlers never unregistered | Passed by value per the Win64 ABI, with arm64/386 encodings. |
| 10 | Critical | All `Get<Interface>()` QI helpers were emitted on `*ICoreWebView2` — controller/settings/environment upgrades were unobtainable | Helpers target the root of each interface's inheritance chain. |
| 11 | Moderate | `COREWEBVIEW2_COLOR`/`RECT` by-value structs passed by address (wrong on every arch / on arm64+386 respectively); `VARIANT` declared 16 bytes instead of the 24-byte win64 layout | Size-driven by-value packing; `VARIANT.Val [2]uintptr`. |
| 12 | Moderate | `[out] LPWSTR**` surfaced as `*string`, `[in/out] T**`/`T***` arrays mismarshalled | Proper `[]string` / `[]*T` array marshalling with CoTaskMem ownership. |

## Testing

```sh
cd webview2/scripts
go test ./...                   # generator + internal pkgs (any OS)
cd ..
GOOS=windows go vet ./pkg/webview2/   # cross-compile bindings AND generated tests
go test ./pkg/webview2/               # on Windows: run the generated suite
```

The generator emits a behavioural test suite next to the bindings: every
generated method gets a round-trip test against a fake recording vtbl (see
`pkg/webview2/comtest_support_test.go` for the harness), plus generated
vtbl-layout, struct-size, QueryInterface and Invoke-dispatch tests. The
suite runs on the Windows CI job; other platforms cross-compile it via
`go vet`.

CI gates on `webview2gen verify` (bindings, tests and capabilities all
match a fresh regeneration) and the test suites above.

Windows VM integration tests against a real WebView2 runtime (COM smoke,
event handlers end-to-end) are tracked separately on WAI-297.

## See also

- `ARCHITECTURE.md` for the design rationale (programmatic emitter,
  file-per-interface, capability sources, etc.).
- Microsoft SDK release notes: <https://learn.microsoft.com/en-us/microsoft-edge/webview2/release-notes>
- WebView2 IDL on NuGet: <https://www.nuget.org/packages/Microsoft.Web.WebView2>
