# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->
- Add Web API examples in `v3/examples/web-apis/` demonstrating 41 browser APIs including Storage (localStorage, sessionStorage, IndexedDB, Cache API), Network (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Device (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performance (Performance API, Mutation Observer, Intersection/Resize Observer), UI (Web Components, Pointer Events, Selection, Dialog, Drag and Drop), and more
- Add WebView API compatibility checker example (`v3/examples/webview-api-check/`) that tests 200+ browser APIs across platforms
- Add `internal/libpath` package for finding native library paths on Linux with parallel search, caching, and support for Flatpak/Snap/Nix
- **WIP:** Add experimental WebKitGTK 6.0 / GTK4 support for Linux, available via `-tags gtk4` (GTK3/WebKit2GTK 4.1 remains the default)
  - Note: On tiling window managers (e.g., Hyprland, Sway), Minimize/Maximize operations may not work as expected since the WM controls window geometry

## Changed
<!-- Changes in existing functionality -->
- **BREAKING**: Map keys in generated JS/TS bindings are now marked optional to accurately reflect Go map semantics. Map value access in Typescript now returns `T | undefined` instead of `T`, requiring null checks or assertions (#4943) by `@fbbdev`

## Fixed
<!-- Bug fixes -->
- Fix file drag-and-drop on Windows not working at non-100% display scaling
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Windows
- Fix file drop coordinates being in wrong pixel space on Windows (physical vs CSS pixels)
- Fix file drag-and-drop on Linux not working reliably with hover effects
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Linux
- Fix DPI scaling on Linux/GTK4 by implementing proper PhysicalBounds calculation and fractional scaling support via `gdk_monitor_get_scale` (GTK 4.14+)
- Fix menu items duplicating when creating new windows on Linux/GTK4
- Fix generation of mapped types with enum keys in JS/TS bindings (#4437) by @fbbdev

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new `SetWindowIcon()` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
