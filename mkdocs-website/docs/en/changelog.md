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
- Add binding generator flag for using the bundled JS runtime script instead of the npm package by [@fbbdev](https://github.com/fbbdev) in [#3334](https://github.com/wailsapp/wails/pull/3334)
- Implement `setIcon` on linux by [@abichinger](https://github.com/abichinger) in [#3354](https://github.com/wailsapp/wails/pull/3354)
- Add flag `-port` to dev command and support environment variable `WAILS_VITE_PORT` by [@abichinger](https://github.com/abichinger) in [#3429](https://github.com/wailsapp/wails/pull/3429)
- Add tests for bound method calls by [@abichinger](https://github.com/abichinger) in [#3431](https://github.com/wailsapp/wails/pull/3431)

### Fixed

- [linux] Fixed theme handling error on NixOS by [tmclane](https://github.com/tmclane) in [#3515)(https://github.com/wailsapp/wails/pull/3515)
- Fixed cross volume project install for windows by [atterpac](https://github.com/atterac) in [#3512](https://github.com/wailsapp/wails/pull/3512)
- Fixed react template css to show footer by [atterpac](https://github.com/atterpac) in [#3477](https://github.com/wailsapp/wails/pull/3477)
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
- Fix crash on windows left clicking the systray icon when not having an attached window [tw1nk](https://github.com/tw1nk) in PR [#3271](https://github.com/wailsapp/wails/pull/3271)
- Fix wrong baseURL when open window twice by @5aaee9 in PR [#3273](https://github.com/wailsapp/wails/pull/3273)
- Fix ordering of if branches in `WebviewWindow.Restore` method by [@fbbdev](https://github.com/fbbdev) in [#3279](https://github.com/wailsapp/wails/pull/3279)
- Correctly compute `startURL` across multiple `GetStartURL` invocations when `FRONTEND_DEVSERVER_URL` is present. [#3299](https://github.com/wailsapp/wails/pull/3299)
- Fix the JS type of the `Screen` struct to match its Go counterpart by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- Fix the `WML.Reload` method to ensure proper cleanup of registered event listeners by [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- Fix custom context menu closing immediately on linux by [@abichinger](https://github.com/abichinger) in [#3330](https://github.com/wailsapp/wails/pull/3330)
- Fix the output path and extension of model files produced by the binding generator by [@fbbdev](https://github.com/fbbdev) in [#3334](https://github.com/wailsapp/wails/pull/3334)
- Fix the import paths of model files in JS code produced by the binding generator by [@fbbdev](https://github.com/fbbdev) in [#3334](https://github.com/wailsapp/wails/pull/3334)
- Fix drag-n-drop on some linux distros by [@abichinger](https://github.com/abichinger) in [#3346](https://github.com/wailsapp/wails/pull/3346)
- Fix missing task for macOS when using `wails3 task dev` by [@hfoxy](https://github.com/hfoxy) in [#3417](https://github.com/wailsapp/wails/pull/3417)
- Fix registering events causing a nil map assignment by [@hfoxy](https://github.com/hfoxy) in [#3426](https://github.com/wailsapp/wails/pull/3426)
- Fix unmarshaling of bound method parameters by [@fbbdev](https://github.com/fbbdev) in [#3431](https://github.com/wailsapp/wails/pull/3431)
- Fix handling of multiple return values from bound methods by [@fbbdev](https://github.com/fbbdev) in [#3431](https://github.com/wailsapp/wails/pull/3431)
- Fix doctor detection of npm that is not installed with system package manager by [@pekim](https://github.com/pekim) in [#3458](https://github.com/wailsapp/wails/pull/3458)

### Changed

- Update linux webkit dependency to webkit2gtk-4.1 over webkitgtk2-4.0 to support Ubuntu 24.04 LTS by [atterpac](https://github.com/atterpac) in [#3461](https://github.com/wailsapp/wails/pull/3461)
- The bundled JS runtime script is now an ESM module: script tags importing it must have the `type="module"` attribute. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- The `@wailsio/runtime` package does not publish its API on the `window.wails` object, and does not start the WML system. This has been done to improve encapsulation. The WML system can be started manually if desired by calling the new `WML.Enable` method. The bundled JS runtime script still performs both operations automatically. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- The Window API module `@wailsio/runtime/src/window` now exposes the containing window object as a default export. It is not possible anymore to import individual methods through ESM named or namespace import syntax.
- The JS window API has been updated to match the current Go `WebviewWindow` API. Some methods have changed name or prototype, specifically: `Screen` becomes `GetScreen`; `GetZoomLevel`/`SetZoomLevel` become `GetZoom`/`SetZoom`; `GetZoom`, `Width` and `Height` now return values directly instead of wrapping them within objects. By [@fbbdev](https://github.com/fbbdev) in [#3295](https://github.com/wailsapp/wails/pull/3295)
- The binding generator now uses calls by ID by default. The `-id` CLI option has been removed. Use the `-names` CLI option to switch back to calls by name. By [@fbbdev](https://github.com/fbbdev) in [#3468](https://github.com/wailsapp/wails/pull/3468)
- New binding code layout: output files were previously organised in folders named after their containing package; now full Go import paths are used, including the module path. By [@fbbdev](https://github.com/fbbdev) in [#3468](https://github.com/wailsapp/wails/pull/3468)
- The struct field `application.Options.Bind` has been renamed to `application.Options.Services`. By [@fbbdev](https://github.com/fbbdev) in [#3468](https://github.com/wailsapp/wails/pull/3468)
- New syntax for binding services: service instances must now be wrapped in a call to `application.NewService`. By [@fbbdev](https://github.com/fbbdev) in [#3468](https://github.com/wailsapp/wails/pull/3468)
- Modified the `contentTypeSniffer` struct to include the `http.CloseNotifier` interface. Now compatible with Gin framework. By [@AnalogJ](https://github.com/AnalogJ) in [#3537](https://github.com/wailsapp/wails/pull/3537)

### Removed

### Deprecated

### Security
