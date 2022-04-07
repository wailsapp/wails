# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v2.0.0-beta.33] - 2022-03-05

### Added

- NSIS Installer support for creating installers for Windows applications - Thanks [@stffabi](https://github.com/stffabi) 🎉
- New frontend:dev:watcher command to spin out 3rd party watchers when using wails dev - Thanks [@stffabi](https://github.com/stffabi)🎉
- Remote templates now support version tags - Thanks [@misitebao](https://github.com/misitebao) 🎉

### Fixed

- A number of fixes for ARM Linux providing a huge improvement - Thanks [@ianmjones](https://github.com/ianmjones) 🎉
- Fixed potential Nil reference when discovering the path to `index.html`
- Fixed crash when using `runtime.Log` methods in a production build
- Improvements to internal file handling meaning webworkers will now work on Windows - Thanks [@stffabi](https://github.com/stffabi)🎉

### Changed

- The Webview2 bootstrapper is now run as a normal user and doesn't require admin rights
- The docs have been improved and updated
- Added troubleshooting guide

[unreleased]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.33...HEAD
[v2.0.0-beta.33]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.32...v2.0.0-beta.33
