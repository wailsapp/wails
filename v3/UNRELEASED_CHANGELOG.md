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
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix crash on Windows when restoring an app that was minimised long enough for WebView2 to suspend or its render/GPU process to be recycled. The minimise/restore DPI resync (#5544) now only touches the WebView2 controller when the window's DPI actually changed, avoiding fatal COM calls into a suspended controller on the common same-DPI restore (#5605)
- Fix repeated native `SIGABRT`/`SIGSEGV` crashes (typically inside `g_object_unref` during the GTK main loop) on long-running Linux apps under frequent asset/media loads. The asset server completed `WebKitURISchemeRequest`s from worker goroutines, calling thread-unsafe WebKit2GTK functions off the GTK main thread; completion (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) now runs on the main thread. Completes the partial fix in #5566. Affects both the GTK3 and GTK4/WebKitGTK 6.0 builds (#5631, #5557)
- Fix intermittent `fatal error: invalid pointer found on stack` in `setupSignalHandlers` on Linux/GTK3. Window IDs passed as signal `user_data` were held in a Go `unsafe.Pointer` local, so the garbage collector aborted when it scanned the (non-pointer) value during a stack copy. The ID is now kept integer-typed (`uintptr_t`) on the Go side, back-porting to the legacy GTK3 path the same fix #4958 applied to the GTK4 path (which switched the C signal functions to `uintptr_t` to clear `-race`/checkptr errors) (#5631)

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
