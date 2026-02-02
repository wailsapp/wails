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
- Add `UseApplicationMenu` option to `WebviewWindowOptions` allowing windows on Windows/Linux to inherit the application menu set via `app.Menu.Set()` by @leaanthony
- Added how to do `One Time Handlers` in the docs for `Listening to Events in JavaScript` by @AbdelhadiSeddar 

## Changed
<!-- Changes in existing functionality -->
- Move `EnabledFeatures`, `DisabledFeatures`, and `AdditionalBrowserArgs` from per-window options to application-level `Options.Windows` (#4559) by @leaanthony
- Changed the use of `Event` into `Events` according to changes in `@wailsio/runtime` and appropriate function calls in the docs in `Listening to Events in Javascript` by @AbdelhadiSeddar

## Fixed
<!-- Bug fixes -->
- Fix OpenFileDialog crash on Linux due to GTK thread safety violation (#3683) by @ddmoney420
- Fix SIGSEGV crash when calling `Focus()` on a hidden or destroyed window (#4890) by @ddmoney420
- Fix potential panic when setting empty icon or bitmap on Linux (#4923) by @ddmoney420
- Fix ErrorDialog crash when called from service binding on macOS (#3631) by @leaanthony
- Make menus to be displayed on Windows OS in `v3\examples\dialogs` by @ndianabasi
- Fix race condition causing TypeError during page reload (#4872) by @ddmoney420
- Fix incorrect output from binding generator tests by removing global state in the `Collector.IsVoidAlias()` method (#4941) by @fbbdev
- Fix `<input type="file">` file picker not working on macOS (#4862) by @leaanthony

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- **BREAKING**: Remove `EnabledFeatures`, `DisabledFeatures`, and `AdditionalLaunchArgs` from per-window `WindowsWindow` options. Use application-level `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures`, and `Options.Windows.AdditionalBrowserArgs` instead. These flags apply globally to the shared WebView2 environment (#4559) by @leaanthony

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
