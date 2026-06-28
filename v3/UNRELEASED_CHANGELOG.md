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
- Add global (system-wide) keyboard shortcuts via `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Shortcuts fire even when the application is not focused. Implemented natively per platform with no third party dependencies: Carbon hot keys on macOS, `RegisterHotKey` on Windows, `XGrabKey` on X11, and the XDG Desktop Portal global shortcuts interface on Wayland.
- Add built-in MCP server: a Model Context Protocol server that starts automatically when the application is built with the `mcp` tag, letting LLM agents test and control a running Wails application — window control, DOM inspection, JavaScript evaluation, bound method calls, events and simulated mouse/keyboard input rendered with an animated on-screen cursor. No user code required: the `mcp` tag is added automatically by `wails3 build`/`wails3 dev` when `WAILS_MCP=1` is set. Configured entirely via environment variables (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix `Menu.Update()` not rebuilding the native menu on GTK4 Linux (#5659, independently diagnosed and fixed by @puneetdixit200 in #5539)
- Fix crash enumerating macOS screens on display change by copying screen id/name strings and snapshotting the count (#5565, independently diagnosed and fixed by @x-haose in #5584)

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
