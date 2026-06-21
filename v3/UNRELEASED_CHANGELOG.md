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
- Add `application.System` for runtime platform detection from shared code: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (the `server` build tag), and `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` to test a single target directly. It compiles on every target so you can branch without build tags. Matching frontend helpers (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) are available in `@wailsio/runtime`
- Add a "Using Other Frontend Frameworks" guide showing how to drop your own Vite project into `frontend/` (covers Solid, Preact, Lit, SvelteKit, Qwik, Angular, etc.)
- The `wails3 setup` wizard now checks the mobile (iOS/Android) toolchain — Xcode and the iOS Simulator runtime, JDK, Android SDK/NDK and emulator — with one-click install and copyable shell-config fixes where applicable
- Generated projects ship a `frontend/.npmrc` setting a 7-day `minimum-release-age` to reduce exposure to freshly published (potentially compromised) packages (honoured by pnpm and bun; harmlessly ignored by npm)

## Changed
<!-- Changes in existing functionality -->
- Redesign all built-in starter templates with a new neon-mountain hero look (web, iOS and Android)
- **TypeScript is now the default for starter templates and owns the bare template name.** `wails3 init` (no `-t`) scaffolds a TypeScript project; `-t vanilla`, `-t react`, `-t vue` and `-t svelte` are TypeScript, with JavaScript variants at `-t vanilla-js`, `-t react-js`, `-t vue-js` and `-t svelte-js`. Built-in templates declare their language with `typescript:` in `template.yaml`; community templates using the `-ts` suffix continue to work as a fallback
- Redesign the `wails3 setup` wizard with the neon "digital Wails" theme (frosted-glass vibrancy over a mountain backdrop)

## Fixed
<!-- Bug fixes -->
- Fix crash on Windows when restoring an app that was minimised long enough for WebView2 to suspend or its render/GPU process to be recycled. The minimise/restore DPI resync (#5544) now only touches the WebView2 controller when the window's DPI actually changed, avoiding fatal COM calls into a suspended controller on the common same-DPI restore (#5605)
- Fix repeated native `SIGABRT`/`SIGSEGV` crashes (typically inside `g_object_unref` during the GTK main loop) on long-running Linux apps under frequent asset/media loads. The asset server completed `WebKitURISchemeRequest`s from worker goroutines, calling thread-unsafe WebKit2GTK functions off the GTK main thread; completion (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) now runs on the main thread. Completes the partial fix in #5566. Affects both the GTK3 and GTK4/WebKitGTK 6.0 builds (#5631, #5557)
- Fix intermittent `fatal error: invalid pointer found on stack` in `setupSignalHandlers` on Linux/GTK3. Window IDs passed as signal `user_data` were held in a Go `unsafe.Pointer` local, so the garbage collector aborted when it scanned the (non-pointer) value during a stack copy. The ID is now kept integer-typed (`uintptr_t`) on the Go side, back-porting to the legacy GTK3 path the same fix #4958 applied to the GTK4 path (which switched the C signal functions to `uintptr_t` to clear `-race`/checkptr errors) (#5631)

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove the `react-swc`, `preact`, `lit`, `solid`, `qwik` and `sveltekit` starter templates (and their `-ts` variants). The supported built-in set is now `vanilla`, `react`, `vue` and `svelte` — each TypeScript by default, with `-js` JavaScript variants. Any other framework can still be used by [bringing your own frontend](https://v3.wails.io/guides/dev/frontend-frameworks) or via a custom template

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
