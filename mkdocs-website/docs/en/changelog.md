# Changelog

<!--
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

-->

## [Unreleased]

### Added
- [darwin] add Event ApplicationShouldHandleReopen to be able to handle dock icon click by @5aaee9 in [#2991](https://github.com/wailsapp/wails/pull/2991)
- [darwin] add getPrimaryScreen/getScreens to impl by @tmclane in [#2618](https://github.com/wailsapp/wails/pull/2618)
- [darwin] add option for showing the toolbar in fullscreen mode on macOS by [@fbbdev](https://github.com/fbbdev) in [#3282](https://github.com/wailsapp/wails/pull/3282)
- [linux] add onKeyPress logic to convert linux keypress into an accelerator @[Atterpac](https://github.com/Atterpac) in[#3022](https://github.com/wailsapp/wails/pull/3022])
- [linux] add task `run:linux` by [@marcus-crane](https://github.com/marcus-crane) in [#3146](https://github.com/wailsapp/wails/pull/3146)
- Export `SetIcon` method by  @almas1992 in [PR](https://github.com/wailsapp/wails/pull/3147)
- Improve `OnShutdown` by  @almas1992 in [PR](https://github.com/wailsapp/wails/pull/3189)
- Restore `ToggleMaximise` method in `Window` interface by [@fbbdev](https://github.com/fbbdev) in [#3281](https://github.com/wailsapp/wails/pull/3281)
- Added more information to `Environment()`. By @leaanthony in [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c) based on [PR](https://github.com/wailsapp/wails/pull/2044) by @Mai-Lapyst
- Expose the `WebviewWindow.IsFocused` method on the `Window` interface by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- Support multiple space-separated trigger events in the WML system by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- Add ESM exports from the bundled JS runtime script by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)

### Fixed

- Fixed zombie processes when working in devmode by updating to latest refresh by [Atterpac](https://github.com/atterpac) in [#3320](https://github.com/wailsapp/wails/pull/3320).
- Fixed appimage webkit file sourcing by [Atterpac](https://github.com/atterpac) in [#3306](https://github.com/wailsapp/wails/pull/3306).
- Fixed Doctor apt package verify by [Atterpac](https://github.com/Atterpac) in [#2972](https://github.com/wailsapp/wails/pull/2972).
- Fixed application frozen when quit (Darwin) by @5aaee9 in [#2982](https://github.com/wailsapp/wails/pull/2982)
- Fixed background colours of examples on Windows by [mmgvh](https://github.com/mmghv) in [#2750](https://github.com/wailsapp/wails/pull/2750).
- Fixed default context menus by [mmgvh](https://github.com/mmghv) in [#2753](https://github.com/wailsapp/wails/pull/2753).
- Fixed hex values for arrow keys on Darwin by [jaybeecave](https://github.com/jaybeecave) in [#3052](https://github.com/wailsapp/wails/pull/3052).
- Set drag-n-drop for windows to working. Added by [@pylotlight](https://github.com/pylotlight) in [PR](https://github.com/wailsapp/wails/pull/3039)
- Fixed bug for linux in doctor in the event user doesn't have proper drivers installed. Added by [@pylotlight](https://github.com/pylotlight) in [PR](https://github.com/wailsapp/wails/pull/3032)
- Fix dpi scaling on start up (windows). Changed by @almas1992 in [PR](https://github.com/wailsapp/wails/pull/3145)
- Fix replace line in `go.mod` to use relative paths. Fixes Windows paths with spaces - @leaanthony.
- Fix MacOS systray click handling when no attached window by [thomas-senechal](https://github.com/thomas-senechal) in PR [#3207](https://github.com/wailsapp/wails/pull/3207)
- Fix failing Windows build due to unknown option by [thomas-senechal](https://github.com/thomas-senechal) in PR [#3208](https://github.com/wailsapp/wails/pull/3208)
- Fix wrong baseURL when open window twice by @5aaee9 in PR [#3273](https://github.com/wailsapp/wails/pull/3273)
- Fix ordering of if branches in `WebviewWindow.Restore` method by [@fbbdev](https://github.com/fbbdev) in [#3279](https://github.com/wailsapp/wails/pull/3279)
- Correctly compute `startURL` across multiple `GetStartURL` invocations when `FRONTEND_DEVSERVER_URL` is present. [#3299](https://github.com/wailsapp/wails/pull/3299)
- Fix the JS type of the `Screen` struct to match its Go counterpart by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- Fix the `WML.Reload` method to ensure proper cleanup of registered event listeners by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)

### Changed

<<<<<<< HEAD
=======
- The bundled JS runtime script is now an ESM module: script tags importing it must have the `type="module"` attribute. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- The `@wailsio/runtime` package does not publish its API on the `window.wails` object, and does not start the WML system. This has been done to improve encapsulation. The WML system can be started manually if desired by calling the new `WML.Enable` method. The bundled JS runtime script still performs both operations automatically. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- The Window API module `@wailsio/runtime/src/window` now exposes the containing window object as a default export. It is not possible anymore to import individual methods through ESM named or namespace import syntax.
- The JS window API has been updated to match the current Go `WebviewWindow` API. Some methods have changed name or prototype, specifically: `Screen` becomes `GetScreen`; `GetZoomLevel`/`SetZoomLevel` become `GetZoom`/`SetZoom`; `GetZoom`, `Width` and `Height` now return values directly instead of wrapping them within objects. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
>>>>>>> ee53801f (Update changelog)

### Removed

### Deprecated

### Security
