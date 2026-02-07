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
- Add `-tags` flag to `wails3 build` command for passing custom Go build tags (e.g., `wails3 build -tags gtk4`) (#4957)

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix file drag-and-drop on Windows not working at non-100% display scaling
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Windows
- Fix file drop coordinates being in wrong pixel space on Windows (physical vs CSS pixels)
- Fix file drag-and-drop on Linux not working reliably with hover effects
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Linux
- Fix window show/hide on Linux/GTK4 sometimes restoring to minimized state by using `gtk_window_present()` (#4957)
- Fix window position get/set on Linux/GTK4 always returning 0,0 by adding X11-conditional support via `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Fix max window size not being enforced on Linux/GTK4 by adding signal-based size clamping to replace removed `gtk_window_set_geometry_hints` (#4957)
- Fix DPI scaling on Linux/GTK4 by implementing proper PhysicalBounds calculation and fractional scaling support via `gdk_monitor_get_scale` (GTK 4.14+)
- Fix menu items duplicating when creating new windows on Linux/GTK4
- Fix generation of mapped types with enum keys in JS/TS bindings (#4437) by @fbbdev

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->
- Restrict GITHUB_TOKEN permissions in workflow files to follow principle of least privilege
- Fix path traversal vulnerability in screen example asset middleware
- Fix command injection vulnerability in setup wizard dependency installation endpoint
- Update rollup to 3.29.5 to fix XSS vulnerability (CVE-2024-47068)

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
