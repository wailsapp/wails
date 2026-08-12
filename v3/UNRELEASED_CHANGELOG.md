# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 release.
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
- Add Streams: bidirectional byte streams between Go and JavaScript with the WebSocket
  programming model and no listening socket. Declare a stream in Go with
  `app.HandleStream(name, handler)` and connect from the frontend with `Stream(name)`,
  which returns a `WebSocket`-shaped object. Go→JS is carried by one held poll per
  window over the asset server, JS→Go by a normal POST; nothing binds a TCP port and
  nothing goes through `evaluateJavaScript`. In server builds (`-tags server`) the same
  handler is served over a real WebSocket instead, so application code is identical
  across builds. by @leaanthony
- Move mailbox changelog entry to Unreleased in [PR](https://github.com/wailsapp/wails/pull/5935) by @leaanthony

## Changed
<!-- Changes in existing functionality -->
- Update docs sidebar autogeneration and blog author type derivation in [PR](https://github.com/wailsapp/wails/pull/5938) by @leaanthony

## Fixed
<!-- Bug fixes -->
- WebView2 initialization uses a deadline and message pump in [PR](https://github.com/wailsapp/wails/pull/5952) by @leaanthony
- WebView2 cookie test skips in CI unless opt-in and locks execution to current OS thread in [PR](https://github.com/wailsapp/wails/pull/5951) by @leaanthony
- Windows menu builders restore command IDs for submenu parent items in [PR](https://github.com/wailsapp/wails/pull/5944) by @gilad-ch
- Align the official cross-compilation image with the GTK 4.14+ Linux support baseline (#5928)
- Configure iOS Xcode project to retain inherited linker flags and add -ObjC in [PR](https://github.com/wailsapp/wails/pull/5915) by @mortenolsrud
- Fix excessive TCP connection churn in the `wails3 dev` asset proxy on large frontends, which could exhaust the host's ephemeral ports and make unrelated processes fail with `EADDRNOTAVAIL`
- Queue per-window event JavaScript for ordered dispatch and backpressure in [PR](https://github.com/wailsapp/wails/pull/5934) by @leaanthony

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove the desktop binary release pipeline: v3 releases are tag-only and the `wails3` CLI is installed with `go install`. Deletes `release-v3.yml` and the nightly step that dispatched it in [PR](https://github.com/wailsapp/wails/pull/5946) by @leaanthony

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
