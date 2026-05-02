# Wails v3 Beta

Wails v3 is in public Beta. The core APIs are stable and ready for building real applications.

## Supported Platforms

| Platform | Status |
|---|---|
| macOS 12+ | Supported |
| Windows 10+ | Supported |
| Linux (GTK3) | Supported (default) |
| Linux (GTK4) | Experimental — build with `-tags gtk4` |

**Not in this Beta:** iOS and Android support is planned for a post-Beta release. WebAssembly is not in scope for v3.

## API Stability

The packages listed in [STABILITY.md](STABILITY.md) have stable APIs. Breaking changes will be documented in the changelog. Everything under `internal/` remains unstable.

## Getting Started

Full documentation is available at **https://v3.wails.io**.

## Migrating from v2

See the [v2 → v3 migration guide](https://v3.wails.io/guides/migrating/) for a step-by-step walkthrough of the breaking changes and how to update your application.

## Known Limitations

- **Wayland — window positioning:** The Wayland protocol does not expose an API for absolute window placement. `SetPosition` and related calls are a no-op on compositors that do not support `xdg-shell` position hints. Use X11/XWayland if you need deterministic positioning.
- **GTK4 — hardware testing:** Phase-10 hardware compatibility testing for GTK4 is still in progress. Expect rough edges on some GPU/driver combinations.
