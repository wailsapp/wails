# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 alpha release.
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

## Changed
Docs: Change to a couple of diagrams on architecture page to use sequence diagram for cleaner display
Docs: Include note about installing D2 as a prerequisite for running 

## Fixed
<!-- Bug fixes -->
- Fix `wails3 generate appimage` on the GTK4 default: the bundler now detects the GTK stack from the binary before searching for runtime files, so it picks `libwebkitgtkinjectedbundle.so` (under `webkitgtk-6.0/`) for GTK4 builds and `libwebkit2gtkinjectedbundle.so` (under `webkit2gtk-4.1/`) for `-tags gtk3` builds. The `.relr.dyn` probe also checks `libgtk-4.so.1` so stripping is correctly disabled on modern toolchains regardless of stack. (#5475)
- Fix `wails3 generate appimage` failing when invoked with a relative `-builddir`: the bundler now resolves `-binary`, `-icon`, `-desktopfile`, `-builddir` and `-outputdir` to absolute paths up-front so the mid-flow `s.CD` doesn't break the AppRun download goroutine or the post-copy `ldd` probe.
- Fix `wails3 generate appimage` failing to move the final AppImage to `-outputdir` when the desktop `Name=` field doesn't match the binary basename: the bundler now forces linuxdeploy's appimage plugin (via the `OUTPUT` env var) to write the AppImage to `<binary>-<arch>.AppImage` instead of the name derived from the desktop file.
- Fix `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` and `Common.SystemDidWake` not firing on Linux after the GTK4 + WebKitGTK 6.0 stack was promoted to the default in alpha.93. The new default `application_linux.go` `run()` wasn't calling `setupCommonEvents()` (which forwards `Linux.*` events to their `Common.*` counterparts) or `monitorPowerEvents()`. The DBus power-monitor helper is now shared between the GTK3 and GTK4 build paths via `application_linux_dbus.go`. (#5474)

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
