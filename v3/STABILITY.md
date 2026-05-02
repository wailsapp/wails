# Wails v3 Beta — API Stability Guarantee

This document defines what "Beta" means for the Wails v3 release and which parts of the API are covered by stability guarantees.

## What Beta Means

Wails v3 Beta signals that the **core API is ready for production evaluation**. Specifically:

- **Stable packages will not receive breaking API changes** without first going through a deprecation cycle (minimum one minor release with a deprecation notice).
- New functionality may still be added in a backwards-compatible way.
- Bug fixes may change observable behaviour when the previous behaviour was clearly wrong.
- Pre-release identifiers (`-beta.N`) will continue until we are confident no further breaking changes are needed, at which point we will release `v3.0.0`.

Beta is **not** a promise that every platform or every package is production-ready. See the sections below for exactly what is and is not covered.

## Stable Packages

The following packages are covered by the Beta stability guarantee:

| Package | Description |
|---|---|
| `v3/pkg/application` | Core application, window, and webview APIs |
| `v3/pkg/events` | Application and window event types |
| `v3/pkg/w32` | Windows-specific Win32 helpers |
| `v3/pkg/icons` | Cross-platform icon loading utilities |
| `v3/pkg/mac` | macOS-specific helpers and types |

Identifiers exported from these packages — types, functions, constants, and interface methods — will not be renamed, removed, or have their signatures changed in a breaking way during the Beta period without a prior deprecation notice.

Internal packages (`v3/internal/...`) and unexported identifiers are never covered by this guarantee.

## Experimental — Not Covered by Beta Guarantees

The following are still evolving and may change without a deprecation cycle:

### GTK4 (Linux)

GTK4 support is opt-in via the `-tags gtk4` build tag. It is **not** the default Linux build target (which uses GTK3/WebKitGTK). GTK4 phase-10 integration testing is still pending and the API surface may shift.

Use it at your own risk; do not rely on GTK4-specific behaviour in production.

### iOS and Android

Mobile platform support (`application_ios.go`, `application_android.go`) is under active development. Architecture is documented in `IOS_ARCHITECTURE.md` and `ANDROID_ARCHITECTURE.md`, but the platform-specific option types and mobile lifecycle APIs are not yet stable. iOS and Android targets are **deferred** from the Beta guarantee and will be promoted to stable in a future release.

### `v3/pkg/services`

The sub-packages under `v3/pkg/services` — `dock`, `fileserver`, `kvstore`, `log`, `notifications`, `sqlite` — are shipping as **preview services**. Their interfaces are reasonably stable today, but we reserve the right to refine them based on community feedback before locking them down. Treat them as "stabilising" and pin to a specific Beta version if you depend on them.

## Breaking Change Policy

### What Counts as Breaking

A breaking change is any change that requires callers to modify their source code or build configuration. Examples:

- Removing or renaming an exported type, function, constant, or method.
- Changing the signature of an exported function or interface method.
- Changing the semantics of a function in a way that silently alters correct programs.
- Changing required build constraints on stable packages.

The following are **not** considered breaking:

- Adding new exported identifiers.
- Adding new fields to a struct (callers should not rely on unkeyed struct literals, which can break if fields are added or reordered).
- Fixing behaviour that was clearly a bug.
- Changes to internal packages.

### Semver Signals

| Version bump | Meaning |
|---|---|
| `v3.0.0-beta.(N+1)` | May contain deprecations; no removals of previously stable API |
| `v3.0.0-rc.N` | API frozen; only bug fixes |
| `v3.0.0` | Full stable release |
| `v3.1.0` | Backwards-compatible additions |
| `v4.0.0` | May remove previously deprecated identifiers |

### Communication Process

1. A GitHub issue labelled **`breaking-change`** is opened at least one minor release before the change lands.
2. The deprecated identifier is annotated with a `// Deprecated:` comment in the source.
3. The change is called out explicitly in the release notes under a **Breaking Changes** heading.
4. If the change affects a common usage pattern, a migration guide or codemod will be provided.

## Known Beta Limitations

- **Wayland window positioning** on native Wayland sessions may be ignored (effectively a no-op, depending on compositor behaviour), so applications should not rely on setting arbitrary window positions. This is not a Wails bug and will not be "fixed."
- **GTK4 phase-10 testing** is pending. GTK4 behaviour on edge-case display configurations is not yet fully validated.
- **iOS and Android** are deferred. Mobile targets will be promoted to stable after their APIs stabilise and platform-specific testing is complete.
- **`v3/pkg/services`** interfaces are subject to revision based on Beta feedback before being locked down.

